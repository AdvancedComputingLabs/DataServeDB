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

package dbstorage

import (
	"fmt"
	"io/ioutil"
	"os"
)

var dbFile = "../../data/re_db/table/users.json"
var tablesPath = "dbsystem/dbstorage/tables/"

type dataType interface{}

// CreateTable to create a table
func CreateTable(tableName string) {
	// 	table := dbtable.NewTableMain(tableName)
	// 	saveToDisk(*table)
}

func SaveToDisk(data []byte) error {
	println(string(data))
	println(dbFile)
	db, err := os.OpenFile(dbFile, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		db, err = os.OpenFile(dbFile, os.O_APPEND, 0644)
	}
	defer db.Close()
	_, err = db.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// tableName or table Id should pass
// TODO :- change to suitable one, as of now it given as table name
func SaveToTable(tableID int, data []byte) error {
	println(tableID)
	// file := fmt.Sprintf("%stable%d.json", tablesPath, tableID)
	// println("file", file)
	db, err := os.OpenFile(dbFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		println(err)
		if !os.IsExist(err) {
			return err
		}
	}
	db.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// tableName or table Id should pass
// TODO :- change to suitable one, as of now it given as table name
func LoadTableFromDisk(tableID int) ([]byte, error) {

	file := fmt.Sprintf("%stable%d.json", tablesPath, tableID)
	return ioutil.ReadFile(file)
}
func LoadTableFromPath() ([]byte, error) {
	return ioutil.ReadFile(dbFile)
}

func LoadTableMeta() ([]byte, error) {
	return ioutil.ReadFile(dbFile)
}

// func getTableMeta(tableName string) ([]byte, error) {
// 	data, err := loadTableMeta()
// 	meta := json.Unmarshal()
// 	if err == nil {
// 		for _, tableMeta := range meta {
// 			if tableMeta.TableName == tableName {
// 				return tableMeta, nil
// 			}
// 		}
// 	}
// 	return json.Marshal(data)
// }
