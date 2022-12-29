package main

import (
	"log"
	"testing"
)

/*
Description: Fuzz tests for table operations.

TODO: Devise useful fuzz test cases.
*/

func FuzzHex(f *testing.F) {
	for i, seed := range testCaseArray {
		f.Add(seed.method, seed.path, seed.body)
		if i == 1000 {
			break
		}
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
