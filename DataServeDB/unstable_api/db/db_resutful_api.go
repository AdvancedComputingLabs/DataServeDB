package db

import (
	"DataServeDB/commtypes"
)

func (me *DB) TablesGet(dbReqCtx *commtypes.DbReqContext) (resultHttpStatus int, resultContent []byte, resultErr error) {

	table, err := me.GetTable(dbReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}

	return table.Get(dbReqCtx)

	//return 0, nil, nil
}
