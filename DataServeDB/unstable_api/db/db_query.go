package db

import (
	"DataServeDB/commtypes"
	"net/http"
)

func (t *DB) TablesQueryGet(QryReqCtx *commtypes.QueryReqContext) (resultHttpStatus int, resultContent []byte, resultErr error) {
	var res string
	table, err := t.GetTable(QryReqCtx.TargetName)
	if err != nil {
		resultErr = err
		return
	}
	if QryReqCtx.TargetId != "" {
		res, resultErr = table.GetRowByPrimaryKeyReturnsJSON(QryReqCtx.TargetId)
	} else {
		res, resultErr = table.GetTableRows()
	}

	return http.StatusOK, []byte(res), resultErr

	//return 0, nil, nil
}
