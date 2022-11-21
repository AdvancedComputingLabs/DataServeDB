package rest

import (
	"fmt"
	"testing"
)

func TestParseKeyValue(t *testing.T) {

	// Test Cases
	// 1. /re_db/tables/Tbl01/Id:
	// 2. /re_db/tables/Tbl01/Id:1
	// 3. /re_db/tables/Tbl01/
	// 4. /re_db/tables/Tbl01

	key, value, err := ParseKeyValue("/re_db/tables/Tbl01/")
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("key:", key, "; value:", value)
	}
}
