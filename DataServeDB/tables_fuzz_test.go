package main

import (
	"log"
	"testing"
)

/*
Description: Fuzz tests for table operations.

TODO: Devise useful fuzz test cases.
*/

func FuzzTable(f *testing.F) {
	for _, seed := range testCaseArray {
		f.Add(seed.method, seed.path, seed.body)
	}
	f.Fuzz(func(t *testing.T, method, path, body string) {
		successResult, err := restApiCall(method, path, body)
		if err != nil {
			t.Fatalf("%v\n", err)
		} else {
			log.Println(successResult)
		}
	})
}
