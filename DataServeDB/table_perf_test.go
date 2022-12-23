package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
)

/*

	Description: Tests performance of table operations.

	Notes:
		- Tests should only test performance and not the correctness of the code.
		- Use REST API as it will be mainly used in production.
		- Go has built in benchmarking api, use that.
		- Table create, change, and delete operations will not be used frequently in production, so they are not tested here.
		- Rows operations are tested here.
			- Use mix of insert, update, and delete operations.
			- Use mix of single and multiple rows operations.
			- Use 30% write operations and 70% read operations.
			- Use 50% write operations and 50% read operations.
			- Use 70% write operations and 30% read operations.
			- Above is to check performance of different write/read ratios. Usually, there will be more read operations than write operations.
			- Do tests with different number of rows in the table. 10, 100, 1,000, 100,000, and 1,000,000 rows.
			- For write operations also update indexed and non-indexed columns and check performance.

*/

//	type testCase struct {
//		method string
//		path   string
//		body   string
//		exp    error
//	}

type user struct {
	Id         int64
	UserName   string
	Occupation string
	Rank       int64
}

var randMIn, randMax int

func (us *user) toString() string {
	b, err := json.Marshal(us)
	if err != nil {
		// fmt.Println(err)
		return ""
	}
	return string(b)
}

var methods = []string{
	"GET",
	"POST",
	"PATCH",
	"PUT",
	"DELETE",
}
var Occupations = []string{
	"Security",
	"Waiter",
	"Electrician",
	"Doorman",
	"Repairman",
	"Mechanic",
}

func generateJSON(count int) (testCaseArray []testCase) {
	deleted := ""
	path := ""
	method := ""

	for i := 0; i < count; {
		method = getMethod(i)
		path = getPath(method, i)
		if method == "DELETE" && deleted == path {
			continue
		} else if method == "DELETE" {
			deleted = path
		}

		testCaseArray = append(testCaseArray,
			testCase{
				method,
				path,
				getBody(method, i),
				nil,
			},
		)
		if method == "POST" {
			i++
		}
	}
	return
}

func getBody(method string, i int) string {
	if method == "POST" {
		rUser := user{
			int64(i),
			"TestUser" + strconv.Itoa(i) + "InTestTable03",
			Occupations[rand.Intn(len(Occupations))],
			int64(i),
		}
		return rUser.toString()
	} else if method == "PATCH" || method == "PUT" {
		return `{ "UserName": "TestUser` + strconv.Itoa(randMax) + `InTestTable03Replaceded" }`
	}
	return ""

}
func getPath(method string, i int) string {
	if method == "GET" || method == "PATCH" || method == "PUT" {
		return "re_db/tables/TestTable03/" + strconv.Itoa(randMax)
	} else if method == "DELETE" {
		return "re_db/tables/TestTable03/" + strconv.Itoa((i/2)-1)
	}
	return "re_db/tables/TestTable03"
}
func getMethod(i int) string {
	if i > 10 {
		randMIn = rand.Intn(i / 2)
		randMax = rand.Intn(i/2) + (i / 2)
		return methods[rand.Intn(len(methods))]
	}
	return "POST"
}

func TestPer(t *testing.T) {
	arr := generateJSON(3000)
	TestDeleteTableRApi(t)
	TestCreateTableRApi(t)
	for _, testCase := range arr {
		t.Run("test "+testCase.method+" "+testCase.path, func(t *testing.T) {
			// testCaseEl := testCase
			// t.Parallel()
			fmt.Println(testCase.method, testCase.path, testCase.body)
			successResult, err := restApiCall(testCase.method, testCase.path, testCase.body)
			if err != testCase.exp {
				t.Errorf("%v\n", err)
			} else {
				log.Println(successResult)
			}
		})
	}
}
