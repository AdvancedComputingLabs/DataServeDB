package main

import (
	"DataServeDB/commtypes"
	"DataServeDB/unstable_api/runtime"
	"fmt"
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
	set := []struct {
		query  string
		result string
	}{
		{`{"Users": {
		"Id": {},
		"UserName": {},
		"Properties": [
		   {
			 "$JOIN": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum",
			 "$WHERE": "Properties.SlNum >= 2"
		   }
		 ]
	   	}}`,
			`{"Users":[{"Id":0,"Properties":[{"Counter":1,"DateAdded":"2021-02-21T10:38:40.0002297Z","GlobalId":"735aaaba-e953-4fa5-8fa6-701566d6ed0c","Name":"JLT1","RoleIDs":"45ae2266-f720-4d4e-9fe4-0900478002ee","SlNum":"acl000"}],"UserName":"captain america"},{"Id":1,"Properties":[{"Counter":2,"DateAdded":"2021-02-21T10:38:40.2654551Z","GlobalId":"0d032a25-3fa5-43c0-ae52-c36be28e60db","Name":"JLTX2","RoleIDs":"ec021e45-739f-45e6-9f51-21bc7426abd0","SlNum":"acl001"}],"UserName":"IRON MAN"},{"Id":2,"Properties":[{"Counter":3,"DateAdded":"2021-02-21T10:38:40.4428508Z","GlobalId":"a8b9a4cd-4121-4a3a-83b1-ff36fce67283","Name":"JLTX3","RoleIDs":"c7fb46ed-3711-42ff-91e1-35249528c04b","SlNum":"acl002"}],"UserName":"professor HULK"},{"Id":3,"Properties":[{"Counter":4,"DateAdded":"2021-02-21T10:38:40.4882851Z","GlobalId":"e5c47104-17f9-4775-9cf1-b5bd6665c43b","Name":"JLTX4","RoleIDs":"85c8efcf-1d4c-4f16-9a7f-7179931fd1e6","SlNum":"acl003"}],"UserName":"peter Parker"}]}`},
		{`{"Users": {}}`,
			`{"Users":[{"Counter":1,"DateAdded":"2021-02-18T08:09:15.8003188Z","GlobalId":"538a7229-b76c-48b9-8ede-667af7102cbd","Id":0,"UserName":"captain america"},{"Counter":2,"DateAdded":"2021-02-18T08:09:15.8452329Z","GlobalId":"221805eb-9482-4c80-b889-7c646cc819eb","Id":1,"UserName":"IRON MAN"},{"Counter":3,"DateAdded":"2021-02-18T08:09:15.8788338Z","GlobalId":"39356318-b3ec-4a10-a6b3-d1ef69ab5216","Id":2,"UserName":"professor HULK"},{"Counter":4,"DateAdded":"2021-02-18T08:09:15.9133239Z","GlobalId":"60cd35cc-4007-416b-be4c-a12e29ebfad9","Id":3,"UserName":"peter Parker"}]}`},
	}
	db, e := runtime.GetDb("re_db")
	if e != nil {
		t.Errorf("%v\n", e)
		return
	}
	// var dbReqCtx *commtypes.DbReqContext
	dbReqCtx := commtypes.NewDbReqContext(
		"", "", "",
		"re_db", db, "", 1)
	for _, testSet := range set {
		_, qry, err := runtime.DecodeJSON([]byte(testSet.query))
		if err != nil {
			t.Errorf("%v\n", err)
			return
		}
		stat, res, err := db.TablesQueryGet(dbReqCtx, qry)
		if err != nil {
			t.Errorf("%v\n", e)
			return
		}
		_ = stat
		if testSet.result != string(res) {
			t.Errorf("error, test Failed")
			fmt.Println(" %V", string(res))
		} else {
			fmt.Println("OK! test passed \n", string(res))
		}
	}

}
