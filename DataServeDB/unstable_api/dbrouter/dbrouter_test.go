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
	return 0, nil, nil
}

func TestRegister(t *testing.T) {
	// runtime.CreateDBmeta()
	//runtime.InitMapOfDB()
	Register("{db_name}/tables/{tbl_name}", TableHandlerTester)
	testMatchPathAndCallHandler(t)
}

func testMatchPathAndCallHandler(t *testing.T) {
	//NOTE: both path formats works:
	// 1) re_db/tables/users
	// 2) /re_db/tables/users
	//TODO: add test cases to detect if this behavior breaks in future update.

	MatchPathAndCallHandler(nil, nil, "/re_db/tables/users/Id:1", "GET")
}
