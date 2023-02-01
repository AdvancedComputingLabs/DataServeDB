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
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"DataServeDB/comminterfaces"
	"DataServeDB/commtypes"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/utils/rest"
	"DataServeDB/utils/rest/dberrors"
)

//TODO: move it to error messages (single location)

//TODO: table main and table data are public now. Can be private?

type TableRow map[string]interface{} //it is by field name.
//type TableRowWithFieldProperties tableRowByInternalIdsWithFieldProperties

type createTableExternalStruct struct {
	_dbPtr                 comminterfaces.DbPtrI // runtime hence _ prefix used.
	_tableStorersInstances []idbstorer.StorerBasic
	TableName              string
	TableStorages          string
	TableColumns           []string
}

// ## Public:

// ### JSON?:

func (t *createTableExternalStruct) AssignDb(dbPtr comminterfaces.DbPtrI) {
	t._dbPtr = dbPtr
}

type DbTable struct {
	TblMain              *tableMain //table structure information to keep it separate from data, so data disk io can be done separately.
	mu                   *sync.RWMutex
	createTableStructure createTableExternalStruct
}

func (t *DbTable) GetId() int {
	return t.TblMain.TableId
}

func (t *DbTable) GetDbTypeDisplayName() string {
	return "dbtable"
}

func CreateTableJSON(jsonStr string, dbPtr comminterfaces.DbPtrI) (*DbTable, *dberrors.DbError) {

	var createTableData createTableExternalStruct
	if err := json.Unmarshal([]byte(jsonStr), &createTableData); err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return nil, dberrors.NewDbError(dberrors.InvalidInput, errors.New("error found in table creation json"))
	}

	createTableData._dbPtr = dbPtr

	//TODO: use global const, better practice
	//NOTE: it is creating table so id is -1.
	dbTable, dberr := createTable(-1, &createTableData)
	if dberr != nil {
		return dbTable, dberr
	}

	//TODO: dbTable.createTableStructure._tableStorersInstances and createTableData._tableStorersInstances are different instances after assignment.
	// FIX, later.

	return dbTable, dberr
}

func DeleteTable(table *DbTable) *dberrors.DbError {

	numberOfStores := len(table.createTableStructure._tableStorersInstances)

	if numberOfStores > 0 {
		var storerResults = make([]*idbstorer.StorerDeleteTableResult, numberOfStores)

		for j := numberOfStores - 1; j >= 0; j-- {
			result := table.createTableStructure._tableStorersInstances[j].DeleteTable(j, storerResults)
			storerResults[j] = &result
			if result.DbErr != nil {
				if dberr := deleteTableRollback(table, storerResults, j); dberr != nil {
					// TODO: invalid state marker for table.
					return dberr
				}
				// TODO: which error to return? There are errors from multiple storages? Maybe not as first error is returned?
				return result.DbErr
			}
		}
	}

	return nil
}

func createTable(tableInternalId int, createTableData *createTableExternalStruct) (*DbTable, *dberrors.DbError) {
	// TODO: could be moved to table file. Name could be better.
	// I think it better belongs here than table.go as it is creating DbTable

	tbl := DbTable{}

	if tblMain, err := validateCreateTableMetaData(tableInternalId, createTableData); err != nil {
		return nil, dberrors.NewDbError(dberrors.InvalidInput, err)
	} else {
		tbl.TblMain = tblMain
	}

	tbl.createTableStructure = *createTableData

	return &tbl, nil
}

func (t *DbTable) DebugPrintInternalFieldsNameMappings() string {
	return fmt.Sprintf("%#v", t.TblMain.TableFieldsMetaData.FieldNameToFieldInternalId)
}

// EventAfterTableIdAssignment It is called from db package, hence, exported. It must be
// called from the db package, because table id is assigned there at table creation time.
func (t *DbTable) EventAfterTableIdAssignment() *dberrors.DbError {
	//IMP-NOTE: must be assigned only once when first table id is added.

	//NOTE: dberr not returned directly; future-proof when it may have more than one function call.
	if dberr := applyStorages(t, true); dberr != nil {
		return dberr
	}

	return nil
}

func GetTableRowWithFieldProperties(table *DbTable, rowByInternalIds tableRowByInternalIds) (commtypes.TableRowWithFieldProperties, error) {
	var meta = &table.TblMain.TableFieldsMetaData
	rowWithProps := commtypes.TableRowWithFieldProperties{}

	meta.mu.RLock()
	defer meta.mu.RUnlock()

	for k, v := range rowByInternalIds {
		if fieldProps, exits := meta.FieldInternalIdToFieldMetaData[k]; exits {
			//NOTE: this does not need conversion because this will probably come from db internal row and will have correct type.
			//But check/test.
			rowWithProps[k] = commtypes.FieldValueAndPropertiesHolder{V: v, TableFieldInternal: fieldProps} //NOTE: used field name stored.
		} else {
			//TODO: Log.
			//TODO: Test if error or panic operations are atomic.
			//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
			panic("this is internal error, shouldn't happen")
			//return nil, errors.New("this is internal error, shouldn't happen")
		}
	}

	return rowWithProps, nil
}

