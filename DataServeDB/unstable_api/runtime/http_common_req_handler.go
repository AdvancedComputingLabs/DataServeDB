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
	"DataServeDB/commtypes"
	"DataServeDB/dbsystem/constants"
	"errors"
	"net/http"
	"strings"

	"DataServeDB/unstable_api/dbrouter"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func httpRestPathParser(httpReqPath string) string {
	//does not include host name, eg: /re_db/tables/users/Id:1
	//NOTE: kept for future if path needs some parsing, then api doesn't need to change.
	return httpReqPath
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

func TableRestPathHandler(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName, targetName string, targetDbResTypeId constants.DbResTypes) (resultHttpStatus int, resultContent []byte, resultErr error) {
	//TODO: dbName empty case

	db, e := GetDb(dbName)

	if e != nil {
		//at the moment it only returns database not found.
		return http.StatusNotFound, nil, e
	}

	//TODO: atomic check if db is close.

	//TODO: increament requent count

	//TODO: check if requires auth

	//db auth
	//TODO: check if separating db auth scheme is better. Keeping all auth tokens under one header have problems with top browsers?

	scheme, authToken, errAuthHeader := getDbAuthFromHttpHeader(r)

	//NOTE: currently auth is must
	if errAuthHeader != nil {
		return http.StatusForbidden, nil, errAuthHeader
	}

	if errAuth := AuthUser(scheme, authToken); errAuth != nil {
		return http.StatusForbidden, nil, errAuth
	}

	/*
		Example:
			resPath: /re_db/tables/users/Id:1
			matchedPath: re_db/tables/users
			dbName: re_db
			targetName: users
			targetDbResTypeId: 1
	 */

	var dbReqCtx *commtypes.DbReqContext
	dbReqCtx = commtypes.NewDbReqContext(
			httpMethod, resPath, matchedPath,
			dbName, db, targetName, targetDbResTypeId)

	switch strings.ToUpper(httpMethod) {
	case "GET":
		dbReqCtx.RestMethodId = constants.RestMethodGet
	}

	return
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

	reqPath := httpRestPathParser(r.URL.String())
	dbrouter.MatchPathAndCallHandler(w, r, reqPath, r.Method)

	return
}
