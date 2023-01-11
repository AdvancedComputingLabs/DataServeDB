package main

import (
	"testing"
)

func BenchmarkFile(b *testing.B) {
	b.Run("test ", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					// The loop body is executed b.N times total across all goroutines.
					_, err := restApiCall("GET", "re_db/files/level1/level2/storer_design_possible_implementation.pdf", "")
					if err != nil {
						b.Errorf("%v\n", err)
					} else {
						//log.Println(successResult)
					}
				}
			})
		}
	})
}
