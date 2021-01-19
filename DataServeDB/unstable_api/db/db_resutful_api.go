package db

import (
	"DataServeDB/commtypes"
	"net/http"
)

func (t *DB) TablesGet(dbReqCtx *commtypes.DbReqContext) (resultHttpStatus int, resultContent []byte, resultErr error) {

	table, err := t.GetTable(dbReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}

	return table.Get(dbReqCtx)

	//return 0, nil, nil
}
func (t *DB) TablesPost(dbReqCtx *commtypes.DbReqContext) (resultHttpStatus int, resultErr error) {
	table, err := t.GetTable(dbReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}
	resultErr = table.InsertRowJSON(dbReqCtx.DataInsert)
	if resultErr == nil {
		resultHttpStatus = http.StatusOK
	}

	return
}
