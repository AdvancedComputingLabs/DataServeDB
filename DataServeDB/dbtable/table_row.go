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

import "DataServeDB/dbstrcmp_base"

type tableRowByInternalIdsWithFieldProperties = map[int]fieldValueAndPropertiesHolder
type tableRowByInternalIds = map[int]interface{}

func fromLabeledByFieldNames(row TableRow, tbl *tableMain, fieldCasingHandler dbstrcmp_base.DbStrCmpInterface) (tableRowByInternalIds, error) {
	//TODO: Primary key should be first. Or do this rearrangement on client side?

	var meta = &tbl.TableFieldsMetaData
	rowById := make(tableRowByInternalIds)

	tmp, e := meta.getRowWithFieldMetadataInternal(row, fieldCasingHandler)
	if e != nil {
		return nil, e
	}

	//execute
	for Id, holder := range tmp {
		if vConverted, errConversion := holder.tableFieldInternal.FieldType.ConvertValue(holder.v, holder.tableFieldInternal.FieldTypeProps); errConversion == nil {
			rowById[Id] = vConverted
		} else {
			return nil, errRplRowDataConversion(holder.tableFieldInternal.FieldName, errConversion)
		}
	}

	return rowById, nil
}

//NOTE: tableRowByInternalIds is a map, so no need to pass it as pointer
func toLabeledByFieldNames(row tableRowByInternalIds, tbl *tableMain) (TableRow, error) {
	//TODO: Primary key should be first.

	//- check if field exists; covered.
	//- check if missing field; not covered.

	var meta *tableFieldsMetadataT = &tbl.TableFieldsMetaData
	rowByNames := make(TableRow)

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
