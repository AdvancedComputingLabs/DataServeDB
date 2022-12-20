package main

import (
	"fmt"
	"log"
	"testing"
	"time"
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
	{"POST", "re_db/tables/TestTable03", `{"Id": 9, "UserName": "TestUser09InTestTable03", "Occupation": "Security", "Rank": 109}`, nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 10, "UserName": "TestUser10InTestTable03", "Occupation": "Security", "Rank": 110}`, nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 11, "UserName": "TestUser11InTestTable03", "Occupation": "Security", "Rank": 111}`, nil},
	{"POST", "re_db/tables/TestTable03", `{"Id": 12, "UserName": "TestUser12InTestTable03", "Occupation": "Security", "Rank": 112}`, nil},
	{"GET", "re_db/tables/TestTable03/1", "", nil},
	{"GET", "re_db/tables/TestTable03/2", "", nil},
	{"GET", "re_db/tables/TestTable03/3", "", nil},
	{"GET", "re_db/tables/TestTable03/4", "", nil},
	{"DELETE", "re_db/tables/TestTable03/5", "", nil},
	{"DELETE", "re_db/tables/TestTable03/6", "", nil},
	{"DELETE", "re_db/tables/TestTable03/7", "", nil},
	{"DELETE", "re_db/tables/TestTable03/8", "", nil},
}

func TestCuncurrency(t *testing.T, f *testing.F) {

	for _, testCaseElem := range testCaseArray {
		go concurrentTesting(t, testCaseElem)
	}

	time.Sleep(time.Second)
}
func concurrentTesting(t *testing.T, testCaseElem testCase) {

	fmt.Println(testCaseElem.method, testCaseElem.path, testCaseElem.body)
	successResult, err := restApiCall(testCaseElem.method, testCaseElem.path, testCaseElem.body)
	if err != testCaseElem.exp {
		t.Errorf("%v\n", err)
	} else {
		log.Println(successResult)
	}

}
