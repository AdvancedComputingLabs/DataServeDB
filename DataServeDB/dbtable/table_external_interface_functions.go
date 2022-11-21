package dbtable

import (
	"DataServeDB/commtypes"
	"DataServeDB/utils/rest/dberrors"
)

func HelloFrom() (success_msg []byte, dberr *dberrors.DbError) {
	return []byte("Hello World!"), nil
}

func callTableFunction(t *DbTable, ctx *commtypes.DbReqContext, key string) ([]byte, *dberrors.DbError) {

	//function names are case sensitive?

	//below will do for now.

	switch key {
	case "HelloFrom":
		return HelloFrom()
	}

	//TODO: nil, nil is ok?
	return nil, nil
}
