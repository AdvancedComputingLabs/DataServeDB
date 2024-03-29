package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
)

/*
	Descripion:

	Plan:
		- create table 			[done]
		- insert row			[done]
		- get row				[done]
		- update row (partial)	[done]
		- get row after update
		- replace row (full)	[done]
		- check primary key is
			not changeable.		[implemented, but not tested]
		- check not nullable
			fields on replace.
		- check fields with
			default value on
			replace.
		- delete row
		- get row after delete
		- delete table

	Notes:
		- Don't need to do the same tests as table_external_storage_test.go,
			only the REST API CRUD operations and http status codes. Error codes needs testing here?

*/

const hostName = "localhost:8080"

func init() {
	// set testing flags
	//testing.T.Log(log.LstdFlags, log.Lshortfile) //coud not find proper way to set formatting for testing log output
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// ## First Level Test Functions

func TestCreateTableRApi(t *testing.T) {

	createTableJSON := `{
	  "TableName": "TestTable03",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Occupation string Length:0..50 Nullable",
		"Rank int32 default: 100 Nullable"
	  ]
	}`

	successResult, err := restApiCall("POST", "re_db/tables", createTableJSON)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func TestDeleteTableRApi(t *testing.T) {

	successResult, err := restApiCall("DELETE", "re_db/tables/TestTable03)", "")
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	} else {
		log.Println(successResult)
	}
}

func TestInsertRecordsRApi(t *testing.T) {

	insertRowJSON := `{
		"Id": 1,
		"UserName": "TestUser01InTestTable03",
		"Occupation": "Security",
		"Rank": 101
	}`

	// NOTE: added rank to the insert json, because it was created as nullable.

	successResult, err := restApiCall("POST", "re_db/tables/TestTable03", insertRowJSON)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func TestGetRecordsRApi(t *testing.T) {

	// get tables
	successResult, err := restApiCall("GET", "re_db/tables", "")
	if err != nil {
		// not implemented yet
		//log.Fatal(err)
		log.Println(err)
	} else {
		log.Println(successResult)
	}

	// get row
	successResult, err = restApiCall("GET", "re_db/tables/TestTable03/1", "")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}

	// get table function result
	successResult, err = restApiCall("GET", "re_db/tables/TestTable03/$HelloFrom()", "")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func TestUpdateRecordsRApi(t *testing.T) {

	updateRowJSON := `{
		"UserName": "TestUser01InTestTable03Updated"
	}`

	successResult, err := restApiCall("PATCH", "re_db/tables/TestTable03/1", updateRowJSON)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func TestReplaceRecordsRApi(t *testing.T) {

	// NOTE: non-nullable field tested but need to have its own test.
	// TODO: add test for non-nullable replace test (should fail; not include non-nullable field in replace json)

	// IMPORTANT-NOTE: Rank was not removed from the record, because it was created as not nullable.
	// 		After deleting the table and recreating it with Rank as nullable, Rank was removed from the record.
	//		-- HY @ 21-Nov-2022

	updateRowJSON := `{
		"UserName": "TestUser01InTestTable03Replaceded"
	}`

	successResult, err := restApiCall("PUT", "re_db/tables/TestTable03/1", updateRowJSON)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}

}

func TestDeleteRecordsRApi(t *testing.T) {

	successResult, err := restApiCall("DELETE", "re_db/tables/TestTable03/1", "")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}

}

// ## Sub Level Test Functions

// ## Helper Functions

func newHttpReqNResp(method, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, fmt.Sprintf("http://%s/%s", hostName, path), body)
	w := httptest.NewRecorder()
	return req, w
}

func restApiCall(method, path, bodyJson string) (string, error) {

	req, w := newHttpReqNResp(method, path, io.NopCloser(bytes.NewReader([]byte(bodyJson))))

	reqPath := rest.HttpRestPathParser(req.URL.String())

	dbrouter.MatchPathAndCallHandler(w, req, reqPath, req.Method)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return "", fmt.Errorf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, string(body))
	} else {
		return fmt.Sprintf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, string(body)), nil
	}
}
