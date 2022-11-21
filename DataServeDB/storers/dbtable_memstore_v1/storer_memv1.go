package dbtable_memstore_v1

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"DataServeDB/commtypes"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/utils/rest/dberrors"
)

const StorerMemV1Name = "StorerMemV1"

type BitsFlag uint32

const (
	RowFlagMarkedForDeletion BitsFlag = 1 << iota
)

type RowFieldsVersion = map[int]any

type TableRow struct {
	RowFlags               BitsFlag
	RowCurrentVersion      int
	RowFieldsDataByVersion []RowFieldsVersion
}

// MemStoreV1DataContainer NOTE: memstorev1 can change so better save it in independent struct
type MemStoreV1DataContainer struct {
	Rows          []TableRow
	PkToRowMapper map[any]int64
}

type MemStoreV1 struct {
	memOnlyTable    bool                      // runtime only, not saved to disk (hence private); TODO: check it is not saved to the disk.
	tableInfoGetter idbstorer.TableInfoGetter // runtime only, not saved to disk (hence private)
	Rows            []TableRow
	PkToRowMapper   map[any]int64
}

type RowHolder struct {
	PkValue any
	Row     TableRow
}

type myRollbackItem struct {
	ColIdOfPk int
	PkValue   any
}

func NewMemStoreV1(tableId int, tableFolderPath string, tableInfo idbstorer.TableInfoGetter) (idbstorer.StorerBasic, *dberrors.DbError) {
	return &MemStoreV1{tableInfoGetter: tableInfo, Rows: []TableRow{}, PkToRowMapper: map[any]int64{}}, nil
}

func (m *MemStoreV1) DisplayName() string {
	return StorerMemV1Name
}

func (m *MemStoreV1) FeaturesAndConstraints(feature idbstorer.TableStorerFeaturesType) bool {

	switch feature {
	case idbstorer.TableStorerFeature_Insert, idbstorer.TableStorerFeature_Update,
		idbstorer.TableStorerFeature_Delete, idbstorer.TableStorerFeature_Get:
		return true
	}

	return false
}

func (m *MemStoreV1) CreateInit(prevStore, nextStore idbstorer.StorerBasic) *dberrors.DbError {
	if nextStore == nil {
		m.memOnlyTable = true
	}
	return nil
}

func (m *MemStoreV1) LoadTable(input any, prevStore, nextStore idbstorer.StorerBasic) *idbstorer.StorerLoadResult {

	// Cases:
	// 1. input is nil, mem only table, no data to load.
	// 2. input is came from another storage, load data from it. Should be of type MemStoreV1DataContainer.

	// identify if mem only table
	if nextStore == nil {
		m.memOnlyTable = true
	}

	if input == nil {
		if m.memOnlyTable {
			// mem only table, no data to load.
			return &idbstorer.StorerLoadResult{}
		}

		return &idbstorer.StorerLoadResult{Err: errors.New("CodingError: input is nil and table is not mem only")}
	}

	if m.memOnlyTable {
		return &idbstorer.StorerLoadResult{Err: errors.New("CodingError: input is not nil and table is mem only")}
	}

	var b []byte

	switch v := input.(type) {
	case []byte:
		b = v
	default:
		panic("CodingError: 'input' must be byte array")
	}

	// TODO: can refectored to its own util function?
	//  Perhaps make a templated util function decodeGobFromBytes?

	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	var tableData MemStoreV1DataContainer

	err := dec.Decode(&tableData)
	if err != nil {
		return &idbstorer.StorerLoadResult{Err: err}
	}

	m.Rows = tableData.Rows
	m.PkToRowMapper = tableData.PkToRowMapper

	return &idbstorer.StorerLoadResult{}
}

func (m *MemStoreV1) DeleteTable(currentStoreIndex int, storesResults []*idbstorer.StorerDeleteTableResult) idbstorer.StorerDeleteTableResult {
	return idbstorer.StorerDeleteTableResult{}
}

