package main

import (
	"fmt"
	"testing"

	"DataServeDB/dbtable"
)

//TODO: test, bad table name. Should fail.
//TODO: test, bad field name. Should fail.
//TODO: test, missing primary key. Should fail.
//TODO: test, more than one primary keys (at the moment only one pk is supported). Should fail.
//TODO: test, bad dbtype. Should fail.
//TODO: test, bad field property. Should fail.

//Note: did the tests this way because testing needs to chained, maybe there is better built in way to do this.

func TestStarter(t *testing.T) {
	testCreateTableJSON(t)
}

func testCreateTableJSON(t *testing.T) {

	createTableJSON := `{
  "tableName": "Tbl01",
  "PrimaryKeyName":"Id",
  "tableFields": [
    {
      "fieldName": "Id",
      "fieldType": "int32",
      "primaryKey": true
    },
    {
      "FieldName": "UserName",
      "FieldType": "string",
      "PrimaryKey": false
    }
  ]
}`

	if tbl01, err := dbtable.CreateTableJSON(createTableJSON); err == nil {
		fmt.Println(tbl01.DebugPrintInternalFieldsNameMappings())
		testInsertRowJSON(tbl01, t)
	} else {
		t.Errorf("%v\n", err)
	}

}

func testInsertRowJSON(tbl *dbtable.DbTable, t *testing.T) {

	row01Json := `{
    "Id" : 1,
    "UserName" : "JohnDoe"
}`

	if e := tbl.InsertRowJSON(row01Json); e == nil {
		fmt.Println("Insert Test Successful")
		testGetRowByPk(tbl, t)
	} else {
		t.Errorf("%v\n", e)
	}
}

func testGetRowByPk(tbl *dbtable.DbTable, t *testing.T) {
	//TODO: test with multiple types to make sure get function is working properly with different types.

	if row, e := tbl.GetRowByPrimaryKey(1); e == nil {
		fmt.Println(row)
		if rowJson, e := tbl.GetRowByPrimaryKeyReturnsJSON(1); e == nil {
			fmt.Println(rowJson)
		} else {
			t.Errorf("%v\n", e)
		}
	} else {
		t.Errorf("%v\n", e)
	}
}
