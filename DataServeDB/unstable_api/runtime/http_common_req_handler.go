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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strings"

	"github.com/golang/gddo/httputil/header"

	"DataServeDB/commtypes"
	"DataServeDB/dbsystem/constants"
	"DataServeDB/unstable_api/dbrouter"
)

const maxMEMORY = 1 * 1024 * 1024

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

func QueryRestPathHandler(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName, targetName string, targetDbResTypeId constants.DbResTypes) (resultHttpStatus int, resultContent []byte, resultErr error) {
	//TODO: resPath if it is more than /query needs to be handled appropriately.
	// var dst interface{}
	println("decode json")
	decodeJSONBody(w, r)

	return
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
		return db.TablesGet(dbReqCtx)
	case "POST":
		if err := r.ParseMultipartForm(maxMEMORY); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		form := getDataInsert(r.MultipartForm.Value)
		for key := range form {
			if httpStatus, err := db.TablesPost(dbReqCtx, form.Get(key)); err != nil {
				return httpStatus, nil, err
			}
		}
	case "PUT":
		if err := r.ParseMultipartForm(maxMEMORY); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		form := getDataInsert(r.MultipartForm.Value)
		for key := range form {
			if httpStatus, err := db.TablesEdit(dbReqCtx, form.Get(key)); err != nil {
				return httpStatus, nil, err
			}
		}
	case "DELETE":
		db.TablesDelete(dbReqCtx)
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
	resultHttpStatus, resultContent, resultErr := dbrouter.MatchPathAndCallHandler(w, r, reqPath, r.Method)

	if resultErr == nil && resultHttpStatus == http.StatusOK {
		w.Write(resultContent)
	}

	return
}

type malformedRequest struct {
	status int
	msg    string
}
type users struct {
	UserId     int
	Name       string
	Properties []string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request) error {
	var dst interface{}
	var result map[string]users
	var Fields map[string]interface{}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	// println(dst)
	// t := reflect.ValueOf(dst)
	// i := t.NumField()
	// for i = 0; i < t.NumField(); i++ {
	// 	fmt.Printf("%v\n", t.Field(i))
	// }
	data, err := json.Marshal(dst)
	if err != nil {
		return err
	}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(data), &result)
	for f, v := range result {
		println(f, " --> ")
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		json.Unmarshal(data, &Fields)

		for fld, val := range Fields {
			println(fld, "--")
			va, err := json.Marshal(val)
			if err != nil {
				return err
			}
			println(string(va))

		}
	}

	return nil
}
func getDataInsert(form url.Values) url.Values {
	return form
}