func (m *MemStoreV1) DeleteTableRollback(currentStoreIndex int, storesResults []*idbstorer.StorerDeleteTableResult) *dberrors.DbError {
	return nil
}

func (m *MemStoreV1) Delete(indexColId int, key any,
	currentStoreIndex int, storesResults []*idbstorer.StorerDeleteResult) idbstorer.StorerDeleteResult {

	//IMP-NOTE: indexColId is primary key at the moment as only primary key is used for indexing.

	rowHolderPtr, err := deleteRow(m, indexColId, key)
	if err != nil {
		return idbstorer.StorerDeleteResult{Count: 0, DbErr: err}
	}

	tableData := MemStoreV1DataContainer{Rows: m.Rows, PkToRowMapper: m.PkToRowMapper}

	return idbstorer.StorerDeleteResult{RollbackItem: rowHolderPtr, Count: 1, Data: &tableData}
}

func (m *MemStoreV1) DeleteRollback(indexColId int, currentStoreIndex int, storesResults []*idbstorer.StorerDeleteResult) *dberrors.DbError {

	// IMP-NOTE: indexColId is primary key at the moment as only primary key is used for indexing.

	rowHolderPtr := storesResults[currentStoreIndex].RollbackItem
	if rowHolderPtr == nil {
		panic("CodingError: rowHolderPtr is nil")
	}

	rowHolder, ok := rowHolderPtr.(*RowHolder)
	if !ok {
		panic("CodingError: rowHolderPtr is not of type *RowHolder")
	}

	m.Rows = append(m.Rows, rowHolder.Row)
	m.PkToRowMapper[rowHolder.PkValue] = int64(len(m.Rows) - 1)

	return nil
}

func (m *MemStoreV1) Get(indexColId int, key any) (int, commtypes.TableRowByFieldsIds, *dberrors.DbError) {
	//indexColId = -1, not used. Currently, uses only primary key.
	rowPos, exists := m.PkToRowMapper[key]
	if !exists {
		return 0, nil, dberrors.NewDbError(dberrors.KeyNotFound, fmt.Errorf("key '%v' not found", key))
	}

	//TODO: check rowPos exists and handle if not? It panic/error if storage is corrupted.
	// Storage should be declared invalid?
	rowRawPtr := &m.Rows[rowPos]
	rowByColsIds := toTableRowByColsIds(rowRawPtr)
	return 1, rowByColsIds, nil
}

func (m *MemStoreV1) Insert(rowWithProps commtypes.TableRowWithFieldProperties,
	currentStoreIndex int, storesResults []*idbstorer.StorerInsertResult) idbstorer.StorerInsertResult {

	row := newTableRow()
	colIdOfPk := -1

	rowFieldsData := RowFieldsVersion{}

	for colId, item := range rowWithProps {
		rowFieldsData[colId] = item.V
		if item.IsPk() {
			colIdOfPk = colId
		}
	}

	//TODO: if colIdOfPk == -1 return error
	if colIdOfPk == -1 {
		// NOTE: I don't think table needs to be put in invalid state, as storage data has not changed.
		panic("CodingError: colIdOfPk is -1")
	}

	row.RowFieldsDataByVersion = append(row.RowFieldsDataByVersion, rowFieldsData)
	//NOTE: current version is 0, so no need to set it.

	numOfRows := int64(len(m.Rows))

	if dberr := addIndices(m, row, colIdOfPk, numOfRows); dberr != nil {
		return idbstorer.StorerInsertResult{Count: 0, DbErr: dberr}
	}

	//TODO: TblData or Rows was giving error after loading table when data file was not there.
	//	There empty data case needs to be considered and dat file must be in the db for the table all the time?
	m.Rows = append(m.Rows, row)

	//IMP NOTE: PkToRowMapper is by reference so saving earlier does not work.
	//			'Rows' object: new row does not show, but if I change the content in m.Rows, it shows
	//			in Rows if it is copied earlier. So Rows internally seems by reference too.

	//NOTE: Rollback. RollbackItem has been added to insert result struct for rollback. It is using PkValue for rollback.

	//NOTE: seemed simpler just to pass the data in the result to the next storer. But could have handled here to pass it to disk store.
	//NOTE: these operations are ok for performance as both are maps and maps are reference types.
	tableData := MemStoreV1DataContainer{Rows: m.Rows, PkToRowMapper: m.PkToRowMapper}

	rbitem := myRollbackItem{ColIdOfPk: colIdOfPk, PkValue: row.RowFieldsDataByVersion[row.RowCurrentVersion][colIdOfPk]}

	return idbstorer.StorerInsertResult{RollbackItem: rbitem, Count: 1, Data: &tableData}
}

