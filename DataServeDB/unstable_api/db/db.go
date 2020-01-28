package db

import (
	"DataServeDB/dbtable"
)

type mapOfTables map[string]dbtable.DbTable

// DB struct for maping
type DB struct {
	dbName       string
	dbInternalID int
	MapOfTables  mapOfTables
}
