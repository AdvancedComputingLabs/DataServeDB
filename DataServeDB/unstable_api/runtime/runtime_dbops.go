package runtime

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	db_rules "DataServeDB/dbsystem/rules"
	"DataServeDB/paths"
	"DataServeDB/unstable_api/db"
	//"DataServeDB/utils/mapwid"
	"DataServeDB/utils/mapwidgen"
)

//TODO:
// - CreateDb fn
// - DeleteDb fn
// - UnmountDb fn
// - RenameDb fn
// Note: right now goal is to get 1 db, tables, and their joined queries working.

//var databases *mapwid.MapWithId
var databases *mapwidgen.MapWithId[*db.DB]
var rwguardDbOps sync.RWMutex

func init() {
	databases = mapwidgen.New[*db.DB]()
}

func GetDb(dbName string) (*db.DB, error) {
	rwguardDbOps.RLock()
	defer rwguardDbOps.RUnlock()

	dbNameCasingHandled := syscasing(dbName)

	//TODO: log error; needs to be here?
	//TODO: error on log needs dbId?
	//TODO: better handle errors
	_, database, e := databases.GetByNameUnsync(dbNameCasingHandled)
	if e == nil {
		return database, nil
	}

	return nil, errors.New("database not found")
}

func loadDatabases() error {

	/*
		//NOTE: keeping it here, incase, in future databases metadata is required.
		// databases metadata has been removed for now.
		// reason: easier to just read the dirs and see if it exits - hy 6-Feb-2020
		 if b, err := storage.LoadDatabasesMetadata(); err != nil {
				return err
			} else {
				if err = json.Unmarshal(b, &DatabasesMetadata); err != nil {
					return err
				}
			}
	*/

	//TODO: can move to common function
	databases_dir := paths.GetDatabasesMainDirPath()

	if databases_dir == "" {
		//TODO: test should be fatal.
		return errors.New("databases dir does not exist")
	}

	//TODO: logging

	dirItems, e := ioutil.ReadDir(databases_dir)
	if e != nil {
		//TODO: refactor to dblog
		log.Fatal(e)
	}

	for _, dirItem := range dirItems {
		if dirItem.IsDir() && db_rules.DbNameIsValid(dirItem.Name()) {
			e = mountDb(dirItem.Name(), databases_dir)
			if e != nil {
				//TODO: log error
			}
		}
	}

	return nil
}

func CreateDb(dbName string) error {

	//check if name is valid.
	//check if database name exist in map and on file system.
	//create database

	//TODO: logging

	rwguardDbOps.Lock()
	defer rwguardDbOps.Unlock()

	if ok := db_rules.DbNameIsValid(dbName); !ok {
		return errors.New("invalid database name")
	}

	dbNameCasingHandled := syscasing(dbName)

	if exists := databases.HasNameUnsync(dbNameCasingHandled); exists {
		return errors.New("database name already exists")
	}
	//TODO: can move to common function
	databases_dir := paths.GetDatabasesMainDirPath()

	if databases_dir == "" {
		//TODO: test should be fatal.
		return errors.New("databases dir does not exist")
	}

	dirItems, e := ioutil.ReadDir(databases_dir)
	if e != nil {
		//TODO: logging
		return errors.New("error reading databases from file system")
	}

	for _, dirItem := range dirItems {
		if dirItem.IsDir() {
			if strings.EqualFold(dirItem.Name(), dbNameCasingHandled) {
				return errors.New("database name already exists")
			}
		}
	}

	//create db
	// 1. create db dir

	return nil
}

func mountDb(dbName, dbPath string) error {
	rwguardDbOps.Lock()
	defer rwguardDbOps.Unlock()

	//convert name to system casing for mapping.
	dbNameCaseHandled := syscasing(dbName)

	//check if db is already in map
	if _, _, e := databases.GetByNameUnsync(dbNameCaseHandled); e.Error() != "name does not exist" { //TODO: should be enum and not string
		//TODO: properly handle this, currently it assumes it is just name already exists.
		return errors.New("database name already exists")
	}

	//initalize db
	//NOTE: actual casing is needed for db dir/file io
	database, e := new(db.DB).Init(dbName, dbPath)
	if e != nil {
		return e
	}

	//add to db map

	//TODO: I see issues. This is runtime id, permanant guid would be better?

	dbId := databases.GetLastIdUnsync()
	triesCount := 0

	for true {
		triesCount++
		dbId++

		if triesCount > 5 {
			//TODO: handle this better.
			break
		}

		//TODO: properly handle this, currently it assumes it is just id already exists.
		//TODO: need to check error type too.
		if e := databases.AddUnsync(dbId, dbNameCaseHandled, database); e == nil {
			break
		}
	}

	return nil
}
