package main

import (
	"log"
	"testing"
)

func FuzzFile(f *testing.F) {

	for _, seed := range testCasefileArray {
		f.Add(seed.method, seed.path, seed.filename)
	}
	f.Fuzz(func(t *testing.T, method, path, body string) {
		if method == "POST" {
			successResult, err := restApiCallMu(method, path, body)
			if err != nil {
				t.Fatalf("%v\n", err)
			} else {
				log.Println(successResult)
			}
		} else {
			successResult, _, err := restApiCall(method, path, body)
			if err != nil {
				t.Fatalf("%v\n", err)
			} else {
				log.Println(successResult)
			}
		}
	})
}
