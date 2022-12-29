// Copyright (c) 2020 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"DataServeDB/commtypes"
	"DataServeDB/dbsystem/constants"
	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
	"DataServeDB/utils/rest/dberrors"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func parseDbAuthStr(authStr string) (scheme, authToken string, e error) {
	toks := strings.SplitN(authStr, " ", 3)

	if len(toks) < 2 {
		e = errors.New("malformed db auth token")
		return
	}

	scheme = toks[0]
	authToken = toks[1]

	return
}

func getDbAuthFromHttpHeader(r *http.Request) (scheme, authToken string, e error) {

	authStrs, authExits := r.Header["DbAuth"]
	if authExits {
		return parseDbAuthStr(authStrs[0])
	}

	authStrs, authExits = r.Header["Authorization"]
	if !authExits {
		return "", "", errors.New("authentication is required")
	}

	return parseDbAuthStr(authStrs[0])
}

// File Resr Path Handler
func FileRestPathHandler(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName string, pathLevels []dbrouter.PathLevel) {

	//TODO: dbName empty test case

	var resultHttpStatus int
	var resultContent []byte
	var resultErr error

	db, err := GetDb(dbName)
	if err != nil {
		resultHttpStatus, resultContent, resultErr = rest.HttpRestDbError(dberrors.NewDbError(dberrors.DatabaseNotFound, err))
		rest.ResponseWriteHelper(w, resultHttpStatus, resultContent, resultErr)
		return
	}
	var dbReqCtx *commtypes.DbReqContext

	dbReqCtx = commtypes.NewDbReqContext(
		httpMethod, resPath, matchedPath,
		dbName, db, pathLevels)

	switch strings.ToUpper(httpMethod) {
	case "GET":
		dbReqCtx.RestMethodId = constants.RestMethodGet
		//db.getFile(dbReqXtx)
		resultHttpStatus, resultContent, resultErr = db.FilesGet(dbReqCtx)
	case "POST":
		fmt.Println("hai-----")
		err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dbReqCtx.RestMethodId = constants.RestMethodPost
		resultHttpStatus, resultContent, resultErr = db.FilesPost(dbReqCtx, r.MultipartForm)
	case "DELETE":
		dbReqCtx.RestMethodId = constants.RestMethodDelete
		resultHttpStatus, resultContent, resultErr = db.FilesDelete(dbReqCtx)
	case "PUT":
		err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dbReqCtx.RestMethodId = constants.RestMethodPut
		resultHttpStatus, resultContent, resultErr = db.FilesPutorPatch(dbReqCtx, r.MultipartForm)
	case "PATCH":
		err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dbReqCtx.RestMethodId = constants.RestMethodPatch
		resultHttpStatus, resultContent, resultErr = db.FilesPutorPatch(dbReqCtx, r.MultipartForm)
	default:
		dbReqCtx.RestMethodId = constants.RestMethodNone
		resultHttpStatus, resultContent, resultErr = rest.HttpRestDbError(
			dberrors.NewDbError(dberrors.InvalidInputHttpMethodNotSupported,
				fmt.Errorf("http method '%s' is not supported", httpMethod)))
	}

	rest.ResponseWriteHelper(w, resultHttpStatus, resultContent, resultErr)
}

