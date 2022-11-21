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

	"DataServeDB/commtypes"
	"DataServeDB/dbtypes"
)

// TODO: add tag annotations.

type tableMain struct {
	TableId             int
	TableName           string
	PkPos               int
	TableFieldsMetaData tableFieldsMetadataT

	//TableStringComparer dbstrcmp_base.DbStrCmpInterface
	//TableDataContainersIds map[string]int
}

// TODO: is this needed here?
type tablesMapper struct {
	tableIdToTable     map[int]tableMain
	tableNameToTableId map[string]int
}

// NOTE: Only creates table object, kept this way for unit testing.
func newTableMain(tableInternalId int, tableName string) *tableMain {

	//if tableInternalId < 0 {
	//	//TODO: table id needs to be done.
	//	//TODO: validation of tableid or/and auto generation of correct one
	//	tableInternalId = 0
	//}

	t := tableMain{
		TableId:   tableInternalId,
		TableName: tableName,
		PkPos:     0, //default position is zero
		TableFieldsMetaData: tableFieldsMetadataT{
			mu:                             sync.RWMutex{},
			FieldInternalIdToFieldMetaData: make(map[int]*commtypes.TableFieldStruct),
			FieldNameToFieldInternalId:     make(map[string]int),
		},
		//TableStringComparer: simpleFold,
	}

	return &t
}

func (t *tableMain) getPkType() (dbtypes.DbTypeI, dbtypes.DbTypePropertiesI) {
	pkFieldInternal := t.TableFieldsMetaData.FieldInternalIdToFieldMetaData[t.PkPos]
	return pkFieldInternal.FieldType, pkFieldInternal.FieldTypeProps
}
