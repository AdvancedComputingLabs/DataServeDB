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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request) (resultHttpStatus int, query db.Query, err error) {
	// var dst interface{}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return http.StatusBadRequest, db.Query{}, &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	resultHttpStatus, query, err = DecodeJSON(data)
	fmt.Println(query)
	if err != nil {
		return
	}
	return
}
func DecodeJSON(dst []byte) (resultHttpStatus int, query db.Query, err error) {
	var result map[string]interface{}
	resultHttpStatus, err = checkJsonData(dst)
	if err != nil {
		return
	}
	fieldRef := getFieldRef(string(dst))
	err = json.Unmarshal(dst, &result)
	if err != nil {
		return
	}

	for i, field := range fieldRef {
		value, ok := result[field]
		if ok {
			query.ItemLabel = field
			query.ItemValue = dst
			children, err := getUsersStuctFields(value, fieldRef[i+1:])
			if err != nil {
				return http.StatusForbidden, query, err
			}
			query.Children = append(query.Children, children...)
		}
	}
	resultHttpStatus = http.StatusOK
	return
}
func getUsersStuctFields(dst interface{}, fieldRef []string) (query []db.Query, err error) {
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
		return getStruct(result, fieldRef)
	} else if _, ok := dst.([]interface{}); ok {
		err = json.Unmarshal([]byte(data), &resArray)
		if err != nil {
			return
		}
		return getArrayStruct(resArray, fieldRef)
	}

	return nil, nil
}
func getStruct(dst map[string]interface{}, fieldRef []string) (query []db.Query, err error) {
	for i, field := range fieldRef {
		nxtRef := fieldRef[i+1:]
		value, ok := dst[field]
		if ok {
			var Qry db.Query = db.Query{}
			Qry.ItemLabel = field

			data, err := json.Marshal(value)
			if err != nil {
				return query, err
			}
			if string(data) != "{}" && string(data) != "[{}]" {
				Qry.ItemValue = data
				qry, err := getUsersStuctFields(value, nxtRef)
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
	}

	return query, err
}
func getArrayStruct(dst []interface{}, fieldRef []string) (query []db.Query, err error) {
	return getUsersStuctFields(dst[0], fieldRef)
}
func checkJsonData(str []byte) (resultHttpStatus int, err error) {
	var dst interface{}
	dec := json.NewDecoder(strings.NewReader(string(str)))
	dec.DisallowUnknownFields()

	err = dec.Decode(&dst)
	if err != nil {
		resultHttpStatus = http.StatusNotAcceptable
		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		resultHttpStatus = http.StatusNotAcceptable
		return
	}
	// data, err := json.Marshal(dst)
	resultHttpStatus = http.StatusOK
	return
}
func getFieldRef(dst string) (fieldRef []string) {
	dec := json.NewDecoder(strings.NewReader(dst))
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if d, ok := t.(string); ok {
			fieldRef = append(fieldRef, d)
		}
	}
	return
}
