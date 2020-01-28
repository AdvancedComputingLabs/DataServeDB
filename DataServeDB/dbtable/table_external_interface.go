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

package dbtable

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	storage "DataServeDB/dbsystem/dbstorage"
	"DataServeDB/dbsystem/systypes/dtIso8601Utc"
	"DataServeDB/dbsystem/systypes/guid"
)

//TODO: move it to error messages (single location)

type TableRow map[string]interface{} //it is by field name.

type createTableExternalStruct struct {
	TableName   string
	TableFields []string
}

type DbTable struct {
	tblMain              *tableMain //table structure information to keep it separate from data, so data disk io can be done separately.
	tblData              *tableDataContainer
	createTableStructure createTableExternalStruct
}

func CreateTableJSON(jsonStr string) (*DbTable, error) {

	var createTableData createTableExternalStruct
	if err := json.Unmarshal([]byte(jsonStr), &createTableData); err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return nil, errors.New("error found in table creation json")
	}

	tdc := &tableDataContainer{
		Rows:          nil,
		PkToRowMapper: map[interface{}]int64{},
	}

	//TODO: use globabl const, better practice
	return createTable(-1, &createTableData, tdc)
}

func createTable(tableInternalId int, createTableData *createTableExternalStruct, tblDataContainer *tableDataContainer) (*DbTable, error) {
	// I think it better belongs here than table.go as it is creating DbTable

	tbl := DbTable{}

	if tblMain, err := validateCreateTableMetaData(tableInternalId, createTableData); err != nil {
		return nil, err
	} else {
		tbl.tblMain = tblMain

	}

	tbl.tblData = tblDataContainer

	tbl.createTableStructure = *createTableData

	return &tbl, nil
}

func (t *DbTable) DebugPrintInternalFieldsNameMappings() string {
	return fmt.Sprintf("%#v", t.tblMain.TableFieldsMetaData.FieldNameToFieldInternalId)
}

func (t *DbTable) InsertRowJSON(jsonStr string) error {

	var rowDataUnmarshalled TableRow
	if e := json.Unmarshal([]byte(jsonStr), &rowDataUnmarshalled); e != nil {
		_ = e
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return errors.New("error occured in parsing row json")
	}
	rowProperTyped, rowInternalIds, e := validateRowData(t.tblMain, rowDataUnmarshalled)
	if e != nil {
		return e
	}
	// check the duplicate primary key before insert
	if _, ok := t.tblData.PkToRowMapper[rowInternalIds[0]]; ok {
		return errors.New("Duplicate Found for Primary Key")
	}
	_ = rowProperTyped // reminds me why it is used.

	{
		numOfRows := int64(len(t.tblData.Rows))
		t.tblData.Rows = append(t.tblData.Rows, rowInternalIds)
		t.tblData.PkToRowMapper[rowInternalIds[0]] = numOfRows // TODO: should get pk or other secondary keys here properly
	}
	/* METOD gob ENCODING */
	gob.Register(dtIso8601Utc.Iso8601Utc{})
	gob.Register(guid.Guid{})
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(t.tblData)
	if err != nil {
		println("error ")
		log.Fatal("encode error:", err)
	}
	/*  METHOD JSON ENCODING
	tb, err := json.Marshal(t.tblData.Rows)
	if err != nil {
		return err
	}
	storage.SaveToTable(t.tblMain.TableId, tb)
	*/
	storage.SaveToTable(t.tblMain.TableId, network.Bytes())

	return nil
}
func (t *DbTable) GetLength() int {
	return int(int64(len(t.tblData.Rows)))
}

func (t *DbTable) GetRowByPrimaryKey(pkValue interface{}) (TableRow, error) {
	dbType, dbTypeProps := t.tblMain.getPkType()

	pkValueCasted, e := dbType.ConvertValue(pkValue, dbTypeProps)
	if e != nil {
		return nil, e
	}

	rowNum, exists := t.tblData.PkToRowMapper[pkValueCasted]
	if !exists {
		return nil, fmt.Errorf("value '%v' not found", pkValue)
	}

	row, e := toLabeledByFieldNames(t.tblData.Rows[rowNum], t.tblMain)
	if e != nil {
		return nil, e
	}

	return row, nil
}

func (t *DbTable) GetRowByPrimaryKeyReturnsJSON(pkValue interface{}) (string, error) {
	// fmt.Printf("here %v\n", t.tblData.Rows)
	row, e := t.GetRowByPrimaryKey(pkValue)
	if e != nil {
		return "", e
	}

	jsonBytes, e := json.Marshal(row)
	if e != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return "", fmt.Errorf("error occured while converting row data to json; primary key value '%v'", pkValue)
	}

	return string(jsonBytes), nil
}
