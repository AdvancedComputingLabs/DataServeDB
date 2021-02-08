package dbrouter

import (
	"fmt"
	"net/http"
	"testing"

	"DataServeDB/dbsystem/constants"
)

func TableHandlerTester(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName, targetName string, targetDbResTypeId constants.DbResTypes) (resultHttpStatus int, resultContent []byte, resultErr error) {
	fmt.Println("Tables handler!")
	fmt.Println("resPath:", resPath)
	fmt.Println("matchedPath:", matchedPath)
	fmt.Println("dbName:", dbName)
	fmt.Println("targetName:", targetName)
	fmt.Println("targetDbResTypeId:", targetDbResTypeId)
	fmt.Println("")
	return 0, nil, nil
}

func QueryHandlerTester(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName, targetName string, targetDbResTypeId constants.DbResTypes) (resultHttpStatus int, resultContent []byte, resultErr error) {
	fmt.Println("Query handler!")
	fmt.Println("resPath:", resPath) //NOTE: resPath: /re_db/query/* needs handling.
	fmt.Println("matchedPath:", matchedPath)
	fmt.Println("dbName:", dbName)
	fmt.Println("targetName:", targetName)
	fmt.Println("targetDbResTypeId:", targetDbResTypeId)
	fmt.Println("")
	return 0, nil, nil
}

func TestRegister(t *testing.T) {
	// runtime.CreateDBmeta()
	//runtime.InitMapOfDB()
	Register("{db_name}/tables/{tbl_name}", TableHandlerTester)
	Register("{db_name}/query", QueryHandlerTester)
	testMatchPathAndCallHandler(t)
}

func testMatchPathAndCallHandler(t *testing.T) {
	//NOTE: both path formats works:
	// 1) re_db/tables/users
	// 2) /re_db/tables/users
	//TODO: add test cases to detect if this behavior breaks in future update.

	MatchPathAndCallHandler(nil, nil, "/re_db/tables/users/Id:1", "GET")
	MatchPathAndCallHandler(nil, nil, "/re_db/query", "POST")
}
