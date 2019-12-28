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
	"encoding/json"
	"errors"
	"fmt"
)

//TODO: move it to error messages (single location)
//TODO: code may need refactoring to better units of code and files

type TableRow map[string]interface{} //it is by field name.

//NOTE: fields need to be public (first letter capital) for json to marshal them.
type createTableExternalInterfaceFieldInfo struct {
	FieldName string
	FieldType string
}

// This it to convert from json.
// Convert from map too?
type createTableExternalInterface struct {
	//TODO: add field properties like incremental for number types and comparision function for string, etc...

	TableName      string
	PrimaryKeyName string
	TableFields    []createTableExternalInterfaceFieldInfo
}

type DbTable struct {
	tblMain *tableMain //table structure information to keep it separate from data, so data disk io can be done separately.
	tblData *tableDataContainer
}

func CreateTableJSON(jsonStr string) (*DbTable, error) {

	var createTableData createTableExternalInterface
	if err := json.Unmarshal([]byte(jsonStr), &createTableData); err != nil {
		//log error for system auditing. This error logging message can be technical.
		//TODO: make error result more user friendly.
		return nil, errors.New("error found in table creation json")
	}

	tbl := DbTable{}

	if tblMain, err := validateCreateTableMetaData(&createTableData); err != nil {
		return nil, err
	} else {
		tbl.tblMain = tblMain
	}

	tbl.tblData = &tableDataContainer{
		Rows:          nil,
		PkToRowMapper: map[interface{}]int64{},
	}

	return &tbl, nil
}

func (t *DbTable) DebugPrintInternalFieldsNameMappings() string {
	return fmt.Sprintf("%#v", t.tblMain.TableFieldsMetaData.fieldNameToFieldInternalId)
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

	_ = rowProperTyped // reminds me why it is used.

	{
		numOfRows := int64(len(t.tblData.Rows))
		t.tblData.Rows = append(t.tblData.Rows, rowInternalIds)
		t.tblData.PkToRowMapper[rowInternalIds[0]] = numOfRows // TODO: should get pk or other secondary keys here properly
	}

	return nil
}

func (t *DbTable) GetRowByPrimaryKey(pkValue interface{}) (TableRow, error) {

	pkValueCasted, e := t.tblMain.getPkType().ConvertValue(pkValue, true)
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
