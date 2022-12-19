package dbrouter

import (
	"fmt"
	"net/http"
	"testing"
)

func getLastPathLevel(pathLevels []PathLevel) (lastPathLevel PathLevel) {
	if pathLevels != nil && len(pathLevels) > 0 {
		lastPathLevel = pathLevels[len(pathLevels)-1]
	}
	return
}

func TableHandlerTester(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName string, pathLevels []PathLevel) {

	lastPathLevel := getLastPathLevel(pathLevels)

	fmt.Println("Table handler!")
	fmt.Println("resPath:", resPath)
	fmt.Println("matchedPath:", matchedPath)
	fmt.Println("dbName:", dbName)
	fmt.Println("lastPathLevel:", lastPathLevel.PathItem)
	fmt.Println("lastPathLevel_DbResTypeId:", lastPathLevel.PathItemTypeId.String())
}

func TablesParentHandlerTester(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName string, pathLevels []PathLevel) {

	lastPathLevel := getLastPathLevel(pathLevels)

	fmt.Println("Tables namespace handler!")
	fmt.Println("resPath:", resPath)
	fmt.Println("matchedPath:", matchedPath)
	fmt.Println("dbName:", dbName)
	fmt.Println("lastPathLevel:", lastPathLevel.PathItem)
	fmt.Println("lastPathLevel_DbResTypeId:", lastPathLevel.PathItemTypeId.String())
}

func TestRegister(t *testing.T) {
	// runtime.CreateDBmeta()
	//runtime.InitMapOfDB()
	Register("{DB_NAME}/tables/{TBL_NAME}/{1}.*", TableHandlerTester)
	Register("{DB_NAME}/tables/{TBL_NAME}", TableHandlerTester)
	Register("{DB_NAME}/tables", TablesParentHandlerTester)
	testMatchPathAndCallHandlerForTableLevelOperations(t)
}

func testMatchPathAndCallHandlerForTableLevelOperations(t *testing.T) {
	//NOTE: both path formats works:
	// 1) re_db/tables/users
	// 2) /re_db/tables/users
	//TODO: add test cases to detect if this behavior breaks in future update.
	//TODO: check if odata object name is case sensitive?

	println("creating table...")
	MatchPathAndCallHandler(nil, nil, "/re_db/tables.$create()", "POST")
	println("")
	println("deleting table...")
	// /re_db/tables.$delete('TableName')
	MatchPathAndCallHandler(nil, nil, "/re_db/tables/TableName", "POST")

	println("")
	println("getting row...")
	MatchPathAndCallHandler(nil, nil, "/re_db/tables/uSers/Id:1", "GET")
}
