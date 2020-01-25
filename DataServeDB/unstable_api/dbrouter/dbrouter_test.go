package dbrouter

import (
	"fmt"
	"net/http"
	"testing"
)

func TableHandlerTester(w http.ResponseWriter, r *http.Request, httpMethod string, dbName string, resPath string) (resultHttpStatus int, resultContent []byte, resultErr error) {
	fmt.Println("Tables handler!")
	return 0, nil, nil
}

func TestRegister(t *testing.T) {
	Register("{db_name}/tables/{tbl_name}", TableHandlerTester)
	testMatchPathAndCallHandler(t)
}

func testMatchPathAndCallHandler(t *testing.T) {
	MatchPathAndCallHandler(nil, nil, "re_db/tables/users", "GET")
}
