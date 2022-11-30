package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"DataServeDB/dbtable"
	"DataServeDB/paths"
	"DataServeDB/unstable_api/db"
	"DataServeDB/unstable_api/runtime"
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

// Storage Cases:
// 	1) Save storage.
// 	2) Load storage after restart, if it is persistant storage.

//Note: did the tests this way because testing needs to chained, maybe there is better built in way to do this.

type row struct {
	Id       int
	UserName string
}

func TestPaths(t *testing.T) {
	fmt.Println(paths.GetDatabasesMainDirPath())
	fmt.Println(paths.GetWorkingDirPath())
	fmt.Println(paths.GetExeDirPath())
}

func TestStarter(t *testing.T) {

	//testCreateTableJSON(t)

	testRestApiGet(t)

}

func TestSaveTableMetadata(t *testing.T) {
	/*
		Problems:
		1) data attaching needs to be figured out.
		2) counter starts again, counter state needs to be saved.
	*/

	createTableJSON := `{
	  "TableName": "Tbl01",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable",
		"DateAdded datetime default:Now() !Nullable",
		"GlobalId guid default:NewGuid() !Nullable"
	  ]
	}`

	// db.CreateTableJSON

	db, e := runtime.GetDb("re_db")
	if e != nil {
		t.Errorf("%v\n", e)
		return
	}

	/* test for save structure
	if jsonStr, err := dbtable.GetTableStorageStructureJson(tbl01); err == nil {
			fmt.Println(jsonStr)
			//testLoadTableMetadata(jsonStr, t)
		} else {
			t.Errorf("%v\n", err)
		}
	*/

	if err := db.CreateTableJSON(createTableJSON, nil); err == nil {

	} else {
		t.Errorf("%v\n", err)
	}
}

/*func testLoadTableMetadata(jsonStr string, t *testing.T) {
	if tbl, err := dbtable.LoadFromJson(jsonStr); err == nil {
		fmt.Printf("table loaded :- %v\n", tbl)
		if false {
			testInsertRowJSON(tbl, t)
		} else {
			testByPk(tbl, t)
		}
	} else {
		t.Errorf("%v\n", err)
	}
}*/

func TestCreateTableJSON(t *testing.T) {

	// TODO: Remove. This test is in table_external_storage_test.go

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var reTableNameExists = regexp.MustCompile(`(?m)table name '\w+' already exits`)

	//"DateAdded datetime default:Now() !Nullable"
	// "DateAdded !Nullable"; insert datetime

	createTableJSON := `{
	  "TableName": "Tbl01",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable",
		"DateAdded datetime default:Now() !Nullable",
		"GlobalId guid default:NewGuid() !Nullable"
	  ]
	}`

	//TODO: check why it is not returning counter.

	db, e := runtime.GetDb("re_db")
	if e != nil {
		t.Errorf("%v\n", e)
		return
	}

	log.Println("creating table...")
	if dberr := db.CreateTableJSON(createTableJSON, nil); dberr == nil {
		//fmt.Println(tbl01.DebugPrintInternalFieldsNameMappings())
		log.Println("table creation successful")
		log.Println("calling testInsertRowJSON(*)...")
		testInsertRowJSON(db, t)
	} else {
		if reTableNameExists.MatchString(dberr.ToError().Error()) {
			log.Println(dberr.ToError().Error())
			log.Println("calling testInsertRowJSON(*)...")
			testInsertRowJSON(db, t)
		} else {
			t.Errorf("%v\n", dberr)
		}
	}
}
func TestDeleteTableJSON(t *testing.T) {
	db, e := runtime.GetDb("re_db")
	if e != nil {
		t.Errorf("%v\n", e)
		return
	}

	dberr := db.DeleteTable("Tbl01")
	if dberr != nil {
		log.Fatal(dberr)
	}
	_, err := db.GetTable("Tbl01")
	if err == nil {
		log.Fatal("table should not exist")
	} else {
		log.Printf("OK. error message: %v.", err)
	}

}

func testRestApiGet(t *testing.T) {

	// TODO: Remove. This test is in table_rest_api_test.go

	// http://localhost:8080/re_db/tables/Id:1

	for !runtime.IsInitalized() {
		fmt.Println("Is initialized: ", runtime.IsInitalized())
		time.Sleep(time.Second * 1)
	}

	fmt.Println("#1")
	time.Sleep(time.Second * 5)
	fmt.Println("#2")
}

func testByPk(tbl *dbtable.DbTable, t *testing.T) {

	// TODO: Remove. This test is in table_external_storage_test.go

	for i := 0; i < tbl.GetLength(); i++ {
		testGetRowByPk(tbl, t, i)
	}
}

func testInsertRowJSON(db *db.DB, t *testing.T) {

	// TODO: Remove. This test is in table_external_storage_test.go

	tbl, e := db.GetTable("Tbl01")
	if e != nil {
		t.Error(e)
		return
	}

	items := [4]string{"captain america", "IRO MAN", "professor HULk", "peter Parker"}
	length := tbl.GetLength()

	for i, item := range items {
		row01 := row{
			Id:       i + length,
			UserName: item,
		}

		row01Json, err := json.Marshal(row01)
		if err != nil {
			t.Error("error converting to json")
		} else {
			log.Println("inserting row...")
			if e := tbl.InsertRowJSON(string(row01Json)); e == nil {
				log.Println("Insert Test Successful")
				testGetRowByPk(tbl, t, row01.Id)
			} else {

				t.Errorf("for row id: %d; error: %v\n", row01.Id, e)
			}
		}
	}

	testDelRowByPk(tbl, t, 5)
}

func testGetRowByPk(tbl *dbtable.DbTable, t *testing.T, i int) {

	// TODO: Remove. This test is in table_external_storage_test.go

	//TODO: test with multiple types to make sure get function is working properly with different types.

	if row, e := tbl.GetRowByPrimaryKey(i); e == nil {
		//fmt.Println(row)
		_ = row // don't need to print this as rowJson is print it.
		if rowJson, e := tbl.GetRowByPrimaryKeyReturnsJSON(i); e == nil {
			fmt.Println("Get With Pk Test Successful")
			fmt.Println(rowJson)
		} else {
			t.Errorf("%v\n", e)
		}
	} else {
		t.Errorf("%v\n", e)
	}
}

func testDelRowByPk(tbl *dbtable.DbTable, t *testing.T, i int) {

	// TODO: Remove. This test is in table_external_storage_test.go

	err := tbl.DeleteRowByPrimaryKey(i)
	if err != nil {
		t.Errorf("error@delete: %v\n", err)
	}
}
