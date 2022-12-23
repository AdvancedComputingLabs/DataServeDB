package main

import (
	"fmt"
	"log"
	"testing"
)

/*
Description: Tests concurrency related bugs and issues for table operations. Mainly race conditions and deadlocks.

Notes:
  - Use REST API calls, if they are causing problems in deadlock detection then direct table calls.
  - Use mix of read and write operations.
  - Use mix of single and multiple rows operations.
  - Use go's own concurrency testing. See: https://forum.golangbridge.org/t/how-do-you-unit-test-a-concurrent-data-structure/26912
  - If there is better 3rd party library for concurrency testing, then get it approved.
*/
type testCase struct {
	method string
	path   string
	body   string
	exp    error
}

var testCaseArray = []testCase{
	{"GET", "re_db/tables/TestTable03/4", "", nil},
	{"GET", "re_db/tables/TestTable03/1", "", nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 11, "UserName": "TestUser11InTestTable03", "Occupation": "Security", "Rank": 111}`, nil},
	{"DELETE", "re_db/tables/TestTable03/6", "", nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 12, "UserName": "TestUser12InTestTable03", "Occupation": "Security", "Rank": 112}`, nil},
	{"GET", "re_db/tables/TestTable03/2", "", nil},
	{"GET", "re_db/tables/TestTable03/3", "", nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 10, "UserName": "TestUser10InTestTable03", "Occupation": "Security", "Rank": 110}`, nil},
	{"DELETE", "re_db/tables/TestTable03/5", "", nil},
	{"GET", "re_db/tables/TestTable03/4", "", nil},
	{"DELETE", "re_db/tables/TestTable03/7", "", nil},
	{"DELETE", "re_db/tables/TestTable03/8", "", nil},
	{"GET", "re_db/tables/TestTable03/3", "", nil},
	{"GET", "re_db/tables/TestTable03/4", "", nil},
	{"GET", "re_db/tables/TestTable03/3", "", nil},
	{"GET", "re_db/tables/TestTable03/2", "", nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 9, "UserName": "TestUser09InTestTable03", "Occupation": "Security", "Rank": 109}`, nil},
	{"GET", "re_db/tables/TestTable03/2", "", nil},
}

func TestCuncurrency(t *testing.T) {

	for i := 0; i < len(testCaseArray); i++ {
		t.Run("test "+testCaseArray[i].method+" "+testCaseArray[i].path, func(t *testing.T) {
			testCaseEl := testCaseArray[i]
			t.Parallel()
			fmt.Println(testCaseEl.method, testCaseEl.path, testCaseEl.body)
			successResult, err := restApiCall(testCaseEl.method, testCaseEl.path, testCaseEl.body)
			if err != testCaseEl.exp {
				t.Errorf("%v\n", err)
			} else {
				log.Println(successResult)
			}
		})
	}
}

// func concurrentTesting(t *testing.T, testCaseElem testCase) {

// 	fmt.Println(testCaseElem.method, testCaseElem.path, testCaseElem.body)
// 	successResult, err := restApiCall(testCaseElem.method, testCaseElem.path, testCaseElem.body)
// 	if err != testCaseElem.exp {
// 		t.Errorf("%v\n", err)
// 	} else {
// 		log.Println(successResult)
// 	}

// }

// func TestInsertRecords(t *testing.T) {

// 	for _, rowJSON := range insertRowJSON {
// 		successResult, err := restApiCall("POST", "re_db/tables/TestTable03", rowJSON)
// 		if err != nil {
// 			log.Fatal(err)
// 		} else {
// 			log.Println(successResult)
// 		}
// 	}
// }

// func TestGetRecords(t *testing.T) {

// }
// func getRecords() {
// 	successResult, err := restApiCall("GET", "re_db/tables/TestTable03/4", "")
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		log.Println(successResult)
// 	}
// }

// func TestReplaceRecords(t *testing.T) {

// 	// NOTE: non-nullable field tested but need to have its own test.
// 	// TODO: add test for non-nullable replace test (should fail; not include non-nullable field in replace json)

// 	// IMPORTANT-NOTE: Rank was not removed from the record, because it was created as not nullable.
// 	// 		After deleting the table and recreating it with Rank as nullable, Rank was removed from the record.
// 	//		-- HY @ 21-Nov-2022

// 	updateRowJSON := []string{
// 		`{ "UserName": "TestUser01InTestTable03Replaceded" }`,
// 		`{ "UserName": "TestUser02InTestTable03Replaceded" }`,
// 		`{ "UserName": "TestUser03InTestTable03Replaceded" }`,
// 		`{ "UserName": "TestUser04InTestTable03Replaceded" }`,
// 	}

// 	for i, rowJSON := range updateRowJSON {
// 		successResult, err := restApiCall("PUT", "re_db/tables/TestTable03/"+strconv.Itoa(i+1), rowJSON)
// 		if err != nil {
// 			log.Fatal(err)
// 		} else {
// 			log.Println(successResult)
// 		}
// 	}

// }

// func TestDeleteRecords(t *testing.T) {

// 	for i := 5; i <= 8; i++ {
// 		successResult, err := restApiCall("DELETE", "re_db/tables/TestTable03/"+strconv.Itoa(i), "")
// 		if err != nil {
// 			t.Error(err)
// 		} else {
// 			log.Println(successResult)
// 		}
// 	}

// }
