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
	"errors"
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

func (me *DB) GetTable(tableName string) (*dbtable.DbTable, error) {
	tblNameCaseHandled := syscasing(tableName)

	_, tbl, e := me.Tables.GetByNameUnsync(tblNameCaseHandled)
	if e != nil {
		//TODO: handle errors properly.
		return nil, fmt.Errorf("table '%s' does not exit", tableName)
	}

	//var tbl *dbtable.DbTable
	//var ok bool
	//
	//if tbl, ok = tblI.(*dbtable.DbTable); !ok {
	//	//TODO: maybe there is better way to do this
	//	log.Fatalf("Casting error while getting table '%s'. This shouldn'me happen there is error in the code.\n", tableName)
	//}

	return tbl, nil
}

type CreateTableCallback = func(table *dbtable.DbTable) error

func (me *DB) CreateTableJSON(jsonStr string, callback CreateTableCallback) error {
	//TODO: multi thread syncing

	//create
	//does not save table to disk
	//TODO: perf optimization, check if table name has conflict.
	tbl, e := dbtable.CreateTableJSON(jsonStr, me)
	if e != nil {
		return e
	}

	tblNameCaseHandled := syscasing(tbl.TblMain.TableName)

	//if _, alreadyExists := me.Tables[tblNameCaseHandled]; alreadyExists {
	//	return errors.New("table name already exits")
	//}
	//me.Tables[tblNameCaseHandled] = tbl

	//TODO: table id management.

	if tbl.TblMain.TableId < 0 {

		tbl.TblMain.TableId = me.Tables.GetLastIdUnsync()
		triesCount := 0

		//try 5 times then exit.
		for true {
			triesCount++
			tbl.TblMain.TableId++

			//TODO: properly handle this, currently it assumes it is just id already exists.
			//TODO: need to check error type too.

			if triesCount > 5 {
				//TODO: handle this better.
				break
			}

			if e = me.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tbl); e == nil {
				break
			}
		}

	} else {
		//NOTE: this is create operation so shouldn'me have tableid > -1 but I kept it here just in case.
		// Probably better to check the logic later and remove it.
		if e = me.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tbl); e != nil {
			//TODO: properly handle and log errors.
			return errors.New("table name already exits")
		}
	}

	//TODO: handle error
	tbl.EventAfterTableIdAssignment()

	if callback != nil {
		e = callback(tbl)
		if e != nil {
			//reverse table operations
			if _, _, e = me.Tables.RemoveByNameUnsync(tblNameCaseHandled); e != nil {
				log.Fatalln("CRITICAL_ERROR: Unable to remove table during callback error; error: ", e)
			}
			return e
		}
	}

	//save to disk
	//NOTE: table creation only creates metadata which is stored in db metadata
	e = me.updateDbMetadataOnDisk()
	if e != nil {
		return e
	}

	return nil
}

func (me *DB) DeleteTable(tableName string) error {

	tblNameCaseHandled := syscasing(tableName)

	//TODO: check if table data needs deleting and dat file cleaning.

	_, tbl, e := me.Tables.RemoveByNameUnsync(tblNameCaseHandled)
	if e != nil {
		//TODO: handle errors properly.
		return fmt.Errorf("table '%s' does not exit", tableName)
	}

	//probably need later to clean up dat file.
	_ = tbl

	//update metadata on disk
	//NOTE: this does not remove data files
	e = me.updateDbMetadataOnDisk()
	if e != nil {
		return e
	}

	return nil
}

//Not sure this is needed, but kept it to review later. HY -- 20-Mar-2022
/*func (t *DB) DeleteTableJSON(jsonStr string) error {

	return nil
}*/
