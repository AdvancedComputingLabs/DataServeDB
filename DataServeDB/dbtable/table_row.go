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
	"DataServeDB/commtypes"
	"DataServeDB/dbstrcmp_base"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/utils/rest/dberrors"
	"fmt"
)

// type tableRowByInternalIdsWithFieldProperties = map[int]fieldValueAndPropertiesHolder
type tableRowByInternalIds = commtypes.TableRowByFieldsIds

func fromLabeledByFieldNames(row TableRow, tbl *tableMain, fieldCasingHandler dbstrcmp_base.DbStrCmpInterface, tableOp idbstorer.TableOperationType) (tableRowByInternalIds, *dberrors.DbError) {
	//TODO: Primary key should be first. Or do this rearrangement on client side?

	var meta = &tbl.TableFieldsMetaData
	rowById := make(tableRowByInternalIds)

	tmp, dberr := meta.getRowWithFieldMetadataInternal(row, fieldCasingHandler, tableOp)
	if dberr != nil {
		return nil, dberr
	}

	//execute
	for fieldId, holder := range tmp {
		if vConverted, errConversion := holder.TableFieldInternal.FieldType.ConvertValue(holder.V, holder.TableFieldInternal.FieldTypeProps); errConversion == nil {
			rowById[fieldId] = vConverted
		} else {
			// TODO: looks like bad way to insert column name into error message; is there better way and safer way?
			return nil, dberrors.NewDbError(dberrors.InvalidInput, fmt.Errorf(errConversion.Error(), holder.Name()))
		}
	}

	return rowById, nil
}

// NOTE: tableRowByInternalIds is a map, so no need to pass it as pointer
func toLabeledByFieldNames(row tableRowByInternalIds, tbl *tableMain) (TableRow, *dberrors.DbError) {
	//TODO: Primary key should be first.

	//TODO: race condition? metadata needs to be locked?

	//- check if field exists; covered.
	//- check if missing field; not covered.

	var meta *tableFieldsMetadataT = &tbl.TableFieldsMetaData
	rowByNames := make(TableRow)

	meta.mu.RLock()
	defer meta.mu.RUnlock()

	for k, v := range row {
		if fieldProps, exits := meta.FieldInternalIdToFieldMetaData[k]; exits {
			//NOTE: this does not need conversion because this will probably come from db internal row and will have correct type.
			//But check/test.
			rowByNames[fieldProps.FieldName] = v //NOTE: used field name stored.
		} else {
			//TODO: Log.
			//TODO: Test if error or panic operations are atomic.
			//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
			panic("this is internal error, shouldn't happen")
			//return nil, errors.New("this is internal error, shouldn't happen")
		}
	}

	return rowByNames, nil
}
