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
func (t *DB) TablesPost(dbReqCtx *commtypes.DbReqContext, insertDataJson string) (resultHttpStatus int, resultErr error) {
	table, err := t.GetTable(dbReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}
	resultErr = table.InsertRowJSON(insertDataJson)
	if resultErr == nil {
		resultHttpStatus = http.StatusOK
	}

	return
}
func (t *DB) TablesEdit(dbReqCtx *commtypes.DbReqContext, insertDataJson string) (resultHttpStatus int, resultErr error) {

	table, err := t.GetTable(dbReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}
	
	resultErr = table.EditRowJSON(insertDataJson)
	if resultErr == nil {
		resultHttpStatus = http.StatusOK
	}

	return
}
