// Copyright (c) 2020 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package dbtable

import (
	"DataServeDB/utils/rest/dberrors"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"DataServeDB/comminterfaces"
	"DataServeDB/paths"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
)

// Description: dbtable package saving and loading file.

// !WARNING: unstable api, but this needs to be in dbtable package.

// TODO: const file to put all the consts in one place?
const tablesDataPathRelative = "tables_data"

type DatabaseTableRecreation struct {
	TableInternalId          int //not runtime, used to save table data.
	DbTableCreationStructure createTableExternalStruct
}

// GetTableStorageStructureJson !INFO: this is just for demonstration, tune it if you think there is a better way to do it.
func GetTableStorageStructureJson(dbtbl *DbTable) (string, error) {
	slStruct := GetTableStorageStructure(dbtbl)

	//TODO: handle error
	if r, e := json.Marshal(slStruct); e == nil {
		return string(r), nil
	}

	return "", errors.New("did not convert to json")
}

func GetTableStorageStructure(dbtbl *DbTable) DatabaseTableRecreation {
	slStruct := DatabaseTableRecreation{}
	slStruct.TableInternalId = dbtbl.TblMain.TableId
	slStruct.DbTableCreationStructure = dbtbl.createTableStructure
	return slStruct
}

func LoadFromJson(dbtblJson string, dbPtr comminterfaces.DbPtrI) (*DbTable, error) {
	var slStruct DatabaseTableRecreation

	if e := json.Unmarshal([]byte(dbtblJson), &slStruct); e != nil {
		//TODO: for later after version 0.5, return structured error, top error json error and in sub structure include the json message.
		return nil, e
	}

	slStruct.DbTableCreationStructure._dbPtr = dbPtr

	return LoadFromTableSaveStructure(slStruct)
}

func LoadFromTableSaveStructure(slStruct DatabaseTableRecreation) (*DbTable, error) {

	dbTable, dberr := createTable(slStruct.TableInternalId, &slStruct.DbTableCreationStructure /*, tblData*/)
	if dberr != nil {
		return dbTable, dberr.ToError()
	}

	if dberr := applyStorages(dbTable, false); dberr != nil {
		return dbTable, dberr.ToError()
	}

	return dbTable, nil
}

func (t *DbTable) GetTableInfo() (*idbstorer.TableInfo, error) {
	// TODO: use of PkPos is best? Or there is better method to get the primary key column id?
	return &idbstorer.TableInfo{TableName: t.TblMain.TableName, PrimaryKeyColId: t.TblMain.PkPos}, nil
}

func applyStorages(t *DbTable, newCreation bool) *dberrors.DbError {

	//check: table is valid
	if t.TblMain.TableId < 0 {
		panic("table id cannot be negative, something wrong in code")
	}

	//table folder path
	//tableName := t.TblMain.TableName
	tableFolderName := fmt.Sprintf("table_%d", t.TblMain.TableId)
	tableFolderPath := paths.Combine(t.createTableStructure._dbPtr.DbPath(), tablesDataPathRelative, tableFolderName)

	//TODO: check when table storers are mention in creation json.
	//TODO: check when table structure is loaded from disk.

	storersNames := strings.Fields(t.createTableStructure.TableStorages)

	if len(storersNames) == 0 {
		//NOTE: default case when no storages are named at creation of the table.
		storersNames = append(storersNames, "StorerMemV1")
		storersNames = append(storersNames, "StorerDiskV1")
	}

	var stores []idbstorer.StorerBasic

	for _, storerKey := range storersNames {
		if storeCreator := idbstorer.GetStoreBasic(storerKey); storeCreator != nil {
			store, dberr := storeCreator(t.TblMain.TableId, tableFolderPath, t)
			if dberr != nil {
				return dberr
			}
			stores = append(stores, store)
		} else {
			//TODO: panic is correct here?
			panic(fmt.Sprintf("storage '%s' does not exist", storerKey))
		}
	}

	if !newCreation {
		//need recursion
		loadStorages(0, stores)
	}

	for i, store := range stores {

		if newCreation {
			var prevStore, nextStore idbstorer.StorerBasic

			if i > 0 {
				prevStore = stores[i-1]
			}

			if (i + 1) < len(stores) {
				nextStore = stores[i+1]
			}

			dberr := store.CreateInit(prevStore, nextStore)
			if dberr != nil {
				return dberr
			}
		} // newCreation

		// for load added here to keep addition of stores in correct order. recursion reverses the order
		t.createTableStructure._tableStorersInstances = append(t.createTableStructure._tableStorersInstances, store)
	}

	return nil
}

func loadStorages(current_idx int, stores []idbstorer.StorerBasic) *idbstorer.StorerLoadResult {
	//TODO: what if something needs to be passed on to next store?

	if current_idx >= len(stores) {
		return nil
	}

	var prevStore, nextStore idbstorer.StorerBasic

	if current_idx > 0 {
		prevStore = stores[current_idx-1]
	}

	store := stores[current_idx]

	if (current_idx + 1) < len(stores) {
		nextStore = stores[current_idx+1]
	}

	var nextStoreResult *idbstorer.StorerLoadResult = nil
	var input any = nil

	if (current_idx + 1) < len(stores) {
		nextStoreResult = loadStorages(current_idx+1, stores)
	}

	if nextStoreResult != nil {
		if nextStoreResult.Err != nil {
			return nextStoreResult
		}

		input = nextStoreResult.Data
	}

	result := store.LoadTable(input, prevStore, nextStore)
	if result != nil {
		if result.Err != nil {
			return result
		}
	}

	return result
}
