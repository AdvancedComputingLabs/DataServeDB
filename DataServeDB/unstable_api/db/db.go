package db

import (
	"DataServeDB/dbtable"
)

//TODO: can somethings can be private?

//TODO: renaming some structs and fields?

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