// delete row json from table.

func (t *DbTable) DeleteRowByPrimaryKey(pkValue any) *dberrors.DbError {

	t.mu.Lock()
	defer t.mu.Unlock()

	dbType, dbTypeProps := t.TblMain.getPkType()

	//TODO: 'castPkValue' seems to do samething. Check and remove.
	pkValueCasted, err := dbType.ConvertValue(pkValue, dbTypeProps)
	if err != nil {
		return dberrors.NewDbError(dberrors.InvalidInput, err)
	}

	numberOfStores := len(t.createTableStructure._tableStorersInstances)

	if numberOfStores > 0 {
		var storerResults = make([]*idbstorer.StorerDeleteResult, numberOfStores)

		pkColId := t.TblMain.PkPos // currently only pk is supported for indexing.

		for i, store := range t.createTableStructure._tableStorersInstances {
			result := store.Delete(pkColId, pkValueCasted, i, storerResults)
			storerResults[i] = &result

			//TODO: check if these errors are correct to pass to user or they are internal errors?

			if result.DbErr != nil {
				dberr := deleteRowRollback(t, pkColId, storerResults, i)
				if dberr != nil {
					//TODO: invalid state marker for table.
					return dberr
				}
				//TODO: probably last error is returned. Error from which store should be returned?
				return result.DbErr
			}
		}
	}

	return nil
}

func (t *DbTable) UpdateRowJsonByPk(pkValue any, jsonStr string, updateType idbstorer.TableOperationType) *dberrors.DbError {
	// TODO: test both patch and replace if there behaviour is according to spec.

	t.mu.Lock()
	defer t.mu.Unlock()

	var rowUpdateDataUnmarshalled TableRow
	if err := json.Unmarshal([]byte(jsonStr), &rowUpdateDataUnmarshalled); err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return dberrors.NewDbError(dberrors.InvalidInput, errors.New("error found in row update json"))
	}

	//validate row update data.
	_, rowUpdateDateInternalIds, dberr := validateRowData(t.TblMain, rowUpdateDataUnmarshalled, updateType)
	if dberr != nil {
		return dberr
	}

	//do row with properties.
	rowWithProps, err := GetTableRowWithFieldProperties(t, rowUpdateDateInternalIds)
	if err != nil {
		return dberrors.NewDbError(dberrors.InternalServerError, err)
	}

	pkValueCasted, dberr := castPkValue(pkValue, t)
	if dberr != nil {
		return dberr
	}

	numberOfStores := len(t.createTableStructure._tableStorersInstances)

	if numberOfStores > 0 {
		var storerResults = make([]*idbstorer.StorerUpdateResult, numberOfStores)

		pkColId := t.TblMain.PkPos // currently only pk is supported for indexing.

		for i, store := range t.createTableStructure._tableStorersInstances {
			result := store.Update(pkColId, pkValueCasted, rowWithProps, updateType, i, storerResults)
			storerResults[i] = &result

			//TODO: check if these errors are correct to pass to user or they are internal errors?

			if result.DbErr != nil {
				dberr = updateRowRollback(t, pkColId, storerResults, i)
				if dberr != nil {
					//TODO: invalid state marker for table.
					return dberr
				}
				return result.DbErr
			}
		}
	}

	return nil
}

func (t *DbTable) InsertRowJSON(jsonStr string) *dberrors.DbError {

	t.mu.Lock()
	defer t.mu.Unlock()

	var rowDataUnmarshalled TableRow
	if err := json.Unmarshal([]byte(jsonStr), &rowDataUnmarshalled); err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return dberrors.NewDbError(dberrors.InvalidInput, errors.New("error found in parsing row json"))
	}

	//returns validated row in the structure of internal ids
	_, rowInternalIds, dberr := validateRowData(t.TblMain, rowDataUnmarshalled, idbstorer.TableOperationInsertRow)
	if dberr != nil {
		return dberr
	}

	rowWProps, err := GetTableRowWithFieldProperties(t, rowInternalIds)
	if err != nil {
		return dberrors.NewDbError(dberrors.InternalServerError, err)
	}

	numberOfStores := len(t.createTableStructure._tableStorersInstances)

	if numberOfStores > 0 {
		var storerResults = make([]*idbstorer.StorerInsertResult, numberOfStores)

		for i, store := range t.createTableStructure._tableStorersInstances {
			result := store.Insert(rowWProps, i, storerResults)
			storerResults[i] = &result

			//TODO: check if these errors are correct to pass to user or they are internal errors?

			if result.DbErr != nil {
				dberr := insertRowRollback(t, storerResults, i)
				if dberr != nil {
					//TODO: invalid state marker for table.
					return dberr
				}
				return result.DbErr
			}
		}
	}

	return nil
}

