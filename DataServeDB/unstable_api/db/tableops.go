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
	"DataServeDB/dbtable"
	"errors"
	"fmt"
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

	tbl, exists := t.MapOfTables[tblNameCaseHandled]
	if !exists {
		return nil, fmt.Errorf("table '%s' does not exit", tableName)
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

	//add to map
	//TODO: name already exists check
	tblNameCaseHandled := syscasing(tbl.TblMain.TableName)
	if _, alreadyExists := t.MapOfTables[tblNameCaseHandled]; alreadyExists {
		return errors.New("table name already exits")
	}
	t.MapOfTables[tblNameCaseHandled] = tbl

	//save to disk
	//NOTE: table creation only creates metadata which is stored in db metadata
	e = t.updateDbMetadataOnDisk()
	if e != nil {
		return e
	}

	return nil
}
