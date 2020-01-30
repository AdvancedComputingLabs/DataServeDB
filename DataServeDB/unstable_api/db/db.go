package db

import (
	"DataServeDB/dbtable"
)

type MapOfTables map[string]dbtable.DbTable

// DB struct for maping
type DB struct {
	DbName       string
	DbInternalID int
	MapOfTables  MapOfTables
}
type MapOfDB map[string]DB
