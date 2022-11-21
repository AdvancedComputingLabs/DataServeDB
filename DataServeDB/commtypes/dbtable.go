package commtypes

import "DataServeDB/dbtypes"

type TableFieldStruct struct {
	FieldInternalId int
	FieldName       string
	FieldType       dbtypes.DbTypeI
	FieldTypeProps  dbtypes.DbTypePropertiesI
}

type FieldValueAndPropertiesHolder struct {
	V                  interface{}
	TableFieldInternal *TableFieldStruct
}

// TableRowWithFieldProperties used to pass on to storages with field property information.
type TableRowWithFieldProperties = map[int]FieldValueAndPropertiesHolder //int is field id

// TableRowByFieldsIds Used for internal operations. It contains row and its field by field ids which does not change.
// Should be used for internal table row storage and operations as it contains all the information. Although each storage can
// implement their own row type with additional properties.
type TableRowByFieldsIds = map[int]any

//NOTE: TableRow is by field names. It is external representation of the row, so it is in 'table_external_interface.go' file. Field names
// can change, so it should be only used for external transactions.

func (t *FieldValueAndPropertiesHolder) IsPk() bool {
	return t.TableFieldInternal.FieldTypeProps.IsPrimaryKey()
}

func (t *FieldValueAndPropertiesHolder) Name() string {
	return t.TableFieldInternal.FieldName
}

func (t *FieldValueAndPropertiesHolder) Value() any {
	return t.V
}
