// Copyright (c) 2019 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package dbtable

import (
	"errors"
	"fmt"

	"DataServeDB/commtypes"
	"DataServeDB/dbsystem"
	db_rules "DataServeDB/dbsystem/rules"
	"DataServeDB/dbtypes"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/utils/rest/dberrors"
)

//TODO: move it to error messages (single location)

// public functions

// private functions

func validateTableName(tableName string) error {
	if !db_rules.TableNameIsValid(tableName) {
		return fmt.Errorf("invalid table name '%s'", tableName)
	}
	return nil
}

func validateFieldMetaData(fieldCreationText string, pkIsSet *bool) (*commtypes.TableFieldStruct, error) {

	fp := newTableFieldProperties()

	//QUESTION: does makes sense to make parse field in this pkg?

	if tableFieldName, tableFieldDbType, tableFieldDbTypeProperties, e := dbtypes.ParseFieldProperties(fieldCreationText); e != nil {
		return nil, e
	} else {
		fp.FieldName = tableFieldName
		fp.FieldType = tableFieldDbType
		fp.FieldTypeProps = tableFieldDbTypeProperties
	}

	if fp.FieldTypeProps.IsPrimaryKey() {
		if *pkIsSet {
			return nil, errors.New("table can only have one primary key")
		}
		*pkIsSet = true
	}

	return fp, nil
}

// validates and creates main object, reasons:
// 1) better code, since adding fields automatically checks certain constraints.
// 2) optimization, since most of the time validation is followed by creation.
// - HY 26-Dec-2019
func validateCreateTableMetaData(tableInternalId int, createTableData *createTableExternalStruct) (*tableMain, error) {
	//first quick checks

	if e := validateTableName(createTableData.TableName); e != nil {
		return nil, e
	}

	//quick checks end

	pkIsSet := false
	dbTbl := newTableMain(tableInternalId, createTableData.TableName)

	for _, fieldCreationText := range createTableData.TableColumns {
		//_ = i

		var fp *commtypes.TableFieldStruct
		var e error

		if fp, e = validateFieldMetaData(fieldCreationText, &pkIsSet); e != nil {
			return nil, e
		}

		if e = dbTbl.TableFieldsMetaData.add(fp, dbsystem.SystemCasingHandler); e != nil {
			return nil, e
		}

		if fp.FieldTypeProps.IsPrimaryKey() {
			dbTbl.PkPos = fp.FieldInternalId
		}

		//NOTE: db type property validation is done during parsing.
	}

	if !pkIsSet {
		return nil, errors.New("table must have primary key")
	}

	return dbTbl, nil
}

// NOTE: TableRow is a map, so no need to pass it as pointer
// WARNING: TableRow (by field name) is not returned unless function succeeds. So don't override r in calling function.
func validateRowData(t *tableMain, r TableRow, tableOp idbstorer.TableOperationType) (TableRow, tableRowByInternalIds, *dberrors.DbError) {
	rowByInternalId, dberr := fromLabeledByFieldNames(r, t, dbsystem.SystemCasingHandler, tableOp)
	if dberr != nil {
		return nil, nil, dberr
	}

	rowConvertedWithCorrectTypes, dberr := toLabeledByFieldNames(rowByInternalId, t)
	if dberr != nil {
		return nil, nil, dberr
	}

	return rowConvertedWithCorrectTypes, rowByInternalId, nil
}
