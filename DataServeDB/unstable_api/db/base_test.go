package db

import "testing"

func TestExtractFullCallStr(t *testing.T) {

	s := extractFullCallStrTablesParent("re_db/tables.$create()", "tables")
	if s != "tables.$create()" {
		t.Fatal("extractFullCallStrTablesParent failed!")
	}

	s = extractFullCallStrTablesParent("re_db/tables.$delete('TableName')", "tables")
	if s != "tables.$delete('TableName')" {
		t.Fatal("extractFullCallStrTablesParent failed!")
	}

	s = extractFullCallStrTablesParent("re_db/tables/", "tables")
	if s != "tables" {
		t.Fatal("extractFullCallStrTablesParent failed!")
	}

	s = extractFullCallStrTablesParent("re_db/tables", "tables")
	if s != "tables" {
		t.Fatal("extractFullCallStrTablesParent failed!")
	}
}
