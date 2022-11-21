package dbtable_disk_store_v1

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"DataServeDB/commtypes"
	storage_utils "DataServeDB/dbsystem/dbstorage"
	"DataServeDB/paths"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/utils/rest/dberrors"
)

const StorerDiskV1Name = "StorerDiskV1"

var DeleteRollbackTest bool = false
var InsertRollbackTest bool = false
var UpdateRollbackTest bool = false

type DiskStoreV1 struct {
	fileName string
	path     string
}

func NewDiskStoreV1(tableId int, tableFolderPath string, _ idbstorer.TableInfoGetter) (idbstorer.StorerBasic, *dberrors.DbError) {
	fileName := fmt.Sprintf("table_%d_%s.dat", tableId, StorerDiskV1Name)
	path := paths.Combine(tableFolderPath, fileName)
	return &DiskStoreV1{fileName: fileName, path: path}, nil
}

func (d DiskStoreV1) DisplayName() string {
	return StorerDiskV1Name
}

func (d DiskStoreV1) FeaturesAndConstraints(featureOrConstraint idbstorer.TableStorerFeaturesType) bool {

	switch featureOrConstraint {
	case idbstorer.TableStorerFeature_Insert, idbstorer.TableStorerFeature_Update, idbstorer.TableStorerFeature_Delete,
		idbstorer.TableStorerConstraint_MustNotBeFirstStore, idbstorer.TableStorerConstraint_MustBeLastStore:
		return true
	}

	return false
}

func (d DiskStoreV1) CreateInit(prevStore, nextStore idbstorer.StorerBasic) *dberrors.DbError {

	if dberr := idbstorer.CheckTableStorageConstraints(d, prevStore, nextStore); dberr != nil {
		return dberr
	}

	return saveToFile("", d.path)
}

func (d DiskStoreV1) LoadTable(input any, prevStore, nextStore idbstorer.StorerBasic) *idbstorer.StorerLoadResult {

	dataBytes, err := storage_utils.LoadFromDisk(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &idbstorer.StorerLoadResult{}
		}
		return &idbstorer.StorerLoadResult{Err: err}
	}

	return &idbstorer.StorerLoadResult{Data: dataBytes}
}

func (d DiskStoreV1) DeleteTable(currentStoreIndex int, storesResults []*idbstorer.StorerDeleteTableResult) idbstorer.StorerDeleteTableResult {
	if err := storage_utils.DeleteDirFromDisk(d.path); err != nil {
		return idbstorer.StorerDeleteTableResult{DbErr: dberrors.NewDbError(dberrors.InternalServerErrorDiskError, err)}
	}
	return idbstorer.StorerDeleteTableResult{}
}

func (d DiskStoreV1) DeleteTableRollback(currentStoreIndex int, storesResults []*idbstorer.StorerDeleteTableResult) *dberrors.DbError {
	return nil
}

func (d DiskStoreV1) Delete(indexColId int, key any,
	currentStoreIndex int, storesResults []*idbstorer.StorerDeleteResult) idbstorer.StorerDeleteResult {

	if DeleteRollbackTest {
		return idbstorer.StorerDeleteResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.TestError, errors.New("disk storage test error"))}
	}

	//IMPORTANT: delete and delete rollback does not preserve the order of the rows for performance reasons.
	// If order is needed then sequential index should be used.

	// deleted row is removed from the table data attached, just save the whole table data to disk.
	// not efficient, but for now it is ok.
	err := saveToDisk(currentStoreIndex, storesResults, d.path)
	if err != nil {
		return idbstorer.StorerDeleteResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.InternalServerErrorDiskError, err)}
	}

	return idbstorer.StorerDeleteResult{Count: 1}
}

func (d DiskStoreV1) DeleteRollback(indexColId int, currentStoreIndex int, storesResults []*idbstorer.StorerDeleteResult) *dberrors.DbError {
	//NOTE: if this is last storage then doesn't need to implement rollback.
	panic("CodingError: DeleteRollback() should not be called.")
}

func (d DiskStoreV1) Get(indexColId int, key any) (int, commtypes.TableRowByFieldsIds, *dberrors.DbError) {
	// NOTE: should not be first storage. First storage with Get() method implementation is used.
	panic("CodingError: Get() should not be called.")
}

