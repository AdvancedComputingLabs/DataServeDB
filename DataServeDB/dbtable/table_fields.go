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
	"sync"

	"DataServeDB/dbstrcmp_base"
	"DataServeDB/dbtypes"
)

type tableFieldProperties struct {
	FieldInternalId int
	FieldName       string
	FieldType       dbtypes.DbTypeI
	FieldTypeProps  dbtypes.DbTypePropertiesI
}

type tableFieldsMetadataT struct {
	mu                             sync.RWMutex
	fieldInternalIdToFieldMetaData map[int]*tableFieldProperties
	fieldNameToFieldInternalId     map[string]int
}

// private static

func getNewInternalId(m map[int]*tableFieldProperties) int {
	i := 0
	for ; i < len(m); i++ {
		if _, exists := m[i]; !exists {
			return i
		}
	}
	return i
}

func newTableFieldProperties() *tableFieldProperties {
	fp := tableFieldProperties{FieldInternalId: -1}
	return &fp
}

// private

// tableFieldsMetadataT

func (t *tableFieldsMetadataT) add(fmd *tableFieldProperties, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	fieldNameKeyCase := fieldCaseHandler.Convert(fmd.FieldName)

	//check if field name exits
	if _, exists := t.fieldNameToFieldInternalId[fieldNameKeyCase]; exists {
		return errRplFieldNameAlreadyExist(fmd.FieldName)
	}

	//set internal id if it is -1
	if fmd.FieldInternalId == -1 {
		fmd.FieldInternalId = getNewInternalId(t.fieldInternalIdToFieldMetaData)
	}

	//check if internal id exits
	if _, exists := t.fieldInternalIdToFieldMetaData[fmd.FieldInternalId]; exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id already exists (this error shouldn't happen as field id is created internally)")
	}

	//add to both maps
	t.fieldInternalIdToFieldMetaData[fmd.FieldInternalId] = fmd
	t.fieldNameToFieldInternalId[fieldNameKeyCase] = fmd.FieldInternalId

	return nil
}

func (t *tableFieldsMetadataT) getFieldInternalId(fieldName, fieldNameKeyCase string) (int, error) {
	var fieldInternalId int
	var exists bool

	//check if field name exits
	if fieldInternalId, exists = t.fieldNameToFieldInternalId[fieldNameKeyCase]; !exists {
		return -1, errRplFieldDoesNotExist(fieldName)
	}

	return fieldInternalId, nil
}

//NOTE: named it with internal just in case external interface requires to get field(s) meta data.
func (t *tableFieldsMetadataT) getFieldMetadataInternal(fieldName string, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) (*tableFieldProperties, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var fieldInternalId int
	var exists bool
	var existsErr error

	fieldNameKeyCase := fieldCaseHandler.Convert(fieldName)

	if fieldInternalId, existsErr = t.getFieldInternalId(fieldName, fieldNameKeyCase); existsErr != nil {
		return nil, existsErr
	}

	var fieldMetadata *tableFieldProperties

	if fieldMetadata, exists = t.fieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
	}

	return fieldMetadata, nil
}

func (t *tableFieldsMetadataT) remove(fieldName string, fieldCaseHandler dbstrcmp_base.DbStrCmpInterface) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	fieldNameKeyCase := fieldCaseHandler.Convert(fieldName)

	fieldInternalId, existsErr := t.getFieldInternalId(fieldName, fieldNameKeyCase)
	if existsErr != nil {
		return existsErr
	}

	if _, exists := t.fieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
	}

	//delete both keys
	delete(t.fieldNameToFieldInternalId, fieldNameKeyCase)
	delete(t.fieldInternalIdToFieldMetaData, fieldInternalId)

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
	var fieldMetadata *tableFieldProperties

	if fieldMetadata, exists = t.fieldInternalIdToFieldMetaData[fieldInternalId]; !exists {
		//TODO: Log.
		//TODO: Test if error or panic operations are atomic.
		//NOTE: Reason for panic: if this error occurred then there is bug in the code, fix.
		panic("internal field id exists but field's metadata doesn't (this error shouldn't happen if code is correct)")
	}

	fieldMetadata.FieldName = newFieldName
	t.fieldNameToFieldInternalId[newFieldNameKeyCase] = fieldInternalId
	delete(t.fieldNameToFieldInternalId, fieldNameKeyCase)

	return nil
}

// tableFieldProperties