func (m *MemStoreV1) InsertRollback(currentStoreIndex int, storesResults []*idbstorer.StorerInsertResult) *dberrors.DbError {

	rbitem, ok := storesResults[currentStoreIndex].RollbackItem.(myRollbackItem)
	if !ok {
		// TODO: should shutdown the db or put the table in invalid state?
		panic("CodingError: rollback item is not of type myRollbackItem")
	}

	if rbitem.PkValue == nil {
		// TODO: should shutdown the db or put the table in invalid state?
		panic("CodingError: pkValue is nil")
	}

	// do not need row holder as it is not used.
	if _, dberr := deleteRow(m, rbitem.ColIdOfPk, rbitem.PkValue); dberr != nil {
		return dberr
	}

	return nil
}

func (m *MemStoreV1) Update(indexColId int, key any, rowWithProps commtypes.TableRowWithFieldProperties, updateType idbstorer.TableOperationType,
	currentStoreIndex int, storesResults []*idbstorer.StorerUpdateResult) idbstorer.StorerUpdateResult {

	//IMP-NOTE: indexColId is primary key at the moment as only primary key is used for indexing.

	rowPos, exists := m.PkToRowMapper[key]
	if !exists {
		return idbstorer.StorerUpdateResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.KeyNotFound, fmt.Errorf("key '%v' not found", key))}
	}

	//TODO: check rowPos exists and handle if not? It panic/error if storage is corrupted.
	// Storage should be declared invalid?

	row := m.Rows[rowPos]

	var rowDataNew RowFieldsVersion

	switch updateType {
	case idbstorer.TableOperationPatchRow:
		{
			err := deepCopy(&rowDataNew, &row.RowFieldsDataByVersion[row.RowCurrentVersion])
			if err != nil {
				return idbstorer.StorerUpdateResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.InternalServerError, err)}
			}

			for fieldId, value := range rowWithProps {
				rowDataNew[fieldId] = value.Value()
			}
		}

	case idbstorer.TableOperationReplaceRow:
		{
			rowDataNew = RowFieldsVersion{}

			// get primary key column id
			if tableInfo, err := m.tableInfoGetter.GetTableInfo(); err != nil {
				return idbstorer.StorerUpdateResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.InternalServerError, err)}
			} else {
				// add primary key to the row data
				rowDataNew[tableInfo.PrimaryKeyColId] = row.RowFieldsDataByVersion[row.RowCurrentVersion][tableInfo.PrimaryKeyColId]
			}

			for fieldId, value := range rowWithProps {
				rowDataNew[fieldId] = value.Value()
			}
		}

	default:
		panic("CodingError: only patch and replace row update types are supported")
	}

	row.RowFieldsDataByVersion = append(row.RowFieldsDataByVersion, rowDataNew)
	row.RowCurrentVersion = len(row.RowFieldsDataByVersion) - 1

	m.Rows[rowPos] = row

	// NOTE: see comment in Insert().
	tableData := MemStoreV1DataContainer{Rows: m.Rows, PkToRowMapper: m.PkToRowMapper}

	return idbstorer.StorerUpdateResult{Count: 1, Data: &tableData, RollbackItem: key}
}

