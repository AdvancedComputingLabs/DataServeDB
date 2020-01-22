package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"DataServeDB/dbtable"
)

//TODO: table id

//TODO: test, bad table name. Should fail.
//TODO: test, bad field name. Should fail.
//TODO: test, missing primary key. Should fail.
//TODO: test, more than one primary keys (at the moment only one pk is supported). Should fail.
//TODO: test, bad dbtype. Should fail.
//TODO: test, bad field property. Should fail.

// TODO: datetime cases:
//		1) user inserts datatime, it must be Iso8601Utc format.
//  	2) datetime is inserted automatically, server will insert in Iso8601Utc by default.
//		3) cases where datetime is not in Iso8601Utc? Currently, imo if user inserted datatime is not in Iso8601Utc format
//			it should fail.

//Note: did the tests this way because testing needs to chained, maybe there is better built in way to do this.

type row struct {
	Id       int
	UserName string
}

func TestStarter(t *testing.T) {
	testCreateTableJSON(t)
}

func TestSaveTableMetadata(t *testing.T) {
	/*
		Problems:
		1) data attaching needs to be figured out.
		2) counter starts again, counter state needs to be saved.
	*/

	createTableJSON := `{
	  "TableName": "Tbl01",
	  "TableFields": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable"
	  ]
	}`

	if tbl01, err := dbtable.CreateTableJSON(createTableJSON); err == nil {
		if jsonStr, err := dbtable.GetSaveLoadStructure(tbl01); err == nil {
			println(jsonStr)
			testLoadTableMetadata(tbl01, t)
		} else {
			t.Errorf("%v\n", err)
		}
	} else {
		t.Errorf("%v\n", err)
	}
}

func testLoadTableMetadata(dbtbl *dbtable.DbTable, t *testing.T) {
	if tbl, err := dbtable.LoadFromJson(dbtbl); err == nil {
		fmt.Printf("table loaded :- %v\n", tbl)
		// testInsertRowJSON(tbl, t)
		testByPk(tbl, t)
	} else {
		t.Errorf("%v\n", err)
	}
}

func testCreateTableJSON(t *testing.T) {

	//"DateAdded datetime default:Now() !Nullable"
	// "DateAdded !Nullable"; insert datetime

	createTableJSON := `{
	  "TableName": "Tbl01",
	  "TableFields": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable"
	  ]
	}`

	//TODO: check why it is not returning counter.

	if tbl01, err := dbtable.CreateTableJSON(createTableJSON); err == nil {
		fmt.Println(tbl01.DebugPrintInternalFieldsNameMappings())
		testInsertRowJSON(tbl01, t)
	} else {
		t.Errorf("%v\n", err)
	}
}
func testByPk(tbl *dbtable.DbTable, t *testing.T) {
	for i := 0; i < 4; i++ {
		testGetRowByPk(tbl, t, i)
	}
}

func testInsertRowJSON(tbl *dbtable.DbTable, t *testing.T) {
	items := [4]string{"captain Marvel", "wanda maximoff", "harry poter", "peter Parker"}

	for i, item := range items {
		row01 := row{
			Id:       i,
			UserName: item,
		}
		row01Json, err := json.Marshal(row01)
		if err != nil {
			t.Error("erroe converting")
		} else {
			if e := tbl.InsertRowJSON(string(row01Json)); e == nil {
				fmt.Println("Insert Test Successful")
				testGetRowByPk(tbl, t, i)
			} else {
				t.Errorf("%v\n", e)
			}
		}
	}
}

func testGetRowByPk(tbl *dbtable.DbTable, t *testing.T, i int) {
	//TODO: test with multiple types to make sure get function is working properly with different types.

	if row, e := tbl.GetRowByPrimaryKey(i); e == nil {
		fmt.Println(row)
		if rowJson, e := tbl.GetRowByPrimaryKeyReturnsJSON(i); e == nil {
			fmt.Println(rowJson)
		} else {
			t.Errorf("%v\n", e)
		}
	} else {
		t.Errorf("%v\n", e)
	}
}