func (d DiskStoreV1) Insert(rowWithProps commtypes.TableRowWithFieldProperties, currentStoreIndex int, storesResults []*idbstorer.StorerInsertResult) idbstorer.StorerInsertResult {

	if InsertRollbackTest {
		// TODO: find better way to test storage errors.
		return idbstorer.StorerInsertResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.TestError, errors.New("disk storage test error"))}
	}

	// new inserted row is in the table data attached, just save the whole table data to disk.
	// not efficient, but for now it is ok.
	err := saveToDisk(currentStoreIndex, storesResults, d.path)
	if err != nil {
		return idbstorer.StorerInsertResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.InternalServerErrorDiskError, err)}
	}

	return idbstorer.StorerInsertResult{Count: 1}
}

func (d DiskStoreV1) InsertRollback(currentStoreIndex int, storesResults []*idbstorer.StorerInsertResult) *dberrors.DbError {
	//NOTE: if this is last storage then doesn't need to implement rollback.
	panic("CodingError: InsertRollback() should not be called.")
}

func (d DiskStoreV1) Update(indexColId int, key any, rowWithProps commtypes.TableRowWithFieldProperties, updateType idbstorer.TableOperationType,
	currentStoreIndex int, storesResults []*idbstorer.StorerUpdateResult) idbstorer.StorerUpdateResult {

	if UpdateRollbackTest {
		return idbstorer.StorerUpdateResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.TestError, errors.New("disk storage test error"))}
	}

	// NOTE: see comment in Insert() method.
	err := saveToDisk(currentStoreIndex, storesResults, d.path)
	if err != nil {
		return idbstorer.StorerUpdateResult{Count: 0, DbErr: dberrors.NewDbError(dberrors.InternalServerErrorDiskError, err)}
	}

	return idbstorer.StorerUpdateResult{Count: 1}
}

func (d DiskStoreV1) UpdateRollback(indexColId int, currentStoreIndex int, storesResults []*idbstorer.StorerUpdateResult) *dberrors.DbError {
	//NOTE: if this is last storage then doesn't need to implement rollback.
	panic("CodingError: UpdateRollback() should not be called.")
}

func (d DiskStoreV1) GetNumberOfRows() int {
	// NOTE: should not be first storage. First storage with  GetNumberOfRows() method implementation is used.
	panic("CodingError: GetNumberOfRows() should not be called.")
}

// ### Private:

func saveToDisk[T idbstorer.StorerDeleteResult | idbstorer.StorerInsertResult | idbstorer.StorerUpdateResult](
	currentStoreIndex int, storesResults []*T, path string) error {

	if currentStoreIndex == 0 {
		panic("'DiskStoreV1' can't be the first store. Currently it does not handle inserts directly as first store")
	}

	var previousStoreTableData any = nil

	switch v := any(storesResults).(type) {
	case []*idbstorer.StorerInsertResult:
		previousStoreTableData = v[currentStoreIndex-1].Data
	case []*idbstorer.StorerDeleteResult:
		previousStoreTableData = v[currentStoreIndex-1].Data
	case []*idbstorer.StorerUpdateResult:
		previousStoreTableData = v[currentStoreIndex-1].Data
	}

	if previousStoreTableData == nil {
		// NOTE: data has not changed, so no need to put table in invalid state.
		return errors.New("CodingError: previousStoreTableData is nil")
	}

	err2 := saveToFile(previousStoreTableData, path)
	if err2 != nil {
		return err2.ToError()
	}

	return nil
}

func saveToFile(data any, path string) *dberrors.DbError {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data) //NOTE: currently this saves all data and index structure to the disk.
	if err != nil {
		// TODO: check is it always invalid input error?
		return dberrors.NewDbError(dberrors.InvalidInput, errors.New(fmt.Sprintf("error while encoding data: %v", err)))
	}

	err = storage_utils.SaveToDisk(buf.Bytes(), path)
	if err != nil {
		return dberrors.NewDbError(dberrors.InternalServerError, errors.New(fmt.Sprintf("error while saving data to disk: %v", err)))
	}
	return nil
}
