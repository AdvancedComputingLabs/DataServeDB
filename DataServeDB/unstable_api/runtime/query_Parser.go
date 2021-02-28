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
	"DataServeDB/unstable_api/db"
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

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request) (query db.Query, err error) {
	var dst interface{}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return db.Query{}, &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
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
			return db.Query{}, &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return db.Query{}, &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return db.Query{}, &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return db.Query{}, &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return db.Query{}, &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return db.Query{}, &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return db.Query{}, err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return db.Query{}, &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}
	return decodeJSON(dst)
	// return err, query
}
func decodeJSON(dst interface{}) (query db.Query, err error) {
	var result map[string]interface{}
	data, err := json.Marshal(dst)
	if err != nil {
		return query, err
	}
	json.Unmarshal([]byte(data), &result)
	for f, v := range result {
		// switch f {
		// case "Users":
		query.ItemLabel = f
		query.ItemValue = data
		children, err := getUsersStuctFields(v)
		if err != nil {
			return db.Query{}, err
		}
		query.Children = append(query.Children, children...)
		// }
	}
	// Unmarshal or Decode the JSON to the user struct.

	return query, nil
}
func getUsersStuctFields(dst interface{}) (query []db.Query, err error) {
	var result map[string]interface{}
	var resArray []interface{}
	data, err := json.Marshal(dst)
	if err != nil {
		return query, err
	}
	if _, ok := dst.(map[string]interface{}); ok {
		// Unmarshal or Decode the JSON to the user struct.
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			return query, err
		}
		return getStruct(result)
	} else if _, ok := dst.([]interface{}); ok {
		err = json.Unmarshal([]byte(data), &resArray)
		if err != nil {
			return
		}
		return getArrayStruct(resArray)
	} else if _, ok := dst.(string); ok {
		return nil, nil
	}
	// Unmarshal or Decode the JSON to the user struct.

	return nil, nil
}
func getStruct(dst map[string]interface{}) (query []db.Query, err error) {
	for f, v := range dst {
		var Qry db.Query = db.Query{}
		Qry.ItemLabel = f

		data, err := json.Marshal(v)
		if err != nil {
			return query, err
		}
		if string(data) != "{}" && string(data) != "[{}]" {
			Qry.ItemValue = data
			qry, err := getUsersStuctFields(v)
			if err != nil {
				return query, err
			}
			Qry.Children = qry
		} else {
			Qry.ItemValue = nil
			Qry.Children = nil
		}
		query = append(query, Qry)
	}

	return query, err
}
func getArrayStruct(dst []interface{}) (query []db.Query, err error) {
	// "Properties": [
	//   {
	//     "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
	//   }
	// ]
	for _, v := range dst {

		var Qry db.Query = db.Query{}
		qry, err := getUsersStuctFields(v)
		if err != nil {
			return query, err
		}
		Qry.Children = qry
		query = append(query, Qry)
	}
	return
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