func updateRowRollback(t *DbTable, indexColId int, results []*idbstorer.StorerUpdateResult, i int) *dberrors.DbError {

	storages := t.createTableStructure._tableStorersInstances

	for j := i - 1; j >= 0; j-- {
		dberr := storages[j].UpdateRollback(indexColId, j, results)
		if dberr != nil {
			return dberr
		}
	}

	return nil
}

func deleteRowRollback(t *DbTable, indexColId int, results []*idbstorer.StorerDeleteResult, i int) *dberrors.DbError {

	storages := t.createTableStructure._tableStorersInstances

	for j := i - 1; j >= 0; j-- {
		dberr := storages[j].DeleteRollback(indexColId, j, results)
		if dberr != nil {
			//append error to error list, continue rollback for other storages or return error?
			//TODO: check back logic
			return dberr
		}
	}

	return nil
}

func deleteTableRollback(t *DbTable, results []*idbstorer.StorerDeleteTableResult, j int) *dberrors.DbError {

	storages := t.createTableStructure._tableStorersInstances

	for i := j; i < len(storages); i++ {
		dberr := storages[i].DeleteTableRollback(i, results)
		if dberr != nil {
			//append error to error list, continue rollback for other storages or return error?
			// TODO: check back logic
			return dberr
		}
	}

	return nil
}

func insertRowRollback(t *DbTable, results []*idbstorer.StorerInsertResult, i int) *dberrors.DbError {
	storages := t.createTableStructure._tableStorersInstances

	for j := i - 1; j >= 0; j-- {
		dberr := storages[j].InsertRollback(j, results)
		if dberr != nil {
			//append error to error list, continue rollback for other storages or return error?
			//TODO: check back logic
			return dberr
		}
	}

	return nil
}

func (t *DbTable) GetLength() int {

	//TODO: why first store? is it by convention (first) or must declared?
	// For now I'm using position 0
	numOfRows := t.createTableStructure._tableStorersInstances[0].GetNumberOfRows()
	return numOfRows
}

func (t *DbTable) GetRowByPrimaryKey(pkValue any) (TableRow, *dberrors.DbError) {
	pkValueCasted, dberr := castPkValue(pkValue, t)
	if dberr != nil {
		return nil, dberr
	}

	pkColId := t.TblMain.PkPos // currently only pk is supported for indexing.

	//TODO: why first store? is it by convention (first) or must declared?
	// For now I'm using position 0
	_, rowRaw, dberr := t.createTableStructure._tableStorersInstances[0].Get(pkColId, pkValueCasted)
	if dberr != nil {
		return nil, dberr
	}

	row, dberr := toLabeledByFieldNames(rowRaw, t.TblMain)
	if dberr != nil {
		return nil, dberr
	}

	return row, nil
}

func castPkValue(pkValue any, t *DbTable) (any, *dberrors.DbError) {
	dbType, dbTypeProps := t.TblMain.getPkType()

	pkValueCasted, err := dbType.ConvertValue(pkValue, dbTypeProps)
	if err != nil {
		//TODO: may not be from input, may be from internal error.
		// Test cases: (1) pk value is not of type of pk column.
		// TODO: make error code specific to key type mismatch? It would be better to have a specific error code for this?
		return nil, dberrors.NewDbError(dberrors.InvalidInput, err)
	}
	return pkValueCasted, nil
}

func (t *DbTable) GetRowByPrimaryKeyReturnsJSON(pkValue any) (string, *dberrors.DbError) {
	//TODO: Not sure this is really needed here.

	row, dberr := t.GetRowByPrimaryKey(pkValue)
	if dberr != nil {
		return "", dberr
	}

	jsonBytes, err := json.Marshal(row)
	if err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		err = fmt.Errorf("error occured while converting row data to json; primary key value '%v'; internal error: %s", pkValue, err.Error())
		dbErr := dberrors.NewDbError(dberrors.InternalServerError, err)
		return "", dbErr
	}

	return string(jsonBytes), nil
}

func (t *DbTable) Get(dbReqCtx *commtypes.DbReqContext) ([]byte, *dberrors.DbError) {

	t.mu.RLock()
	defer t.mu.RUnlock()

	// TODO: check all paths are handled.

	key, value, err := rest.ParseKeyValue(dbReqCtx.ResPath)
	if err != nil {
		return nil, dberrors.NewDbError(dberrors.InvalidInput, err)
	}

	//NOTE: there is no empty index 'indexName:' access at the moment.
	if value == "" {
		return nil, dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, errors.New("key is not provided"))
	}

	if key == "" {
		//primary key or function
		if value[0] != '$' {
			row, dberr := t.GetRowByPrimaryKeyReturnsJSON(value)
			if dberr != nil {
				return nil, dberr
			}

			return []byte(row), nil
		} else {
			if isFunction(value) {
				fnName, err := extractFunctionName(value)
				if err != nil {
					return nil, dberrors.NewDbError(dberrors.InvalidInput, err)
				}
				return callTableFunction(t, dbReqCtx, fnName)
			}

			//TODO: add support for $row (keywords)
		}

	} else {
		//

	}

	//TODO: check if nil, nil return is ok?
	return nil, nil
}
