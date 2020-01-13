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

type tableRowByInternalIds = map[int]interface{}

func fromLabeledByFieldNames(row TableRow, tbl *tableMain, fieldCasingHandler dbstrcmp_base.DbStrCmpInterface) (tableRowByInternalIds, error) {
	//TODO: Primary key should be first.

	//- check if field exists; covered.
	//- check if missing field; not covered.

	var meta = &tbl.TableFieldsMetaData
	rowById := make(tableRowByInternalIds)

	for k, v := range row {
		if field, err := meta.getFieldMetadataInternal(k, fieldCasingHandler); err == nil {
			//TODO: validation by constraints.
			if vConverted, errConversion := field.FieldType.ConvertValue(v, field.FieldTypeProps, true); errConversion == nil {
				rowById[field.FieldInternalId] = vConverted
			} else {
				return nil, errRplRowDataConversion(field.FieldName, errConversion)
			}
		} else {
			return nil, err
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
		if fieldProps, exits := meta.fieldInternalIdToFieldMetaData[k]; exits {
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
