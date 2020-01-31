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

type Meta map[string]DbMeta
type DbMeta struct {
	TableMeta []TableMeta
}
type TableMeta struct {
	Table string
}
