package main

import (
	"errors"
	"log"
	"testing"

	"DataServeDB/storers/dbtable_disk_store_v1"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/unstable_api/db"
	"DataServeDB/unstable_api/runtime"
)

/*
	Descripion: Test basic storage operations. Create table, insert rows, check rows, delete table,
				saving and loading of data is working correctly.
                Does not go into detailed tests of database operations and functionality.

	            TODO: Needs multiple tests to run to test restart of server. Can be done in automated way?
                      Check exit of runtime after each test run.

	Plan:
		- Create table. [done: first run]
		- Insert few rows. [done: first run]
		- Save table. [done: first run]
		- Restart server. [done: second run]
		- Load table. [done: second run]
		- Randomly check few entries to check if table is loaded correctly. [done: second run]
		- Delete and update few rows.
		- Restart server.
		- Randomly check few entries to check if table is loaded correctly.
		- Delete table.
		- Restart server.
		- Check if table is deleted.
*/

// - Check storage is saved and loaded correctly after restart.
// - Check storage constraints.
// -- Constraints:
// --- Must not be first storage.
// --- Must be last storage.
// --- Requires another storage before it.

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// ## First Level Test Functions

func TestDeleteTable(t *testing.T) {

	// TODO: at the moment, data of the storage is not removed. Should it be removed?

	// server is already initialized, so we can use it.
	// load database through runtime api.
	db_, err := runtime.GetDb("re_db")
	if err != nil {
		log.Fatal(err)
	}

	// delete table.
	err = deleteTable(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check if table exists.
	_, err = db_.GetTable("TestTable02")
	if err == nil {
		log.Fatal("table should not exist")
	} else {
		log.Printf("OK. error message: %v.", err)
	}

}

// TestFirstRun tests first run after creating a table. Creates table, inserts rows,
// checks rows, and shuts down. Running it twice should result in error on table creation.
func TestFirstRun(t *testing.T) {

	// server is already initialized, so we can use it.
	// load database through runtime api.
	db_, err := runtime.GetDb("re_db")
	if err != nil {
		log.Fatal(err)
	}

	// create table.
	err = createTable(db_)
	if err != nil {
		// In this run table should not exist, so it should fail.
		log.Fatal(err)
	}

	// insert few rows.
	err = insertRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check inserted rows.
	err = checkRows(db_)
	if err != nil {
		log.Fatal(err)
	}
}

// TestSecondRun tests second run (after restart of server).
// Loads table, checks rows, and shuts down. Running it multiple times should not result in error.
func TestSecondRun(t *testing.T) {

	// server is already initialized, so we can use it.
	// load database through runtime api.
	db_, err := runtime.GetDb("re_db")
	if err != nil {
		log.Fatal(err)
	}

	// shoud get create error.
	err = createTable(db_)
	if err == nil {
		log.Fatal("table should exist")
	} else {
		log.Printf("OK. error message: %v.", err)
	}

	// check inserted rows.
	err = checkRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	//do insert rollback test
	err = insertRowsRollback(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check inserted rows.
	err = checkRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	//do delete rollback test
	err = deleteRowsRollback(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check inserted rows.
	err = checkRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	// do update rollback test
	err = updateRowsRollback(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check inserted rows.
	err = checkRows(db_)
	if err != nil {
		log.Fatal(err)
	}

}

func TestThirdRun(t *testing.T) {

	// server is already initialized, so we can use it.
	// load database through runtime api.
	db_, err := runtime.GetDb("re_db")
	if err != nil {
		log.Fatal(err)
	}

	// check inserted rows.
	err = checkRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	// delete rows.
	err = deleteRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check rows after delete.
	err = checkRowsAfterDelete(db_)
	if err != nil {
		log.Fatal(err)
	}

	// update rows.
	err = updateRows(db_)
	if err != nil {
		log.Fatal(err)
	}

	// check rows after update.
	err = checkRowsAfterUpdate(db_)
	if err != nil {
		log.Fatal(err)
	}

}

func TestTableCreateFails(t *testing.T) {

	// server is already initialized, so we can use it.
	// load database through runtime api.
	db_, err := runtime.GetDb("re_db")
	if err != nil {
		log.Fatal(err)
	}

	// check if table exists.
	_, err = db_.GetTable("TestFailingTable01")
	if err == nil {
		// delete table.
		db_.DeleteTable("TestFailingTable01")
	}

	createTableJSON := `{
	  "TableName": "TestFailingTable01",
	  "TableStorages": "StorerDiskV1",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable"
	  ]
	}`

	// create table.
	dberr := db_.CreateTableJSON(createTableJSON, nil)

	if dberr == nil {
		log.Fatal("create table should give error")
	} else {
		log.Printf("OK. error message: %v.", dberr)
	}

	createTableJSON = `{
	  "TableName": "TestFailingTable01",
	  "TableStorages": "StorerMemV1 StorerDiskV1 StorerMemV1",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable"
	  ]
	}`

	// create table.
	dberr = db_.CreateTableJSON(createTableJSON, nil)

	if dberr == nil {
		log.Fatal("create table should give error")
	} else {
		log.Printf("OK. error message: %v.", dberr)
	}
}

// ## Sub Level Test Functions

func checkRows(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	// check inserted rows.
	row, dberr := tbl.GetRowByPrimaryKey(1)
	if dberr != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser01" {
		return errors.New("invalid row data")
	}

	row, dberr = tbl.GetRowByPrimaryKey(2)
	if dberr != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser02" {
		return errors.New("invalid row data")
	}

	row, dberr = tbl.GetRowByPrimaryKey(3)
	if err != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser03" {
		return errors.New("invalid row data")
	}

	return nil
}

func checkRowsAfterDelete(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	// check rows.
	// check row 1.
	row, dberr := tbl.GetRowByPrimaryKey(1)
	if dberr != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser01" {
		return errors.New("invalid row data")
	}

	// check row 2.
	row, dberr = tbl.GetRowByPrimaryKey(2)
	if dberr == nil {
		return errors.New("row should not exist")
	}

	// check row 3.
	row, dberr = tbl.GetRowByPrimaryKey(3)
	if dberr != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser03" {
		return errors.New("invalid row data")
	}

	return nil
}

