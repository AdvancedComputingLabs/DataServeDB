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
	"log"

	storage "DataServeDB/dbsystem/dbstorage"
	"DataServeDB/dbsystem/systypes/dtIso8601Utc"
	"DataServeDB/dbsystem/systypes/guid"
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
func LoadTableFromDB(dbtblJson string) (*DbTable, error) {
	var slStruct DbTableRecreation

	if e := json.Unmarshal([]byte(dbtblJson), &slStruct); e != nil {
		//TODO: for later after version 0.5, return structured error, top error json error and in sub structure include the json message.
		return nil, e
	}
	// path := fmt.Sprintf("../../data/%s/table/%s.json", dbName, slStruct.CreationStructure.TableName)
	row, err := storage.LoadTableFromPath()
	if err != nil {
		return nil, err
	}
	network := bytes.NewReader(row) // Stand-in for a network connection
	dec := gob.NewDecoder(network)  // Will read from network.
	gob.Register(dtIso8601Utc.Iso8601Utc{})
	gob.Register(guid.Guid{})
	var rowDataUnmarshalled tableDataContainer
	err = dec.Decode(&rowDataUnmarshalled)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}
	tblData := &tableDataContainer{
		Rows:          rowDataUnmarshalled.Rows,
		PkToRowMapper: rowDataUnmarshalled.PkToRowMapper,
	}
	return createTable(slStruct.TableInternalId, &slStruct.CreationStructure, tblData)
}

func LoadFromJson(dbtblJson string) (*DbTable, error) {
	var slStruct DbTableRecreation

	if e := json.Unmarshal([]byte(dbtblJson), &slStruct); e != nil {
		//TODO: for later after version 0.5, return structured error, top error json error and in sub structure include the json message.
		return nil, e
	}

	row, err := storage.LoadTableFromDisk(slStruct.TableInternalId)
	if err != nil {
		return nil, err
	}

	/**** BY METHOD ON gob ENCODE  *****************/
	network := bytes.NewReader(row) // Stand-in for a network connection
	dec := gob.NewDecoder(network)  // Will read from network.
	gob.Register(dtIso8601Utc.Iso8601Utc{})
	gob.Register(guid.Guid{})
	var rowDataUnmarshalled tableDataContainer
	err = dec.Decode(&rowDataUnmarshalled)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}
	tblData := &tableDataContainer{
		Rows:          rowDataUnmarshalled.Rows,
		PkToRowMapper: rowDataUnmarshalled.PkToRowMapper,
	}
	/*
		// METHOD BY DOING JSON ENCODE ROWS AND VALIDATE PK
		// var rowDataUnmarshalled []map[int]interface{}
		// if e := json.Unmarshal(row, &rowDataUnmarshalled); e != nil {
		// 	_ = e
		// 	//log error for system auditing. This error logging message can be technical.
		// 	//TODO: make error result more user friendly.
		// 	return nil, e
		// }
		// fmt.Printf("table --> %v\n", dbtbl)
		// for _, rowData := range rowDataUnmarshalled {
		// 	// fmt.Printf("table --> %t\n", rowData)
		// 	var row = TableRow{}
		// 	for i, data := range rowData {
		// 		row[dbtbl.tblMain.TableFieldsMetaData.FieldInternalIdToFieldMetaData[i].FieldName] = data
		// 	}
		// 	_, rowInternalIds, e := validateRowData(dbtbl.tblMain, row)
		// 	if e != nil {
		// 		println(e.Error())
		// 		return nil, e
		// 	}

		// 	numOfRows := int64(len(dbtbl.tblData.Rows))
		// 	dbtbl.tblData.Rows = append(dbtbl.tblData.Rows, rowInternalIds)
		// 	dbtbl.tblData.PkToRowMapper[rowInternalIds[0]] = numOfRows
		// }

		// fmt.Printf("Rows ->> %v\n", dbtbl.tblData.Rows)
		// fmt.Printf("index ->> %t\n", dbtbl.tblData.PkToRowMapper)

	*/
	return createTable(slStruct.TableInternalId, &slStruct.CreationStructure, tblData)

	// if e != nil {
	// 	return nil, e
	// }

	// return dbtbl, nil
}
