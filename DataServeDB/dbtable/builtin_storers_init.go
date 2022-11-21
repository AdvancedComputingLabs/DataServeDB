package dbtable

import (
	dskstorv1 "DataServeDB/storers/dbtable_disk_store_v1"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	memstorv1 "DataServeDB/storers/dbtable_memstore_v1"
)

func init() {
	//TODO: error handling
	//idbstorer.RegisterStoreBasic(storerMemAppendLogName, NewMemAppendLog)
	idbstorer.RegisterStoreBasic(memstorv1.StorerMemV1Name, memstorv1.NewMemStoreV1)
	idbstorer.RegisterStoreBasic(dskstorv1.StorerDiskV1Name, dskstorv1.NewDiskStoreV1)
}
