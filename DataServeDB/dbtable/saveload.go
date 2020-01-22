package dbtable

import (
	storage "DataServeDB/dbsystem/dbstorage"
	"encoding/json"
	"errors"
	"fmt"
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

	slStruct.TableInternalId = dbtbl.tblMain.TableId
	slStruct.CreationStructure = dbtbl.createTableStructure

	//TODO: handle error
	if r, e := json.Marshal(slStruct); e == nil {
		return string(r), nil
	}

	return "", errors.New("did not convert to json")
}

func LoadFromJson(dbtbl *DbTable) (*DbTable, error) {
	var slStruct DbTableRecreation
	row, err := storage.LoadTable(slStruct.TableInternalId)
	if err != nil {
		return nil, err
	}

	var rowDataUnmarshalled []map[int]interface{}
	if e := json.Unmarshal(row, &rowDataUnmarshalled); e != nil {
		_ = e
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return nil, e
	}

	fmt.Printf("table --> %v\n", dbtbl)
	for _, rowData := range rowDataUnmarshalled {
		// fmt.Printf("table --> %t\n", rowData)
		var row = TableRow{}
		for i, data := range rowData {
			row[dbtbl.tblMain.TableFieldsMetaData.FieldInternalIdToFieldMetaData[i].FieldName] = data
		}
		_, rowInternalIds, e := validateRowData(dbtbl.tblMain, row)
		if e != nil {
			println(e.Error())
			return nil, e
		}

		numOfRows := int64(len(dbtbl.tblData.Rows))
		dbtbl.tblData.Rows = append(dbtbl.tblData.Rows, rowInternalIds)
		dbtbl.tblData.PkToRowMapper[rowInternalIds[0]] = numOfRows
	}

	// fmt.Printf("Rows ->> %v\n", dbtbl.tblData.Rows)
	// fmt.Printf("index ->> %t\n", dbtbl.tblData.PkToRowMapper)
	// dbtbl, e := createTable(slStruct.TableInternalId, &slStruct.CreationStructure, tdc)
	// if e != nil {
	// 	return nil, e
	// }

	return dbtbl, nil
}
