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
	UserNil bool
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

	targetName = "Tbl01"
	var qryReqCtx *commtypes.QueryReqContext
	qryReqCtx = commtypes.NewQueryReqContext(
		dbName, db, targetName, "")
	err, User, Flags := decodeJSONBody(w, r)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	if Flags.Id {
		if Flags.IdSpecd {
			qryReqCtx.TargetId = fmt.Sprint(User.UserId)
		}
	}
	resultHttpStatus, resultContent, resultErr = db.TablesQueryGet(qryReqCtx)
	println(string(resultContent))

	return
}
func decodeJSONBody(w http.ResponseWriter, r *http.Request) (err error, User users, Flags flags) {
	var dst interface{}
	var result map[string]interface{}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}, User, Flags
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, User, Flags

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, User, Flags

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, User, Flags

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, User, Flags

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, User, Flags

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}, User, Flags

		default:
			return err, User, Flags
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}, User, Flags
	}
	// err, Flgs := getUsersStuctFields(v)
	// if err != nil {
	// 	return err, users{}, flags{}
	// }
	// return nil, v, Flgs
	data, err := json.Marshal(dst)
	if err != nil {
		return err, User, Flags
	}

	// Unmarshal or Decode the JSON to the user struct.
	json.Unmarshal([]byte(data), &result)
	for f, v := range result {
		println(f, " --> ")
		switch f {
		case "Users":
			err, User, Flags = getUsersStuctFields(v)
			if err != nil {
				return err, User, Flags
			}
			break
			// return nil, User, Flags
		}

	}

	return nil, User, Flags
}

func getUsersStuctFields(User interface{}) (err error, UserStruct users, Flags flags) {
	var Fields map[string]interface{}
	Flags.UserNil = true

	data, err := json.Marshal(User)
	if err != nil {
		return err, UserStruct, Flags
	}
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(data, &Fields)

	for fld, val := range Fields {
		Flags.UserNil = false
		switch fld {
		case "UserId":
			Flags.Id = true
			data, err := json.Marshal(val)
			if err != nil {
				return err, users{}, flags{}
			}
			if string(data) != "{}" {
				Flags.IdSpecd = true
				// Unmarshal or Decode the JSON to the interface.
				json.Unmarshal(data, &UserStruct.UserId)
				println(UserStruct.UserId)
			}
		case "Name":
			Flags.Name = true
		case "properties":
			Flags.props = true
		}
	}
	return nil, UserStruct, Flags
}
