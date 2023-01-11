package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
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

type user struct {
	Id         int64
	UserName   string
	Occupation string
	Rank       int64
}
type testCase struct {
	Method string
	Path   string
	Body   string
	Exp    error
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
func getJson() []testCase {
	jsonFile, err := os.Open("test_data.json")
	if err != nil {
		log.Println(err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []testCase
	json.Unmarshal([]byte(byteValue), &result)

	return result

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
	} else if method == "PATCH" {
		return `{ "UserName": "TestUser` + strconv.Itoa(randMax) + `InTestTable03Replaceded" }`
	}
	return ""

}
func getPath(method string, i int) string {
	if method == "GET" || method == "PATCH" {
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
	//arr := generateJSON(1000)
	// for _, v := range arr {
	// 	b, err := json.Marshal(v)
	// 	if err != nil {
	// 		fmt.Printf("Error: %s", err)
	// 		return
	// 	}
	// 	fmt.Println(string(b))
	// }`
	// db, err := os.OpenFile("test_data.json", os.O_EXCL|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	if !os.IsExist(err) {
	// 		return
	// 	}
	// }
	// defer db.Close()
	// enc := json.NewEncoder(db)
	// err = enc.Encode(&arr)
	// if err != nil {
	// 	return
	// }
	// return
	arr := getJson()
	// }
	TestDeleteTableRApi(t)
	TestCreateTableRApi(t)
	for _, testCase := range arr {
		t.Run("test "+testCase.Method+" "+testCase.Path, func(t *testing.T) {
			fmt.Println(testCase.Method, testCase.Path, testCase.Body)
			successResult, err := restApiCall(testCase.Method, testCase.Path, testCase.Body)
			if err != testCase.Exp {
				t.Errorf("%v\n", err)
			} else {
				log.Println(successResult)
			}
		})
	}
}

type calculation struct {
	index int
	mutex sync.Mutex
}

func BenchmarkPerf(b *testing.B) {
	array := getJson()
	c := calculation{}

	b.Run("test ", func(b *testing.B) {
		fmt.Println("benchmarkN --> ", b.N)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.mutex.Lock()
				fmt.Println(array[c.index].Method, array[c.index].Path, array[c.index].Body)
				successResult, err := restApiCall(array[c.index].Method, array[c.index].Path, array[c.index].Body)
				if err != array[c.index].Exp {
					b.Errorf("%v\n", err)
				} else {
					log.Println(successResult)
				}
				c.index++
				c.mutex.Unlock()
			}
		})
	})
}
