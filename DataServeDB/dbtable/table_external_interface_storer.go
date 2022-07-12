package dbtable

/*
	Storer NOTES:
 		1) indexName is used to specify the name of secondary index. Empty value is used for primary key.
		2) StoreBasic must be inherited by all extended storers with extend features.
		3) Non-implemented operations must return error of non-implemented.
		4) Stores can only be added at table creation time.
		5) At table creation time during adding, it should be checked for minimum supported operations for the table.

		TODO: body io.ReadCloser may or may not be idea. Need to see how it works out in practice.

*/

type TableStorerFeatures int

const (
	TableStorerFeatureNone = iota
	TableStorerFeatureDelete
	TableStorerFeatureGet
	TableStorerFeatureInsert
	TableStorerFeaturepUpdate
)

type TableRowWithFieldProperties tableRowByInternalIdsWithFieldProperties

type StorerBasic interface {
	Implemented(feature TableStorerFeatures) bool

	Delete(indexName string, key string) (int, error)
	Get(indexName string, key string) (int, any, error)
	Insert(rowWithProps TableRowWithFieldProperties, data any) (int, error) //data any is used in case storage needs data any other the row data
	Update(rowWithProps TableRowWithFieldProperties, data any) (int, error)
}

func GetTableRowWithFieldProperties(table *DbTable, rowByInternalIds tableRowByInternalIds) (TableRowWithFieldProperties, error) {
	var meta = &table.TblMain.TableFieldsMetaData
	rowWithProps := TableRowWithFieldProperties{}

	meta.mu.RLock()
	defer meta.mu.RUnlock()

	for k, v := range rowByInternalIds {
		if fieldProps, exits := meta.FieldInternalIdToFieldMetaData[k]; exits {
			//NOTE: this does not need conversion because this will probably come from db internal row and will have correct type.
			//But check/test.
			rowWithProps[k] = fieldValueAndPropertiesHolder{v: v, tableFieldInternal: fieldProps} //NOTE: used field name stored.
		} else {
			//TODO: Log.
			//TODO: Test if error or panic operations are atomic.
			//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
			panic("this is internal error, shouldn't happen")
			//return nil, errors.New("this is internal error, shouldn't happen")
		}
	}

	return rowWithProps, nil
}
