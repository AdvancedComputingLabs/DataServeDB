package dbtypes

import (
	"fmt"

	"DataServeDB/dbsystem"
)

// Section: declarations

var syscasing = dbsystem.SystemCasingHandler.Convert

var dbtypes_map = map[string]DbTypeI{}

func addDbTypeToMap(dbtype DbTypeI) {
	cased_type_name := syscasing(dbtype.GetDbTypeDisplayName())
	dbtypes_map[cased_type_name] = dbtype
}

func getDbType(dbtype_name string) (DbTypeI, error) {
	cased_type_name := syscasing(dbtype_name)
	if dt, ok := dbtypes_map[cased_type_name]; ok {
		return dt, nil
	}
	return nil, fmt.Errorf("variable type '%s' doesn't exist", dbtype_name)
}

func init() {
	addDbTypeToMap(Bool)
	addDbTypeToMap(DateTime)
	addDbTypeToMap(Int32)
	addDbTypeToMap(String)
}
