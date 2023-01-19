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
		f.Add(seed.Method, seed.Path, seed.Body)
	}
	f.Fuzz(func(t *testing.T, method, path, body string) {
		successResult, _, err := restApiCall(method, path, body)
		if err != nil {
			t.Fatalf("%v\n", err)
		} else {
			log.Println(successResult)
		}
	})
}
