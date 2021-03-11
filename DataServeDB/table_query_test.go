package main

import (
	"DataServeDB/commtypes"
	"DataServeDB/unstable_api/runtime"
	"testing"
)

// var Use= {
// 	Id: {},
// 	"UserName": {},
// 	"Properties": [
// 	   {
// 		 "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
// 	   }
// 	 ]
// 	   }}
func TestQuary(t *testing.T) {
	query := `{"Users": {
		"Id": 1,
		"UserName": {},
		"Properties": [
		   {
			 "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
		   }
		 ]
	   	}}`
	//var dst interface{}

	// dec := json.NewDecoder(strings.NewReader(query))
	// for {
	// 	t, err := dec.Token()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%T: %v", t, t)
	// 	if dec.More() {
	// 		fmt.Printf(" (more)")
	// 	}
	// 	fmt.Printf("\n")
	// }

	_, qry, err := runtime.DecodeJSON([]byte(query))
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	db, e := runtime.GetDb("re_db")
	if e != nil {
		t.Errorf("%v\n", e)
		return
	}
	var dbReqCtx *commtypes.DbReqContext
	dbReqCtx = commtypes.NewDbReqContext(
		"", "", "",
		"re_db", db, "", 1)
	// b, e := json.Marshal(query)
	// if e != nil {
	// 	t.Errorf("%v\n", e)
	// }
	// println(string(b))
	stat, res, err := db.TablesQueryGet(dbReqCtx, qry)
	if err != nil {
		t.Errorf("%v\n", e)
		return
	}
	_ = stat
	println(string(res))

}
