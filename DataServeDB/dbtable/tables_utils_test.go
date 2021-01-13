package dbtable

import (
	"fmt"
	"testing"
)

func TestParseKeyValue(t *testing.T) {
	key, value, err := parseKeyValue("/re_db/tables/Tbl01/Id:")
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("key:", key, "; value:", value)
	}
}


