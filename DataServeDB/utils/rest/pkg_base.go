package rest

import (
	"errors"
	"net/http"
	"strings"

	"DataServeDB/utils/rest/dberrors"
)

func HttpRestError(errCode dberrors.RestError, err error) (resultHttpStatus int, resultContent []byte, resultErr error) {
	resultErr = errors.New(errCode.Error() + ": " + err.Error())
	resultHttpStatus = dberrors.GetStatusCode(errCode)
	return
}

func HttpRestDbError(dbErr *dberrors.DbError) (resultHttpStatus int, resultContent []byte, resultErr error) {
	resultErr = errors.New(dbErr.ToError().Error())
	resultHttpStatus = dberrors.GetStatusCode(dbErr.ErrCode)
	return
}

func HttpRestOk(resultContent []byte) (resultHttpStatus int, resultContentOut []byte, resultErr error) {
	resultHttpStatus = 200
	resultContentOut = resultContent
	return
}

func HttpRestOkNoContent() (resultHttpStatus int, resultContent []byte, resultErr error) {
	resultHttpStatus = 204
	return
}

func HttpRestPathParser(httpReqPath string) string {
	//does not include host name, eg: /re_db/tables/users/Id:1
	//NOTE: kept for future if path needs some parsing, then api doesn't need to change.
	return httpReqPath
}

func ParseKeyValue(resPath string) (key string, value string, err error) {
	pos := strings.LastIndex(resPath, "/") + 1

	if pos == 0 {
		err = errors.New("key path is in wrong format")
		return
	}

	if pos >= len(resPath) {
		err = errors.New("key or value is not provided")
		return
	}

	splitted := strings.SplitN(resPath[pos:], ":", 2)

	if len(splitted) == 1 {
		value = splitted[0]
	} else {
		key = splitted[0]
		value = splitted[1]
	}

	return
}

func ResponseWriteHelper(w http.ResponseWriter, httpStatus int, content []byte, err error) {
	// TODO: alternatively function can be passed instead of rest result params, which will save some code duplication.
	//  But need to check performance impact. - HY 19-Nov-2022.

	if err != nil {
		content = []byte(err.Error())
	}

	w.WriteHeader(httpStatus)
	if _, err2 := w.Write(content); err2 != nil {
		//TODO: check cases which can cause this error; how int should be handled?
		// TODO: log error
	}
}
