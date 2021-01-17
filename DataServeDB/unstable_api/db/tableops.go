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

func (t *DB) GetTable(tableName string) (*dbtable.DbTable, error) {
	tblNameCaseHandled := syscasing(tableName)

	_, tblI, e := t.Tables.GetByNameUnsync(tblNameCaseHandled)
	if e != nil {
		//TODO: handle errors properly.
		return nil, fmt.Errorf("table '%s' does not exit", tableName)
	}

	var tbl *dbtable.DbTable
	var ok bool

	if tbl, ok = tblI.(*dbtable.DbTable); !ok {
		//TODO: maybe there is better way to do this
		log.Fatalf("Casting error while getting table '%s'. This shouldn't happen there is error in the code.\n", tableName)
	}

	return tbl, nil
}

func (t *DB) CreateTableJSON(jsonStr string) error {
	//TODO: multi thread syncing

	//create
	//does not save table to disk
	//TODO: perf optimization, check if table name has conflict.
	tbl, e := dbtable.CreateTableJSON(jsonStr, t)
	if e != nil {
		return e
	}

	tblNameCaseHandled := syscasing(tbl.TblMain.TableName)

	//if _, alreadyExists := t.Tables[tblNameCaseHandled]; alreadyExists {
	//	return errors.New("table name already exits")
	//}
	//t.Tables[tblNameCaseHandled] = tbl

	//TODO: table id management.

	if tbl.TblMain.TableId < 0 {

		tbl.TblMain.TableId = t.Tables.LastId
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

			if e = t.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tbl); e == nil {
				break
			}
		}

	} else {
		//NOTE: this is create operation so shouldn't have tableid > -1 but I kept it here just in case.
		// Probably better to check the logic later and remove it.
		if e = t.Tables.AddUnsync(tbl.TblMain.TableId, tblNameCaseHandled, tbl); e != nil {
			//TODO: properly handle and log errors.
			return errors.New("table name already exits")
		}
	}

	//save to disk
	//NOTE: table creation only creates metadata which is stored in db metadata
	e = t.updateDbMetadataOnDisk()
	if e != nil {
		return e
	}

	return nil
}
