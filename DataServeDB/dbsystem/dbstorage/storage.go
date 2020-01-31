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

//TODO: file paths here

var mataFile = "../../data/meta.json"
var tablesPath = "dbsystem/dbstorage/tables/"

type dataType interface{}

// CreateTable to create a table
// func CreateTable(tableName string) {
// 	// 	table := dbtable.NewTableMain(tableName)
// 	// 	saveToDisk(*table)
// }

func SaveToDisk(data []byte) error {
	println(string(data))
	db, err := os.OpenFile(mataFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func LoadMata() ([]byte, error) {
	return ioutil.ReadFile(mataFile)
}

// tableName or table Id should pass
// TODO :- change to suitable one, as of now it given as table name
func SaveToTable(tableRoot, tableName string, data []byte) error {
	file := fmt.Sprintf("../../data/%stable/%s.json", tableRoot, tableName)
	// println("file", file)
	db, err := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
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
func LoadTableFromPath(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// func LoadTableMeta() ([]byte, error) {
// 	return ioutil.ReadFile(dbFile)
// }

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
