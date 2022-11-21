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

package db

import (
	"DataServeDB/utils/rest/dberrors"
	"fmt"
	"log"

	"DataServeDB/dbtable"
)

/*
	table ops:
	* create table
		- data object is created?
		- data file created?

	* recreate from structure
		- empty table, data object?
		- data file is there?

	* insert row

	* update row
	* delete row
	* change table
	* delete table
*/

func (d *DB) GetTable(tableName string) (*dbtable.DbTable, error) {
	tblNameCaseHandled := syscasing(tableName)

	_, tbl, e := d.Tables.GetByNameUnsync(tblNameCaseHandled)
	if e != nil {
		//TODO: handle errors properly.
		return nil, fmt.Errorf("table '%s' does not exit", tableName)
	}

	//var tbl *dbtable.DbTable
	//var ok bool
	//
	//if tbl, ok = tblI.(*dbtable.DbTable); !ok {
	//	//TODO: maybe there is better way to do this
	//	log.Fatalf("Casting error while getting table '%s'. This shouldn'd happen there is error in the code.\n", tableName)
	//}

	return tbl, nil
}

type CreateTableCallback = func(table *dbtable.DbTable) error

func (d *DB) CreateTableJSON(jsonStr string, callback CreateTableCallback) *dberrors.DbError {
	//TODO: multi thread syncing
	//TODO: log errors
	//TODO: remove callback function parameter?

	//create
	//does not save table to disk
	//TODO: perf optimization, check if table name has conflict.
	tbl, dberr := dbtable.CreateTableJSON(jsonStr, d)
	if dberr != nil {
		return dberr
	}

	tblNameCaseHandled := syscasing(tbl.TblMain.TableName)

	if tbl.TblMain.TableId < 0 {

		tbl.TblMain.TableId = d.Tables.GetLastIdUnsync()
		triesCount := 0

		//try 5 times then exit.
		for true {
			triesCount++
			tbl.TblMain.TableId++

			//TODO: make number of tries configurable?
			if triesCount > 5 {
				if dberr != nil {
					return dberr
				}
				//TODO: test this
				panic("CODE_ERROR: Unable to find free table id, should have returned error before this.")
			}

			if dberr = d.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tbl.TblMain.TableName, tbl); dberr == nil {
				break
			} else {
				if dberr.ErrCode != dberrors.TableIdAlreadyExists {
					return dberr
				}
			}
		}

	} else {
		//NOTE: this is create operation so shouldn'd have tableid > -1 but I kept it here just in case.
		// Probably better to check the logic later and remove it.
		if dberr = d.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tbl.TblMain.TableName, tbl); dberr != nil {
			//TODO: it is probably load operation, conflict normally should not happen. If it does then it is probably error in the code/logic of the code.
			// Check if create table supplies table id. Shouldn't but there might be some edge cases. Check later.
			// Identify the difference in operation for the log. -- HY 12-Oct-2022
			return dberr
		}
	}

	dberr = tbl.EventAfterTableIdAssignment()
	if dberr != nil {
		goto ERROR_EXIT
	}

	//if callback != nil {
	//	dberr = callback(tbl)
	//	if dberr != nil {
	//		goto ERROR_EXIT
	//	}
	//}

ERROR_EXIT:
	if dberr != nil {
		//reverse table operations
		if _, _, err2 := d.Tables.RemoveByNameUnsync(tblNameCaseHandled); err2 != nil {
			log.Fatalln("CRITICAL_ERROR: Unable to remove table during callback error; error: ", err2)
		}
		return dberr
	}

	//save to disk
	//NOTE: table creation only creates metadata which is stored in db metadata
	dberr = d.updateDbMetadataOnDisk()
	if dberr != nil {
		return dberr
	}

	return nil
}

func (d *DB) DeleteTable(tableName string) *dberrors.DbError {

	tblNameCaseHandled := syscasing(tableName)

	// IMP-NOTE: data files are not deleted here. They are supposed to be handled by the storage(s) themselves.

	_, tbl, err := d.Tables.RemoveByNameUnsync(tblNameCaseHandled)
	if err != nil {
		// TODO: log internal error for debugging? 'err' is not used here.
		return dberrors.NewDbError(dberrors.TableNotFound, fmt.Errorf("table '%s' does not exit", tableName))
	}

	//probably need later to clean up dat file.
	if dberr := dbtable.DeleteTable(tbl); dberr != nil {
		// add table back to the list
		if dberr2 := d.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tableName, tbl); dberr2 != nil {
			//TODO: handle this better. Perhaps do not remove table first.
			log.Fatalln("CRITICAL_ERROR: Unable to add table back to the list during table delete error; error: ", dberr2)
		}
		return dberr
	}

	//update metadata on disk
	//IMP-NOTE: This does not remove data files. See imp-note above.
	if dberr := d.updateDbMetadataOnDisk(); dberr != nil {
		return dberr
	}

	return nil
}

//Not sure this is needed, but kept it to review later. HY -- 20-Mar-2022
/*func (t *DB) DeleteTableJSON(jsonStr string) error {

	return nil
}*/
