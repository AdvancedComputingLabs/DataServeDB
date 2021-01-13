// Copyright (c) 2020 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package runtime

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"DataServeDB/dbsystem"
	db_rules "DataServeDB/dbsystem/rules"
	"DataServeDB/paths"
	"DataServeDB/unstable_api/db"
	"DataServeDB/unstable_api/dbrouter"
)

//TODO: change to get functions

var syscasing = dbsystem.SystemCasingHandler.Convert

var isInitalized = false

// mapOfDatabases is exported
var mapOfDatabases = make(map[string]*db.DB)
var rwguard sync.RWMutex

func GetDb(dbName string) (*db.DB, error) {
	//TODO: syncing
	dbNameCasingHandled := syscasing(dbName)
	if data, ok := mapOfDatabases[dbNameCasingHandled]; ok {
		return data, nil
	}
	return nil, errors.New("database not found")
}

func mountDb(dbName string, dbPath string) error {
	rwguard.Lock()
	defer rwguard.Unlock()

	//convert name to system casing for mapping.
	dbNameCaseHandled := syscasing(dbName)

	//check if db is already in map
	if _, alreadyExists := mapOfDatabases[dbNameCaseHandled]; alreadyExists {
		return errors.New("database already mounted")
	}

	//initalize db
	//NOTE: actual casing is needed for db dir/file io
	database, e := new(db.DB).Init(dbName, dbPath)
	if e != nil {
		return e
	}

	//add to db map
	mapOfDatabases[dbNameCaseHandled] = database

	return nil
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

func IsInitalized() bool {
	return isInitalized
}

func Start() error {
	//TODO: check if this needs go process to independently initalize db server; there could be hanging issue?
	fmt.Println("Starting DataServeDB server ...")

	//TODO: list/log all the databases being mounted.
	//TODO: refactor

	//TODO: error handling
	loadDatabases()

	//routing
	dbrouter.Register("{db_name}/tables/{tbl_name}", TableRestPathHandler)

	//http server and rest api routing
	StartHttpServer()

	cliProcessor()

	//log.Println("Closing DataServeDB server ...")

	isInitalized = true

	return nil
}
