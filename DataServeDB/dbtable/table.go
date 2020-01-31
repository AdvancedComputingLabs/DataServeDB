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

	"DataServeDB/dbtypes"
	//"DataServeDB/dbstrcmp_base"
)

// TODO: add tag annotations.

type tableMain struct {
	TableId             int
	TableName           string
	TableFieldsMetaData tableFieldsMetadataT
	//TableRoot to Get to Know this table belongs to which DB
	TableRoot string //TODO: references db, can be made better or is it really needed?
	//TableStringComparer dbstrcmp_base.DbStrCmpInterface
	//TableDataContainersIds map[string]int
}

// Separate, makes it easier to save it separately than table metadata.
type tableDataContainer struct {
	Rows          []tableRowByInternalIds
	PkToRowMapper map[interface{}]int64
}

type tablesMapper struct {
	tableIdToTable     map[int]tableMain
	tableNameToTableId map[string]int
}

//Only creates table object, kept this way for unit testing.
func newTableMain(tableInternalId int, tableName string) *tableMain {

	if tableInternalId <= -1 {
		//TODO: table id needs to be done.
		//TODO: validation of tableid or/and auto generation of correct one
		tableInternalId = 0
	}

	t := tableMain{
		TableId:   tableInternalId,
		TableName: tableName,
		TableFieldsMetaData: tableFieldsMetadataT{
			mu:                             sync.RWMutex{},
			FieldInternalIdToFieldMetaData: make(map[int]*tableFieldStruct),
			FieldNameToFieldInternalId:     make(map[string]int),
		},
		//TableStringComparer: simpleFold,
	}

	return &t
}

func (tm *tableMain) getPkType() (dbtypes.DbTypeI, dbtypes.DbTypePropertiesI) {
	// TODO: if pk position is not always zero then change it find pk.
	// it should be implemented as always zero.
	pkFieldInternal := tm.TableFieldsMetaData.FieldInternalIdToFieldMetaData[0]
	return pkFieldInternal.FieldType, pkFieldInternal.FieldTypeProps
}
