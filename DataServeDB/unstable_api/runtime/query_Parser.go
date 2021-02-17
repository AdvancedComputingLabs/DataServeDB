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

	"strings"

	"github.com/golang/gddo/httputil/header"
)

type malformedRequest struct {
	status int
	msg    string
}

// type users struct {
// 	UserId     int
// 	Name       string
// 	Properties []string
// }

// type flags struct {
// 	UserNil bool
// 	Id      bool
// 	IdSpecd bool
// 	Name    bool
// 	props   bool
// }
type query struct {
	itemLabel string
	itemType  string
	itemValue []byte // json Converted
	children  []query
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request) (err error, Query query) {
	var dst interface{}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}, query{}
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
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, query{}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, query{}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, query{}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, query{}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}, query{}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}, query{}

		default:
			return err, query{}
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}, query{}
	}
	err, Query = decodeJSON(dst)
	return err, Query
}
func decodeJSON(dst interface{}) (err error, Query query) {
	var result map[string]interface{}
	data, err := json.Marshal(dst)
	if err != nil {
		return err, Query
	}
	json.Unmarshal([]byte(data), &result)
	for f, v := range result {
		switch f {
		case "Users":
			Query.itemLabel = f
			Query.itemValue = data
			err, children := getUsersStuctFields(v)
			if err != nil {
				return err, query{}
			}
			Query.children = append(Query.children, children...)
		}
	}
	// Unmarshal or Decode the JSON to the user struct.

	return nil, Query
}
func getUsersStuctFields(dst interface{}) (err error, Query []query) {
	var result map[string]interface{}
	data, err := json.Marshal(dst)
	if err != nil {
		return err, Query
	}

	// Unmarshal or Decode the JSON to the user struct.
	json.Unmarshal([]byte(data), &result)
	for f, v := range result {
		var Qry query = query{}
		Qry.itemLabel = f

		data, err := json.Marshal(v)
		if err != nil {
			return err, Query
		}
		if string(data) != "{}" && string(data) != "[{}]" {
			Qry.itemValue = data
			err, qry := getUsersStuctFields(v)
			if err != nil {
				return err, Query
			}
			Qry.children = qry
			if Qry.children != nil {
				Qry.itemType = "struct"
			}
		} else {
			Qry.itemValue = nil
			Qry.children = nil
		}
		Query = append(Query, Qry)
	}

	return nil, Query
}

// func getUsersStuctFields(User interface{}) (err error, UserStruct users, Flags flags) {
// 	var Fields map[string]interface{}
// 	Flags.UserNil = true

// 	data, err := json.Marshal(User)
// 	if err != nil {
// 		return err, UserStruct, Flags
// 	}
// 	// Unmarshal or Decode the JSON to the interface.
// 	json.Unmarshal(data, &Fields)

// 	for fld, val := range Fields {
// 		Flags.UserNil = false
// 		switch fld {
// 		case "UserId":
// 			Flags.Id = true
// 			data, err := json.Marshal(val)
// 			if err != nil {
// 				return err, users{}, flags{}
// 			}
// 			if string(data) != "{}" {
// 				Flags.IdSpecd = true
// 				// Unmarshal or Decode the JSON to the interface.
// 				json.Unmarshal(data, &UserStruct.UserId)
// 				println(UserStruct.UserId)
// 			}
// 		case "Name":
// 			Flags.Name = true
// 		case "properties":
// 			Flags.props = true
// 		}
// 	}
// 	return nil, UserStruct, Flags
// }
