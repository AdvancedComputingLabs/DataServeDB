// Copyright (c) 2019 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

/*
Operations:
	- Loading
		- Create Table
		- Attach Table when table is mounted
	- Change Operations
		- Delete Table
		- Alter Table
*/

package dbtable

import (
	"DataServeDB/paths"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"DataServeDB/comminterfaces"
	"DataServeDB/commtypes"
)

//TODO: move it to error messages (single location)

//TODO: table main and table data are public now. Can be private?

type TableRow map[string]interface{} //it is by field name.

type createTableExternalStruct struct {
	_dbPtr       comminterfaces.DbPtrI // runtime hence _ prefix used.
	TableName    string
	TableStorers []StorerBasic
	TableFields  []string
}

func (t *createTableExternalStruct) AssignDb(dbPtr comminterfaces.DbPtrI) {
	t._dbPtr = dbPtr
}

type DbTable struct {
	TblMain              *tableMain //table structure information to keep it separate from data, so data disk io can be done separately.
	TblData              *tableDataContainer
	createTableStructure createTableExternalStruct
}

func (t *DbTable) GetId() int {
	return t.TblMain.TableId
}

func CreateTableJSON(jsonStr string, dbPtr comminterfaces.DbPtrI) (*DbTable, error) {

	var createTableData createTableExternalStruct
	if err := json.Unmarshal([]byte(jsonStr), &createTableData); err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return nil, errors.New("error found in table creation json")
	}

	createTableData._dbPtr = dbPtr

	tdc := &tableDataContainer{
		Rows:          nil,
		PkToRowMapper: map[interface{}]int64{},
	}

	//TODO: use global const, better practice
	//NOTE: it is creating table so id is -1.
	dbTable, err := createTable(-1, &createTableData, tdc)
	if err != nil {
		return dbTable, err
	}

	//TODO: dbTable.createTableStructure.TableStorers and createTableData.TableStorers are different instances after assignment.
	// FIX, later.

	return dbTable, err
}

func createTable(tableInternalId int, createTableData *createTableExternalStruct, tblDataContainer *tableDataContainer) (*DbTable, error) {
	// TODO: could be moved to table file. Name could be better.
	// I think it better belongs here than table.go as it is creating DbTable

	tbl := DbTable{}

	if tblMain, err := validateCreateTableMetaData(tableInternalId, createTableData); err != nil {
		return nil, err
	} else {
		tbl.TblMain = tblMain
	}

	tbl.TblData = tblDataContainer

	tbl.createTableStructure = *createTableData

	return &tbl, nil
}

func (t *DbTable) DebugPrintInternalFieldsNameMappings() string {
	return fmt.Sprintf("%#v", t.TblMain.TableFieldsMetaData.FieldNameToFieldInternalId)
}

func addIndices(table *DbTable, rowInternalIds tableRowByInternalIds, rowNumber int64) error {
	//NOTE: tableRowByInternalIds is passed by reference. -HY 22-Apr-2020

	// check the duplicate primary key before insert
	if _, ok := table.TblData.PkToRowMapper[rowInternalIds[table.TblMain.PkPos]]; ok {
		return errors.New("duplicate primary key")
	}

	table.TblData.PkToRowMapper[rowInternalIds[table.TblMain.PkPos]] = rowNumber

	return nil
}

func (t *DbTable) EventAfterTableIdAssignment() error {

	//IMP-NOTE: must be assigned only once when first table id is added.
	//TODO: check when table storers are mention in creation json.
	//TODO: check when table structure is loaded from disk.
	if len(t.createTableStructure.TableStorers) == 0 {
		fileName := fmt.Sprintf("table_%d.dat", t.TblMain.TableId)
		path := paths.Combine(t.createTableStructure._dbPtr.DbPath(), tablesDataPathRelative, fileName)

		diskStorePtr, err := NewDiskStoreV1(fileName, path)
		if err != nil {
			return err
		}
		t.createTableStructure.TableStorers = append(t.createTableStructure.TableStorers, diskStorePtr)
	}

	return nil
}

//type InsertRowCallback = func(table *DbTable, row *TableRowWithFieldProperties) error

