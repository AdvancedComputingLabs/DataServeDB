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
	//TableStringComparer dbstrcmp_base.DbStrCmpInterface
	//TableDataContainersIds map[string]int
}

// Separate, makes it easier to save it separately than table metadata.
type tableDataContainer struct {
	Rows []tableRowByInternalIds
	PkToRowMapper map[interface{}]int64
}

type tablesMapper struct {
	tableIdToTable map[int]tableMain
	tableNameToTableId map[string]int
}

//Only creates table object, kept this way for unit testing.
func newTableMain(tableName string) *tableMain {

	t := tableMain{
		TableId:   0,
		TableName: tableName,
		TableFieldsMetaData: tableFieldsMetadataT{
			mu:                             sync.RWMutex{},
			fieldInternalIdToFieldMetaData: make(map[int]*tableFieldProperties),
			fieldNameToFieldInternalId:     make(map[string]int),
		},
		//TableStringComparer: simpleFold,
	}

	return &t
}

func (tm *tableMain) getPkType() dbtypes.DbTypeInterface {
	// TODO: if pk position is not always zero then change it find pk.
	// it should be implemented as always zero.
	pkProperties := tm.TableFieldsMetaData.fieldInternalIdToFieldMetaData[0]
	return pkProperties.FieldType
}
