package db

import (
	"DataServeDB/commtypes"
)

func (t *DB) TablesQueryGet(QryReqCtx *commtypes.QueryReqContext) (resultHttpStatus int, resultContent []byte, resultErr error) {

	table, err := t.GetTable(QryReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}

	return table.GetUser(QryReqCtx)

	//return 0, nil, nil
}
