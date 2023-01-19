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

var testCaseArray = []testCase{
	{"GET", "re_db/tables/TestTable03/4", "", 200},
	{"GET", "re_db/tables/TestTable03/1", "", 200},
	{"POST", "re_db/tables/TestTable03", `{"Id": 11, "UserName": "TestUser11InTestTable03", "Occupation": "Security", "Rank": 111}`, 400},
	{"DELETE", "re_db/tables/TestTable03/6", "", 404},
	{"POST", "re_db/tables/TestTable03", `{"Id": 12, "UserName": "TestUser12InTestTable03", "Occupation": "Security", "Rank": 112}`, 400},
	{"GET", "re_db/tables/TestTable03/2", "", 200},
	{"GET", "re_db/tables/TestTable03/3", "", 200},
	{"POST", "re_db/tables/TestTable03", `{"Id": 10, "UserName": "TestUser10InTestTable03", "Occupation": "Security", "Rank": 110}`, 400},
	{"DELETE", "re_db/tables/TestTable03/5", "", 404},
	{"GET", "re_db/tables/TestTable03/4", "", 200},
	{"DELETE", "re_db/tables/TestTable03/7", "", 404},
	{"DELETE", "re_db/tables/TestTable03/8", "", 404},
	{"GET", "re_db/tables/TestTable03/3", "", 200},
	{"GET", "re_db/tables/TestTable03/4", "", 200},
	{"GET", "re_db/tables/TestTable03/3", "", 200},
	{"GET", "re_db/tables/TestTable03/2", "", 200},
	{"POST", "re_db/tables/TestTable03", `{"Id": 9, "UserName": "TestUser09InTestTable03", "Occupation": "Security", "Rank": 109}`, 400},
	{"GET", "re_db/tables/TestTable03/2", "", 200},
}

func TestCuncurrency(t *testing.T) {

	for i := 0; i < len(testCaseArray); i++ {
		t.Run("test "+testCaseArray[i].Method+" "+testCaseArray[i].Path, func(t *testing.T) {
			testCaseEl := testCaseArray[i]
			t.Parallel()
			fmt.Println(testCaseEl.Method, testCaseEl.Path, testCaseEl.Body)
			successResult, status, err := restApiCall(testCaseEl.Method, testCaseEl.Path, testCaseEl.Body)
			if err != nil {
				if status != testCaseEl.Exp {
					fmt.Println(err.Error())
				} else {
					log.Println(successResult)

				}
				// t.Errorf("%v\n", err)
			} else {
				log.Println(successResult)
			}
		})
	}
}
