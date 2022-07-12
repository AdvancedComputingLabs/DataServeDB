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

package db

import (
	"encoding/json"
	"fmt"

	"DataServeDB/dbsystem"
	"DataServeDB/dbsystem/dbstorage"
	"DataServeDB/dbtable"
	"DataServeDB/paths"
	"DataServeDB/utils/mapwidgen"
)

//TODO: can somethings can be private?

//TODO: renaming some structs and fields?

const dbMetadataJson = "metadata/db.json" //TODO: this should be here or single file paths locations?

//type TablesMapOld map[string]*dbtable.DbTable

// DB struct for maping
type DB struct {
	DbName       string //runtime at the moment
	dbPath       string //runtime at the moment
	DbInternalID int    //runtime at the moment
	//Tables       TablesMap
	Tables *mapwidgen.MapWithId[*dbtable.DbTable] //loaded at runtime, but not sure TableId is persistant.
}

type DatabaseMetaSaveStructure struct {
	Tables []dbtable.DatabaseTableRecreation
}

var syscasing = dbsystem.SystemCasingHandler.Convert

func (me *DB) DbPath() string {
	return me.dbPath
}

func (me *DB) GetId() int {
	return me.DbInternalID
}

func (me *DB) loadDbMetadata() error {
	dbMetaPath := paths.Combine(me.dbPath, dbMetadataJson)

	dbMetaBuf, e := dbstorage.LoadFromDisk(dbMetaPath)
	if e != nil {
		return e
	}

	var dbMeta DatabaseMetaSaveStructure

	if e := json.Unmarshal(dbMetaBuf, &dbMeta); e != nil {
		//TODO: handle error better
		return e
	}

	//load tables
	for _, stbl := range dbMeta.Tables {

		//NOTE: table name validation happens 'LoadFromTableSaveStructure' so it is not needed here.

		stbl.DbTableCreationStructure.AssignDb(me)
		dtbl, e := dbtable.LoadFromTableSaveStructure(stbl) //could be called attach table data
		if e != nil {
			//TODO: handle error better
			continue
		}

		//no need to check if table already exits in the map, since this is at load time.
		//me.Tables[syscasing(dtbl.TblMain.TableName)] = dtbl

		//TODO: here table loaded from disk is added to the map, error can happen with add to map operation.
		me.Tables.AddUnsync(dtbl.TblMain.TableId, syscasing(dtbl.TblMain.TableName), dtbl)
	}

	//at the moment there are only tables

	return nil
}

func (me *DB) Init(dbName, dbsPath string) (*DB, error) {

	me.DbName = dbName

	//me.Tables = make(TablesMap)
	me.Tables = mapwidgen.New[*dbtable.DbTable]()

	me.dbPath = paths.ConstructDbPath(dbName, dbsPath)

	// ! for dev time stuff
	//fmt.Println(dbMetaPath)
	//me.createDbMetadata()
	//me.createDbMetadataGeneric()

	e := me.loadDbMetadata()
	if e != nil {
		return nil, e
	}

	return me, nil
}

func (me *DB) getTablesSaveStructureJson() string {
	//todo: synching

	var result string
	dbMeta := DatabaseMetaSaveStructure{}

	//for _, tblI := range me.Tables.GetItemsUnsync() {
	//	if tbl, ok := tblI.(*dbtable.DbTable); ok {
	//		tblStructure := dbtable.GetTableStorageStructure(tbl)
	//		dbMeta.Tables = append(dbMeta.Tables, tblStructure)
	//	}
	//}

	for _, tbl := range me.Tables.GetItemsUnsync() {
		tblStructure := dbtable.GetTableStorageStructure(tbl)
		dbMeta.Tables = append(dbMeta.Tables, tblStructure)
	}

	{
		r, e := json.Marshal(dbMeta)
		if e != nil {
			//TODO: handle error better
			panic(e)
		}
		result = string(r)
	}

	return result
}

func (me *DB) updateDbMetadataOnDisk() error {
	//TODO: perhaps it requires some optimization.

	tablesStructureJson := me.getTablesSaveStructureJson()
	dbMetaPath := paths.Combine(me.dbPath, dbMetadataJson)
	e := dbstorage.SaveToDisk([]byte(tablesStructureJson), dbMetaPath)
	if e != nil {
		return e
	}
	return nil
}

//mocks for dev only
func (me *DB) createDbMetadata() {

	createTable01JSON := `{
	  "TableName": "Tbl01",
	  "TableFields": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable",
		"DateAdded datetime default:Now() !Nullable",
		"GlobalId guid default:NewGuid() !Nullable"
	  ]
	}`

	createTable02JSON := `{
	  "TableName": "Tbl02",
	  "TableFields": [
		"PropertyId int32 PrimaryKey",
		"PropertName string Length:5..50 !Nullable"
	  ]
	}`

	tbl01, e := dbtable.CreateTableJSON(createTable01JSON, me)
	_ = e
	//me.Tables[tbl01.TblMain.TableName] = tbl01
	tbl01.TblMain.TableId = 0
	me.Tables.AddUnsync(tbl01.TblMain.TableId, tbl01.TblMain.TableName, tbl01)

	tbl02, e := dbtable.CreateTableJSON(createTable02JSON, me)
	_ = e
	//me.Tables[tbl02.TblMain.TableName] = tbl02
	tbl02.TblMain.TableId = 0
	me.Tables.AddUnsync(tbl02.TblMain.TableId, tbl02.TblMain.TableName, tbl02)

	//TODO: make it save all the tables
	//TODO: change map to table pointer
	//TODO: using tableid instead of name
	//does dbtable stores more than 1 table?

	e = me.updateDbMetadataOnDisk()
	if e != nil {
		fmt.Println(e)
	}

}

func (me *DB) createDbMetadataGeneric() {
	createTable01JSON := `{
	  "TableName": "Tbl01",
	  "TableFields": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable",
		"DateAdded datetime default:Now() !Nullable",
		"GlobalId guid default:NewGuid() !Nullable"
	  ]
	}`

	tbl01, e := dbtable.CreateTableJSON(createTable01JSON, me)
	_ = e
	tbl01.TblMain.TableId = 0

	fmt.Println("Created table id:", tbl01.TblMain.TableId, "Created table name:", tbl01.TblMain.TableName)

	tables := mapwidgen.New[*dbtable.DbTable]()

	tables.AddUnsync(tbl01.TblMain.TableId, tbl01.TblMain.TableName, tbl01)

	//tblNameCaseHandled := syscasing(tbl01.TblMain.TableName)

	tableId, tbl01Getted, _ := tables.GetByNameUnsync(tbl01.TblMain.TableName)
	fmt.Println("Getted table returned id:", tableId)
	fmt.Println("Getted table id:", tbl01Getted.TblMain.TableId, "Getted table name:", tbl01Getted.TblMain.TableName)
}
