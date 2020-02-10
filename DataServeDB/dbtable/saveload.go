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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"DataServeDB/comminterfaces"
	"DataServeDB/dbsystem/dbstorage"
	"DataServeDB/paths"
)

// Description: dbtable package saving and loading file.

// !WARNING: unstable api, but this needs to be in dbtable package.

//TODO: const file to put all the consts in one place?
const tablesDataPathRelative = "tables_data"

type DatabaseTableRecreation struct {
	TableInternalId          int //not runtime, used to save table data.
	DbTableCreationStructure createTableExternalStruct
}

//!INFO: this is just for demonstration, tune it if you think there is better way to do it.
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

	tblData, err2 := loadTableDataFromDisk(&slStruct)
	if err2 != nil {
		return nil, err2
	}

	return createTable(slStruct.TableInternalId, &slStruct.DbTableCreationStructure, tblData)
}

func loadTableDataFromDisk(slStruct *DatabaseTableRecreation) (*tableDataContainer, error) {
	fileName := fmt.Sprintf("table_%d.dat", slStruct.TableInternalId)
	path := paths.Combine(slStruct.DbTableCreationStructure._dbPtr.DbPath(), tablesDataPathRelative, fileName)

	//NOTE: this could be empty if data is not saved.
	tblDataBytes, err := dbstorage.LoadFromDisk(path)
	if err != nil {
		//if path error then empty table
		if os.IsNotExist(err) {
			tdc := &tableDataContainer{
				Rows:          nil,
				PkToRowMapper: map[interface{}]int64{},
			}
			return tdc, nil
		}
		return nil, err
	}

	var tblData tableDataContainer

	//TODO: can refectored to its own util function?
	buf := bytes.NewReader(tblDataBytes)
	dec := gob.NewDecoder(buf)

	err = dec.Decode(&tblData)
	if err != nil {
		return nil, err
	}

	return &tblData, nil
}
