package runtime

import (
	"DataServeDB/dbtable"
	"DataServeDB/unstable_api/db"
	"errors"
	"fmt"
)

// MapOfdb is exported
var MapOfdb = make(map[string]db.DB)

func GetDB(dbName string) (db.DB, error) {
	if data, ok := MapOfdb[dbName]; ok {
		return data, nil
	}
	return db.DB{}, errors.New("db Name not found")
}
func LoadDB() error {
	var mapoftable = make(db.MapOfTables)
	createTableJSON := `{
	  "TableName": "users",
	  "TableFields": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable",
		"DateAdded datetime default:Now() !Nullable",
		"GlobalId guid default:NewGuid() !Nullable"
	  ]
	}`

	if tbl01, err := dbtable.CreateTableJSON(createTableJSON); err == nil {
		if jsonStr, err := dbtable.GetSaveLoadStructure(tbl01); err == nil {
			fmt.Println(jsonStr)
			if tbl, err := dbtable.LoadTableFromDB(jsonStr); err == nil {
				fmt.Printf("table loaded :- %v\n", tbl)
				mapoftable["users"] = *tbl
				MapOfdb["re_db"] = db.DB{
					DbName:       "re_db",
					DbInternalID: 0,
					MapOfTables:  mapoftable,
				}
			} else {
				println("err", err.Error())
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
	// fmt.Printf("%v\n", MapOfdb)
	// gob.Register(dtIso8601Utc.Iso8601Utc{})
	// gob.Register(guid.Guid{})
	// gob.Register()
	// var network bytes.Buffer        // Stand-in for a network connection
	// enc := gob.NewEncoder(&network) // Will write to network.
	// err := enc.Encode(MapOfdb)
	// if err != nil {
	// 	println("error ")
	// 	log.Fatal("encode error:", err)
	// }
}

func InitMapOfDB() {
	LoadDB()
}

// func LoadData() {

// }