func checkRowsAfterUpdate(db_ *db.DB) error {

	// get table.
	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	//check all rows just to be sure update operation has not affected other rows.

	// check rows.
	// check row 1.
	row, dberr := tbl.GetRowByPrimaryKey(1)
	if dberr != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser01" {
		return errors.New("invalid row data")
	}

	// check row 2, deleted.
	row, dberr = tbl.GetRowByPrimaryKey(2)
	if dberr == nil {
		return errors.New("row should not exist")
	}

	// check row 3 after update.
	row, dberr = tbl.GetRowByPrimaryKey(3)
	if dberr != nil {
		return dberr.ToError()
	}
	if row["UserName"] != "TestUser03Updated" {
		return errors.New("invalid row data")
	}

	return nil
}

func createTable(db_ *db.DB) error {

	createTableJSON := `{
	  "TableName": "TestTable02",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable"
	  ]
	}`

	return db_.CreateTableJSON(createTableJSON, nil).ToError()
}

func deleteRows(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	// delete rows.
	dberr := tbl.DeleteRowByPrimaryKey(2)
	if dberr != nil {
		return dberr.ToError()
	}

	return nil
}

func deleteRowsRollback(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	// set rollback flag.
	dbtable_disk_store_v1.DeleteRollbackTest = true
	defer func() {
		dbtable_disk_store_v1.DeleteRollbackTest = false
	}() // reset rollback flag.

	// delete rows.
	dberr := tbl.DeleteRowByPrimaryKey(2)
	if dberr == nil {
		return errors.New("delete should fail")
	} else {
		log.Printf("OK. error message: %v.", dberr)
	}

	return nil
}

func deleteTable(db_ *db.DB) error {
	// check if table exists.
	_, err := db_.GetTable("TestTable02")
	if err != nil {
		return nil
	}

	// delete table.
	return db_.DeleteTable("TestTable02").ToError()
}

func insertRows(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	insertRowJSON := `{
		"Id": 1,
		"UserName": "TestUser01"
	}`

	dberr := tbl.InsertRowJSON(insertRowJSON)
	if dberr != nil {
		return dberr.ToError()
	}

	insertRowJSON = `{
		"Id": 2,
		"UserName": "TestUser02"
	}`

	dberr = tbl.InsertRowJSON(insertRowJSON)
	if dberr != nil {
		return dberr.ToError()
	}

	insertRowJSON = `{
		"Id": 3,
		"UserName": "TestUser03"
	}`

	dberr = tbl.InsertRowJSON(insertRowJSON)
	if dberr != nil {
		return dberr.ToError()
	}

	return nil
}

func insertRowsRollback(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	insertRowJSON := `{
		"Id": 4,
		"UserName": "TestUser04"
	}`

	dbtable_disk_store_v1.InsertRollbackTest = true
	defer func() {
		dbtable_disk_store_v1.InsertRollbackTest = false
	}() // reset

	dberr := tbl.InsertRowJSON(insertRowJSON)
	if dberr == nil {
		return errors.New("rollback test failed")
	} else {
		// TODO: this is correct or should return message?
		log.Printf("OK. error message: %v.", dberr)
	}

	return nil
}

//TODO: test replace update.

func updateRows(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	updateRowJSON := `{
		"UserName": "TestUser03Updated"
	}`

	dberr := tbl.UpdateRowJsonByPk(3, updateRowJSON, idbstorer.TableOperationPatchRow)
	if dberr != nil {
		return dberr.ToError()
	}

	return nil
}

func updateRowsRollback(db_ *db.DB) error {

	tbl, err := db_.GetTable("TestTable02")
	if err != nil {
		return err
	}

	updateRowJSON := `{
		"UserName": "TestUser03Updated"
	}`

	//TODO: test with pk value in json same and different.

	dbtable_disk_store_v1.UpdateRollbackTest = true
	defer func() {
		dbtable_disk_store_v1.UpdateRollbackTest = false
	}() // reset

	dberr := tbl.UpdateRowJsonByPk(1, updateRowJSON, idbstorer.TableOperationPatchRow)
	if dberr == nil {
		return errors.New("rollback test failed")
	} else {
		log.Printf("OK. error message: %v.", dberr)
	}

	return nil
}