func (m *MemStoreV1) UpdateRollback(indexColId int, currentStoreIndex int, storesResults []*idbstorer.StorerUpdateResult) *dberrors.DbError {

	pkValue := storesResults[currentStoreIndex].RollbackItem
	if pkValue == nil {
		panic("CodingError: pkValue is nil")
	}

	row := m.Rows[m.PkToRowMapper[pkValue]]

	row.RowFieldsDataByVersion = row.RowFieldsDataByVersion[:row.RowCurrentVersion]
	row.RowCurrentVersion = len(row.RowFieldsDataByVersion) - 1

	// TODO: if current version is -1, it is coding error?

	m.Rows[m.PkToRowMapper[pkValue]] = row

	return nil
}

func (m *MemStoreV1) GetNumberOfRows() int {
	return int(int64(len(m.Rows)))
}

// ### Private Section:

func addIndices(m *MemStoreV1, row TableRow, colIdOfPk int, rowNumber int64) *dberrors.DbError {
	//NOTE: tableRowByInternalIds is passed by reference. -HY 22-Apr-2020

	if colIdOfPk == -1 {
		panic("something wrong in the code, should have primary key col id")
	}

	// check the duplicate primary key before insert
	if _, ok := m.PkToRowMapper[row.RowFieldsDataByVersion[row.RowCurrentVersion][colIdOfPk]]; ok {
		return dberrors.NewDbError(dberrors.InvalidInputDuplicateKey, errors.New("duplicate primary key"))
	}

	m.PkToRowMapper[row.RowFieldsDataByVersion[row.RowCurrentVersion][colIdOfPk]] = rowNumber

	return nil
}

func deepCopy(dstPtr, srcPtr any) (err error) {
	buf := bytes.Buffer{}
	if err = gob.NewEncoder(&buf).Encode(srcPtr); err != nil {
		return
	}
	return gob.NewDecoder(&buf).Decode(dstPtr)
}

func deleteRow(m *MemStoreV1, pkColId int, pkValue any) (*RowHolder, *dberrors.DbError) {

	rowNumber, exists := m.PkToRowMapper[pkValue]
	if !exists {
		//TODO: error message to differentiate between primary key vs other index?
		return nil, dberrors.NewDbError(dberrors.KeyNotFound, errors.New("key not found"))
	}

	numOfRows := int64(len(m.Rows))

	if rowNumber >= numOfRows {
		panic("something wrong in the code, should have row here")
	}

	//save deleting row item
	rh := RowHolder{PkValue: pkValue, Row: m.Rows[rowNumber]}

	if rowNumber == numOfRows-1 {
		//it is last row, just truncate.
		m.Rows = m.Rows[:rowNumber]
		delete(m.PkToRowMapper, pkValue)
	} else {
		//swap with last row and truncate.
		m.Rows[rowNumber] = m.Rows[numOfRows-1]
		m.Rows = m.Rows[:numOfRows-1]

		//update the swapped row number in the mapper
		row := m.Rows[rowNumber]
		m.PkToRowMapper[row.RowFieldsDataByVersion[row.RowCurrentVersion][pkColId]] = rowNumber

		delete(m.PkToRowMapper, pkValue)
	}

	return &rh, nil
}

func newTableRow() TableRow {
	return TableRow{RowCurrentVersion: 0, RowFieldsDataByVersion: make([]RowFieldsVersion, 0)}
}

func toTableRowByColsIds(rowPtr *TableRow) commtypes.TableRowByFieldsIds {
	//TODO: need to check if rowPtr is nil? -HY 14-Oct-2022

	//TODO: make rowPtr.RowFieldsDataByVersion[rowPtr.RowCurrentVersion] more readable
	//	by using method? -HY 15-Oct-2022

	if rowPtr.RowCurrentVersion < 0 {
		panic("CodingError: rowPtr.RowCurrentVersion < 0")
	}

	rowByColsIds := make(commtypes.TableRowByFieldsIds, len(rowPtr.RowFieldsDataByVersion[rowPtr.RowCurrentVersion]))
	for colId, value := range rowPtr.RowFieldsDataByVersion[rowPtr.RowCurrentVersion] {
		rowByColsIds[colId] = value
	}
	return rowByColsIds
}
