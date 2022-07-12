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
	"sync"

	"DataServeDB/dbstrcmp_base"
	"DataServeDB/dbtypes"
)

type tableFieldStruct struct {
	FieldInternalId int
	FieldName       string
	FieldType       dbtypes.DbTypeI
	FieldTypeProps  dbtypes.DbTypePropertiesI
}

type tableFieldsMetadataT struct {
	mu                             sync.RWMutex
	FieldInternalIdToFieldMetaData map[int]*tableFieldStruct
	FieldNameToFieldInternalId     map[string]int
}

type fieldValueAndPropertiesHolder struct {
	v                  interface{}
	tableFieldInternal *tableFieldStruct
}

// private static

func getNewInternalId(m map[int]*tableFieldStruct) int {
	i := 0
	for ; i < len(m); i++ {
		if _, exists := m[i]; !exists {
			return i
		}
	}
	return i
}

func newTableFieldProperties() *tableFieldStruct {
	fp := tableFieldStruct{FieldInternalId: -1}
	return &fp
}

// private

// tableFieldsMetadataT

func (t *tableFieldsMetadataT) add(fmd *tableFieldStruct, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	fieldNameKeyCase := fieldCaseHandler.Convert(fmd.FieldName)

	//check if field name exits
	if _, exists := t.FieldNameToFieldInternalId[fieldNameKeyCase]; exists {
		return errRplFieldNameAlreadyExist(fmd.FieldName)
	}

	//set internal id if it is -1
	if fmd.FieldInternalId == -1 {
		fmd.FieldInternalId = getNewInternalId(t.FieldInternalIdToFieldMetaData)
	}

	//check if internal id exits
	if _, exists := t.FieldInternalIdToFieldMetaData[fmd.FieldInternalId]; exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id already exists (this error shouldn't happen as field id is created internally)")
	}

	//add to both maps
	//TODO: Q: can internal id change after reload?
	t.FieldInternalIdToFieldMetaData[fmd.FieldInternalId] = fmd
	t.FieldNameToFieldInternalId[fieldNameKeyCase] = fmd.FieldInternalId

	return nil
}

// for future use, done for loading table but easier was to just store json creation text for now.
func (t *tableFieldsMetadataT) getCopyOfFieldsMetadataSafe() []tableFieldStruct {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []tableFieldStruct

	for _, v := range t.FieldInternalIdToFieldMetaData {
		result = append(result, *v)
	}

	return result
}

func (t *tableFieldsMetadataT) getFieldInternalId(fieldName, fieldNameKeyCase string) (int, error) {
	var fieldInternalId int
	var exists bool

	//check if field name exits
	if fieldInternalId, exists = t.FieldNameToFieldInternalId[fieldNameKeyCase]; !exists {
		return -1, errRplFieldDoesNotExist(fieldName)
	}

	return fieldInternalId, nil
}

// depricated
//TODO: only used in field tests; update tests and remove this method.
//NOTE: named it with internal just in case external interface requires to get field(s) meta data.
func (t *tableFieldsMetadataT) getFieldMetadataInternal(fieldName string, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) (*tableFieldStruct, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var fieldInternalId int
	var exists bool
	var existsErr error

	fieldNameKeyCase := fieldCaseHandler.Convert(fieldName)

	if fieldInternalId, existsErr = t.getFieldInternalId(fieldName, fieldNameKeyCase); existsErr != nil {
		return nil, existsErr
	}

	var fieldMetadata *tableFieldStruct

	if fieldMetadata, exists = t.FieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
	}

	return fieldMetadata, nil
}

func (t *tableFieldsMetadataT) getRowWithFieldMetadataInternal(userSentRow TableRow, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) (tableRowByInternalIdsWithFieldProperties, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var fieldInternalId int
	var exists bool
	var existsErr error

	rowByRowIdsWVAP := make(tableRowByInternalIdsWithFieldProperties)

	//goes through user provided fields and values
	for fieldName, v := range userSentRow {

		fieldNameKeyCase := fieldCaseHandler.Convert(fieldName)

		if fieldInternalId, existsErr = t.getFieldInternalId(fieldName, fieldNameKeyCase); existsErr != nil {
			return nil, existsErr
		}

		var fieldMetadata *tableFieldStruct

		if fieldMetadata, exists = t.FieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
			//TODO: Log.
			//TODO: Test if error or panic operations are atomic.
			//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
			panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
		}

		rowByRowIdsWVAP[fieldInternalId] = fieldValueAndPropertiesHolder{v: v, tableFieldInternal: fieldMetadata}
	}

	//check and fill missing fields
	for id, p := range t.FieldInternalIdToFieldMetaData {
		if _, exists := rowByRowIdsWVAP[id]; !exists {
			rowByRowIdsWVAP[id] = fieldValueAndPropertiesHolder{v: nil, tableFieldInternal: p}
		}
	}

	return rowByRowIdsWVAP, nil
}

// for future use, done for loading table but easier was to just store json creation text for now.
func (t *tableFieldsMetadataT) loadFieldsMetadataSafe(fields []tableFieldStruct) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.FieldInternalIdToFieldMetaData) > 0 || len(t.FieldNameToFieldInternalId) > 0 {
		//should be empty, if not error
		return errors.New("fields data is already loaded, it should be loaded once only")
	}

	for _, v := range fields {
		t.FieldNameToFieldInternalId[v.FieldName] = v.FieldInternalId
		t.FieldInternalIdToFieldMetaData[v.FieldInternalId] = &v
	}

	return nil
}

func (t *tableFieldsMetadataT) remove(fieldName string, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	fieldNameKeyCase := fieldCaseHandler.Convert(fieldName)

	fieldInternalId, existsErr := t.getFieldInternalId(fieldName, fieldNameKeyCase)
	if existsErr != nil {
		return existsErr
	}

	if _, exists := t.FieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
	}

	//delete both keys
	delete(t.FieldNameToFieldInternalId, fieldNameKeyCase)
	delete(t.FieldInternalIdToFieldMetaData, fieldInternalId)

	return nil
}

func (t *tableFieldsMetadataT) updateFieldName(fieldName string, newFieldName string, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var fieldInternalId int
	var existsErr error

	fieldNameKeyCase := fieldCaseHandler.Convert(fieldName)
	newFieldNameKeyCase := fieldCaseHandler.Convert(newFieldName)

	if fieldInternalId, existsErr = t.getFieldInternalId(fieldName, fieldNameKeyCase); existsErr != nil {
		return existsErr
	}

	if _, e := t.getFieldInternalId(newFieldName, newFieldNameKeyCase); e == nil {
		return errRplFieldNameAlreadyExist(newFieldName)
	}

	var exists bool
	var fieldMetadata *tableFieldStruct

	if fieldMetadata, exists = t.FieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
	}

	fieldMetadata.FieldName = newFieldName
	t.FieldNameToFieldInternalId[newFieldNameKeyCase] = fieldInternalId
	delete(t.FieldNameToFieldInternalId, fieldNameKeyCase)

	return nil
}

// tableFieldStruct

func (t *fieldValueAndPropertiesHolder) IsPk() bool {
	return t.tableFieldInternal.FieldTypeProps.IsPrimaryKey()
}

func (t *fieldValueAndPropertiesHolder) Name() string {
	return t.tableFieldInternal.FieldName
}

func (t *fieldValueAndPropertiesHolder) Value() any {
	return t.v
}
