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
	"log"

	"DataServeDB/dbsystem/dbstorage"
)

// Description: dbtable package saving and loading file.

// !WARNING: unstable api, but this needs to be in dbtable package.

type DbTableRecreation struct {
	TableInternalId   int
	CreationStructure createTableExternalStruct
}

//!INFO: this is just for demonstration, tune it if you think there is better way to do it.
func GetSaveLoadStructure(dbtbl *DbTable) (string, error) {
	slStruct := DbTableRecreation{}

	slStruct.TableInternalId = dbtbl.TblMain.TableId
	slStruct.CreationStructure = dbtbl.createTableStructure

	//TODO: handle error
	if r, e := json.Marshal(slStruct); e == nil {
		return string(r), nil
	}

	return "", errors.New("did not convert to json")
}

//TODO: whats the difference betwen LoadTableFromDB and LoadFromJson?
// Seems like only dbname difference. Then, can be reduced to one function.

//TODO: change function name
//TODO: dbname should be part of it?
func LoadTableFromDB(dbtblJson, dbName string) (*DbTable, error) {
	var slStruct DbTableRecreation

	if e := json.Unmarshal([]byte(dbtblJson), &slStruct); e != nil {
		//TODO: for later after version 0.5, return structured error, top error json error and in sub structure include the json message.
		return nil, e
	}

	//TODO: path shouldn't be here
	path := fmt.Sprintf("../../data/%s/table/%s.json", dbName, slStruct.CreationStructure.TableName)
	row, err := dbstorage.LoadTableFromPath(path)
	if err != nil {
		return nil, err
	}

	var tblData tableDataContainer

	//TODO: can refectored to its own util function
	buf := bytes.NewReader(row)
	dec := gob.NewDecoder(buf)

	//gob.Register(dtIso8601Utc.Iso8601Utc{})
	//gob.Register(guid.Guid{})

	err = dec.Decode(&tblData)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	//tblData := &tableDataContainer{
	//	Rows:          rowDataUnmarshalled.Rows,
	//	PkToRowMapper: rowDataUnmarshalled.PkToRowMapper,
	//}
	return createTable(slStruct.TableInternalId, &slStruct.CreationStructure, &tblData)
}

func LoadFromJson(dbtblJson string) (*DbTable, error) {
	var slStruct DbTableRecreation

	if e := json.Unmarshal([]byte(dbtblJson), &slStruct); e != nil {
		//TODO: for later after version 0.5, return structured error, top error json error and in sub structure include the json message.
		return nil, e
	}

	row, err := dbstorage.LoadTableFromDisk(slStruct.TableInternalId)
	if err != nil {
		return nil, err
	}

	var tblData tableDataContainer

	//TODO: can refectored to its own util function
	buf := bytes.NewReader(row)
	dec := gob.NewDecoder(buf)

	//gob.Register(dtIso8601Utc.Iso8601Utc{})
	//gob.Register(guid.Guid{})

	err = dec.Decode(&tblData)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	return createTable(slStruct.TableInternalId, &slStruct.CreationStructure, &tblData)
}
