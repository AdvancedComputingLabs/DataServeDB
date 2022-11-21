/*
	Storer NOTES:
 		1) indexName is used to specify the name of secondary index. Empty value is used for primary key.
		2) StoreBasic must be inherited by all extended storers with extend features.
		3) Non-implemented operations must return error of non-implemented.
		4) Stores can only be added at table creation time.
		5) At table creation time during adding, it should be checked for minimum supported operations for the table.

		### Initialization of storer implementation.
		- If they need to be registered, they should be in do their registration in init function in their implementation file.

		### Issues:
		- to the storer implementation, only file name should be passed or it needs path also? Path should be internally handled?
		-- Currently, it is passing path.

*/

package dbtable_interface_storer

import (
	"DataServeDB/commtypes"
	"DataServeDB/utils/rest/dberrors"
)

type TableStorerFeaturesType int

type TableOperationType int

// noinspection GoCommentStart
const (
	TableStorerFeatureNone TableStorerFeaturesType = iota

	// CRUD operations
	TableStorerFeature_Delete
	TableStorerFeature_Get
	TableStorerFeature_Insert
	TableStorerFeature_Update

	// constriants
	TableStorerConstraint_MustNotBeFirstStore
	TableStorerConstraint_MustBeLastStore
)

const (
	TableOperationNone TableOperationType = iota // used as default value which is 0 here

	TableOperationCreateTable
	TableOperationDeleteTable
	TableOperationListTables
	TableOperationGetTable

	TableOperationInsertRow
	TableOperationDeleteRow
	TableOperationGetRow
	TableOperationListRows
	TableOperationPatchRow
	TableOperationReplaceRow

	// not all table operations are used, but they are there if needed.
)

type TableInfo struct {
	TableName       string
	PrimaryKeyColId int
}

type StorerLoadResult struct {
	Data []byte
	Err  error
}

type StorerDeleteTableResult struct {
	RollbackItem any
	DbErr        *dberrors.DbError
}

type StorerInsertResult struct {
	Count int
	Data  any
	//Rollback bool //might not be needed. Error might serve as indicator to rollback or not. Or include flag to rollback.
	RollbackItem any
	DbErr        *dberrors.DbError
	//TODO: add skip flag to skip if it has been handled in the insert of the current storer
}

type StorerDeleteResult struct {
	Count        int
	Data         any
	RollbackItem any
	DbErr        *dberrors.DbError
}

type StorerUpdateResult struct {
	Count        int
	Data         any
	RollbackItem any
	UpdateType_  TableOperationType
	DbErr        *dberrors.DbError
}

type TableInfoGetter interface {
	// GetTableInfo TODO: this one has a potential to run into deadlock. May need to handle it.
	GetTableInfo() (*TableInfo, error)
}

// NewStoreBasic NOTE: can add argument of object (hash or/and regex patterns) to check against invalid names of file
//
//	but how useful it will be?
//
// TODO: check if 'StorerBasic' is pointer or it needs '*StorerBasic'
type NewStoreBasic func(tableId int, tableFolderPath string, tableInfo TableInfoGetter) (StorerBasic, *dberrors.DbError)

/*
	Get call and other calls that need to use index key:
		- currently uses indexColId to specify the index to use.
		-- is it better to specify index name instead of indexColId?
		-- ColId does not change but name can change. So ColId is better?
		- if indexColId is -1, it means primary key is used.
		- for 'key' string is better or any type?
		-- At the moment, I'll use same as Get call.
		- return numeric value for del, insert, update is 1 when successful and
			-1 when error. Should be -1 or 0 when error?

		- function calls must be transactional. If error occurs, changes are supposed to be reversed before exiting the function.
		-- one execption, row list order need not be preserved for performance reasons.
*/

type StorerBasic interface {
	DisplayName() string

	FeaturesAndConstraints(feature TableStorerFeaturesType) bool

	// CreateInit Initialization after creating of table, if it is needed. NewStoreBasic is used to create new store instance.
	// 	If storage implements table creation, it should be implemented here. All table information needed to create a table is passed
	// 	in NewStoreBasic function. Furthermore, 'GetTableInfo()' function can be used to get additional information.
	CreateInit(prevStore, nextStore StorerBasic) *dberrors.DbError
	LoadTable(input any, prevStore, nextStore StorerBasic) *StorerLoadResult
	DeleteTable(currentStoreIndex int, storesResults []*StorerDeleteTableResult) StorerDeleteTableResult
	DeleteTableRollback(currentStoreIndex int, storesResults []*StorerDeleteTableResult) *dberrors.DbError

	Delete(indexColId int, key any, currentStoreIndex int, storesResults []*StorerDeleteResult) StorerDeleteResult
	DeleteRollback(indexColId int, currentStoreIndex int, storesResults []*StorerDeleteResult) *dberrors.DbError
	Get(indexColId int, key any) (int, commtypes.TableRowByFieldsIds, *dberrors.DbError)
	Insert(rowWithProps commtypes.TableRowWithFieldProperties, currentStoreIndex int,
		storesResults []*StorerInsertResult) StorerInsertResult
	InsertRollback(currentStoreIndex int, storesResults []*StorerInsertResult) *dberrors.DbError
	Update(indexColId int, key any, rowWithProps commtypes.TableRowWithFieldProperties, updateType_ TableOperationType,
		currentStoreIndex int, storesResults []*StorerUpdateResult) StorerUpdateResult
	UpdateRollback(indexColId int, currentStoreIndex int, storesResults []*StorerUpdateResult) *dberrors.DbError

	GetNumberOfRows() int
}
