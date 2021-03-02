package main

import (
	"DataServeDB/commtypes"
	"DataServeDB/unstable_api/runtime"
	"encoding/json"
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
		"Id": {},
		"UserName": {},
		"Properties": [
		   {
			 "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
		   }
		 ]
	   	}}`
	var dst interface{}
	json.Unmarshal([]byte(query), &dst)

	qry, err := runtime.DecodeJSON(dst)
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
	println("qry-> ", qry.ItemLabel)
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