func TableRestPathHandler(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName string, pathLevels []dbrouter.PathLevel) {
	//TODO: dbName empty test case

	var resultHttpStatus int
	var resultContent []byte
	var resultErr error

	db, err := GetDb(dbName)
	if err != nil {
		resultHttpStatus, resultContent, resultErr = rest.HttpRestDbError(dberrors.NewDbError(dberrors.DatabaseNotFound, err))
		rest.ResponseWriteHelper(w, resultHttpStatus, resultContent, resultErr)
		return
	}

	//TODO: atomic check if db is close.

	//TODO: increament requent count

	//TODO: check if requires auth

	//db auth
	//TODO: check if separating db auth scheme is better. Keeping all auth tokens under one header have problems with top browsers?

	//scheme, authToken, errAuthHeader := getDbAuthFromHttpHeader(r)
	//
	////NOTE: currently auth is must
	//if errAuthHeader != nil {
	//	return http.StatusForbidden, nil, errAuthHeader
	//}
	//
	//if errAuth := AuthUser(scheme, authToken); errAuth != nil {
	//	return http.StatusForbidden, nil, errAuth
	//}

	/*
		Examples:

		CREATE TABLE:
			resPath: /re_db/tables
			matchedPath: /re_db/tables
			dbName: re_db
			targetName: tables
			targetDbResTypeId: constants.DbResTypes_Table

		GET/SELECT ROW:
			resPath: /re_db/tables/users/Id:1 or /re_db/tables/users/1
			matchedPath: re_db/tables/users
			dbName: re_db
			targetName: users
			targetDbResTypeId: 1
	*/

	var dbReqCtx *commtypes.DbReqContext

	dbReqCtx = commtypes.NewDbReqContext(
		httpMethod, resPath, matchedPath,
		dbName, db, pathLevels)

	switch strings.ToUpper(httpMethod) {
	case "GET":
		dbReqCtx.RestMethodId = constants.RestMethodGet
		resultHttpStatus, resultContent, resultErr = db.TablesGet(dbReqCtx)
	case "POST":
		dbReqCtx.RestMethodId = constants.RestMethodPost
		resultHttpStatus, resultContent, resultErr = db.TablesPost(dbReqCtx, r.Body)
	case "DELETE":
		dbReqCtx.RestMethodId = constants.RestMethodDelete
		resultHttpStatus, resultContent, resultErr = db.TablesDelete(dbReqCtx)
	case "PUT":
		dbReqCtx.RestMethodId = constants.RestMethodPut
		resultHttpStatus, resultContent, resultErr = db.TablesPut(dbReqCtx, r.Body)
	case "PATCH":
		dbReqCtx.RestMethodId = constants.RestMethodPatch
		resultHttpStatus, resultContent, resultErr = db.TablesPatch(dbReqCtx, r.Body)
	default:
		dbReqCtx.RestMethodId = constants.RestMethodNone
		resultHttpStatus, resultContent, resultErr = rest.HttpRestDbError(
			dberrors.NewDbError(dberrors.InvalidInputHttpMethodNotSupported,
				fmt.Errorf("http method '%s' is not supported", httpMethod)))
	}

	rest.ResponseWriteHelper(w, resultHttpStatus, resultContent, resultErr)
}

func commonHttpServReqHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(w)

	if r.Method == "OPTIONS" {
		return
	}

	//********
	//	path := r.URL.String()
	//
	//	table, key, err := requestParser(path)
	//	if err != nil {
	//		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//	if strings.ToUpper(table) == "SIGNIN" {
	//		result, err := Signin(w, r)
	//		if err != nil {
	//			return
	//		}
	//		w.Write(result)
	//		return
	//	} else if strings.ToUpper(table) == "AUTHTOKEN" {
	//		_, err := AuthenticateToken(r)
	//		if err != nil {
	//			w.WriteHeader(http.StatusUnauthorized)
	//			return
	//		}
	//		w.WriteHeader(http.StatusOK)
	//		return
	//	}
	//
	//	/***************************************************************************/
	//	/* session cookie checking
	//	/*******************************************************************************/
	//	// We can obtain the session token from the requests cookies, which come with every request
	//
	//	claimID, err := AuthenticateToken(r)
	//	if err != nil {
	//		w.WriteHeader(http.StatusUnauthorized)
	//		return
	//	}
	//	// after check get data according to authenticated user

	//TODO: where does http restful api user authentications go?
	//TODO: http restful api return results handling.

	reqPath := rest.HttpRestPathParser(r.URL.String())
	dbrouter.MatchPathAndCallHandler(w, r, reqPath, r.Method)
}
