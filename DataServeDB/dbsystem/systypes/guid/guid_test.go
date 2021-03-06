package guid

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseString(t *testing.T) {
	type DtTester struct {
		Title      string
		GlobalId Guid
	}

	g, e := ParseString("9f821cfd-9215-4566-84e4-c6f67ee25914")
	tester := DtTester{Title:"GuidTest", GlobalId: *g}
	if e != nil {
		t.Error(e)
	}

	{ // json

		testerJson, e := json.Marshal(tester)
		if e != nil {
			t.Error(e)
		} else {
			fmt.Println(string(testerJson))
		}

		var UnmarshalledTester DtTester
		e = json.Unmarshal(testerJson, &UnmarshalledTester)
		if e != nil {
			t.Error(e)
		} else {
			fmt.Println(UnmarshalledTester)
		}
	}

	{ // gob

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		e := enc.Encode(tester)
		if e != nil {
			t.Error(e)
		} else {
			fmt.Println(buf.String())
		}

		var testerDecoded DtTester
		bufDecode := bytes.NewReader(buf.Bytes())
		dec := gob.NewDecoder(bufDecode)
		e = dec.Decode(&testerDecoded)
		if e != nil {
			t.Error(e)
		} else {
			fmt.Println(testerDecoded)
		}
	}
}
