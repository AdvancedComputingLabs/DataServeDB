package runtime

import (
	"encoding/json"
	"errors"
	"fmt"

	storage "DataServeDB/dbsystem/dbstorage"
	"DataServeDB/dbtable"
	"DataServeDB/unstable_api/db"
)

//TODO: change to get functions

// MapOfdb is exported
var MapOfdb = make(map[string]db.DB)
var MetaData = make(db.Meta)

func GetDB(dbName string) (db.DB, error) {
	if data, ok := MapOfdb[dbName]; ok {
		return data, nil
	}
	return db.DB{}, errors.New("db Name not found")
}

func LoadDB() error {
	//TODO: refactor
	if b, err := storage.LoadMata(); err != nil {
		return err
	} else {
		if err = json.Unmarshal(b, &MetaData); err != nil {
			return err
		}
	}

	var mapoftable = make(db.MapOfTables)

	for dbName, DbMeta := range MetaData {
		for _, tableMeta := range DbMeta.TableMeta {
			if tbl01, err := dbtable.CreateTableJSON(tableMeta.Table); err == nil {
				//TODO: why it is using GetSaveLoadStructure? it shouldn't be need to load.
				if jsonStr, err := dbtable.GetSaveLoadStructure(tbl01); err == nil {
					fmt.Println(jsonStr)
					// println(tableMeta.Table, dbName)
					//TODO: continued: maybe because of this function.
					if tbl, err := dbtable.LoadTableFromDB(jsonStr, dbName); err == nil {
						fmt.Printf("table loaded :- %v\n", tbl)
						mapoftable["users"] = *tbl
						MapOfdb["re_db"] = db.DB{
							DbName:       "re_db", //TODO: hard coded.
							DbInternalID: 0, //TODO: 0?
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
		}
	}

	return nil
}

func InitMapOfDB() {
	LoadDB()
}

//TODO: testing related?
func CreateDBmeta() {
	table := db.TableMeta{
		Table: `{"TableName": "users","TableRoot": "re_db", "TableFields": ["Id int32 PrimaryKey","UserName string Length:5..50 !Nullable","Counter int32 default:Increment(1,1) !Nullable","DateAdded datetime default:Now() !Nullable","GlobalId guid default:NewGuid() !Nullable"]}`,
	}

	dbMeta := db.DbMeta{}

	dbMeta.TableMeta = append(dbMeta.TableMeta, table)

	// dbMeta
	MetaData["re_db"] = dbMeta
	if db, err := json.Marshal(MetaData); err != nil {

	} else {
		storage.SaveToDisk(db)
	}
}
