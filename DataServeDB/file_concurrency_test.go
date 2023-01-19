package main

import (
	"log"
	"testing"
)

type testCaseFile struct {
	method   string
	path     string
	filename string
	exp      error
}

var testCasefileArray = []testCaseFile{
	{"GET", "re_db/files/level1", "", nil},
	{"GET", "re_db/files/level1/level2", "", nil},
	{"POST", "re_db/files/level1", "PropertiesData.ods", nil},
	{"DELETE", "re_db/files/level1/ACLPropertyAppDummyData.xlsx", "", nil},
	{"POST", "re_db/files/level1/level2", "myJsonFile0.json", nil},
	{"GET", "re_db/files/level1", "", nil},
	{"GET", "re_db/files/level1/level2", "", nil},
	{"POST", "re_db/files/level1/level2", "storer_design_possible_implementation.pdf", nil},
	{"DELETE", "re_db/files/lvel1/list.txt", "", nil},
	{"GET", "re_db/files/level1/README.md", "", nil},
	{"DELETE", "re_db/files/level1/ipconfig.txt", "", nil},
	{"DELETE", "re_db/files/level1/Home_props.PNG", "", nil},
	{"GET", "re_db/files/level1/netstat.txt", "", nil},
	{"GET", "re_db/files/level1/netstat.txt", "", nil},
	{"GET", "re_db/files/level1/netstat.txt", "", nil},
	{"GET", "re_db/files/level1/netstat.txt", "", nil},
	{"GET", "re_db/files/level1/netstat.txt", "", nil},
	{"GET", "re_db/files/level1/netstat.txt", "", nil},
	{"GET", "re_db/files/level1/ping.txt", "", nil},
	{"GET", "re_db/files/level1/token.json", "", nil},
	{"GET", "re_db/files/level1/trace.txt", "", nil},
	{"POST", "re_db/files/level1", "task.txt", nil},
	{"GET", "re_db/files//level1/token.json", "", nil},
}

func TestFileCuncurrency(t *testing.T) {

	for _, testCase := range testCasefileArray {
		t.Run("test "+testCase.method+" "+testCase.path, func(t *testing.T) {
			testCaseEl := testCase
			t.Parallel()
			if testCaseEl.method == "POST" {
				successResult, err := restApiCallMu(testCaseEl.method, testCaseEl.path, testCaseEl.filename)
				if err != nil {
					t.Fatal(err)
				} else {
					log.Println(successResult)
				}
			} else {
				successResult, _, err := restApiCall(testCaseEl.method, testCaseEl.path, testCaseEl.filename)
				if err != testCaseEl.exp {
					log.Println(err)
				} else {
					log.Println(successResult)
				}
			}
		})
	}

}
