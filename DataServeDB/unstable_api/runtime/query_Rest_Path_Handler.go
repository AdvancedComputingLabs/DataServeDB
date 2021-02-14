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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"strings"

	"github.com/golang/gddo/httputil/header"
)

type malformedRequest struct {
	status int
	msg    string
}
type users struct {
	UserId     int
	Name       string
	Properties []string
}

type flags struct {
	Id      bool
	IdSpecd bool
	Name    bool
	props   bool
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}
func QueryRestPathHandler(w http.ResponseWriter, r *http.Request, httpMethod, resPath, matchedPath, dbName, targetName string, targetDbResTypeId constants.DbResTypes) (resultHttpStatus int, resultContent []byte, resultErr error) {
	//TODO: resPath if it is more than /query needs to be handled appropriately.
	// var dst interface{}
	db, e := GetDb(dbName)
	if e != nil {
		//at the moment it only returns database not found.
		return http.StatusNotFound, nil, e
	}

	var dbReqCtx *commtypes.QueryReqContext
	dbReqCtx = commtypes.NewQueryReqContext(
		httpMethod, resPath, matchedPath,
		dbName, db, targetName, targetDbResTypeId)

	dbReqCtx.RestMethodId = constants.RestMethodGet

	err, _, Flags := decodeJSONBody(w, r)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	db.TablesQueryGet(dbReqCtx)
	if Flags.Id {
		if Flags.IdSpecd {
			// dbReqCtx.TargetName =
			// db.TablesGet()
		}
	}

	return
}
func decodeJSONBody(w http.ResponseWriter, r *http.Request) (error, users, flags) {
	var dst interface{}
	var result map[string]users

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}, users{}, flags{}
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
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, users{}, flags{}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, users{}, flags{}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, users{}, flags{}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, users{}, flags{}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, users{}, flags{}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}, users{}, flags{}

		default:
			return err, users{}, flags{}
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}, users{}, flags{}
	}

	data, err := json.Marshal(dst)
	if err != nil {
		return err, users{}, flags{}
	}

	// Unmarshal or Decode the JSON to the user struct.
	json.Unmarshal([]byte(data), &result)
	for f, v := range result {
		println(f, " --> ")
		switch f {
		case "Users":
			err, Flgs := getUsersStuctFields(v)
			if err != nil {
				return err, users{}, flags{}
			}
			return nil, v, Flgs
		}

	}

	return nil, users{}, flags{}
}

func getUsersStuctFields(Users users) (error, flags) {
	var Fields map[string]interface{}
	var Flags flags

	data, err := json.Marshal(Users)
	if err != nil {
		return err, flags{}
	}
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(data, &Fields)

	for fld, val := range Fields {
		switch fld {
		case "UserId":
			if val != "" {
				Flags.IdSpecd = true
			}
		case "Name":
			Flags.Name = true
		case "properties":
			Flags.props = true
		}
		// va, err := json.Marshal(val)
		// if err != nil {
		// 	return err, users{}
		// }
		// println(string(va))
	}
	return nil, Flags
}
