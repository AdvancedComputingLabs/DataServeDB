package runtime

import "DataServeDB/unstable_api/db"

// MapOfdb is exported
var MapOfdb map[string]db.DB

// func initMapOfDB(){
// 	MapOfdb["re_db"] = db.DB{
// 		dbName "re_db",
// 		dbInternalID 0,
// 		MapOfTables 
// 	}
// }