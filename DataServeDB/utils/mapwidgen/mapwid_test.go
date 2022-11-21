package mapwidgen

import "testing"

func mwidGen01Test(t *testing.T) {
	createTable01JSON := `{
	  "TableName": "Tbl01",
	  "TableColumns": [
		"Id int32 PrimaryKey",
		"UserName string Length:5..50 !Nullable",
		"Counter int32 default:Increment(1,1) !Nullable",
		"DateAdded datetime default:Now() !Nullable",
		"GlobalId guid default:NewGuid() !Nullable"
	  ]
	}`

	_ = createTable01JSON
}