func (t *DbTable) InsertRowJSON(jsonStr string /*, callback InsertRowCallback*/) error {

	var rowDataUnmarshalled TableRow
	if e := json.Unmarshal([]byte(jsonStr), &rowDataUnmarshalled); e != nil {
		_ = e
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return errors.New("error occured in parsing row json")
	}

	//returns validated row in the structure of internal ids
	_, rowInternalIds, e := validateRowData(t.TblMain, rowDataUnmarshalled)
	if e != nil {
		return e
	}

	numOfRows := int64(len(t.TblData.Rows))

	if e := addIndices(t, rowInternalIds, numOfRows); e != nil {
		return e
	}

	//TODO: TblData or Rows was giving error after loading table when data file was not there.
	//	There empty data case needs to be considered and dat file must be in the db for the table all the time?
	t.TblData.Rows = append(t.TblData.Rows, rowInternalIds)

	////TODO: path for table needs its own function?
	//fileName := fmt.Sprintf("table_%d.dat", t.TblMain.TableId)
	//path := paths.Combine(t.createTableStructure._dbPtr.DbPath(), tablesDataPathRelative, fileName)
	//
	////TODO: refector this into own function with binary or json option.
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	//err := enc.Encode(t.TblData)
	//if err != nil {
	//	//TODO: better handling needed
	//	//TODO: clean up event
	//	println("error ")
	//	log.Fatal("encode error:", err)
	//}
	//
	////TODO: disk error handling is needed.
	//storage.SaveToDisk(buf.Bytes(), path)

	rowWProps, e := GetTableRowWithFieldProperties(t, rowInternalIds)
	if e != nil {
		return e
	}

	//TODO: see below (memstore needs some thinking of indexing issues). Indices should be part of memstore or table's base structure?
	//WARNING AND NOTE: currently there is no memstore in the list, when there is this need to be handled differently.
	for _, store := range t.createTableStructure.TableStorers {
		//TODO: need to handle rollbacks
		//TODO: handle error
		store.Insert(rowWProps, t.TblData) //NOTE: t.TblData is a pointer.
	}

	return nil
}

func (t *DbTable) GetLength() int {
	//TODO: with empty there is potential it to panic with nil
	return int(int64(len(t.TblData.Rows)))
}

func (t *DbTable) GetRowByPrimaryKey(pkValue interface{}) (TableRow, error) {
	dbType, dbTypeProps := t.TblMain.getPkType()

	pkValueCasted, e := dbType.ConvertValue(pkValue, dbTypeProps)
	if e != nil {
		return nil, e
	}

	rowNum, exists := t.TblData.PkToRowMapper[pkValueCasted]
	if !exists {
		return nil, fmt.Errorf("value '%v' not found", pkValue)
	}

	row, e := toLabeledByFieldNames(t.TblData.Rows[rowNum], t.TblMain)
	if e != nil {
		return nil, e
	}

	return row, nil
}

func (t *DbTable) GetRowByPrimaryKeyReturnsJSON(pkValue interface{}) (string, error) {
	//TODO: Not sure this is really needed here.

	// fmt.Printf("here %v\n", t.tblData.Rows)
	row, e := t.GetRowByPrimaryKey(pkValue)
	if e != nil {
		return "", e
	}

	jsonBytes, e := json.Marshal(row)
	if e != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return "", fmt.Errorf("error occured while converting row data to json; primary key value '%v'", pkValue)
	}

	return string(jsonBytes), nil
}

func (t *DbTable) Get(dbReqCtx *commtypes.DbReqContext) (resultHttpStatus int, resultContent []byte, resultErr error) {
	key, value, err := parseKeyValue(dbReqCtx.ResPath)
	if err != nil {
		resultErr = err
		return
	}

	//NOTE: there is no empty index 'indexName:' access at the moment.
	if value == "" {
		resultErr = errors.New("value is not provided")
	}

	if key == "" {
		//primary key
		row, err := t.GetRowByPrimaryKeyReturnsJSON(value)
		if err != nil {
			resultErr = err
			return
		}

		resultContent = []byte(row)
		resultHttpStatus = http.StatusOK
	} else {
		//
	}

	return
}
