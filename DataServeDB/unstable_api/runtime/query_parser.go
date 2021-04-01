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
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

func decodeJSONBody(w http.ResponseWriter, r *http.Request) (resultHttpStatus int, query *db.Query, err error) {
	// var dst interface{}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "content-Type header is not application/json"
			return http.StatusBadRequest, &db.Query{}, &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
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
func DecodeJSON(dst []byte) (resultHttpStatus int, query *db.Query, err error) {
	var result map[string]interface{}

	err = json.Unmarshal(dst, &result)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return http.StatusBadRequest, nil, errors.New(msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("request body contains badly-formed JSON")
			return http.StatusBadRequest, nil, errors.New(msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return http.StatusBadRequest, nil, errors.New(msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			return http.StatusBadRequest, nil, errors.New(msg)

		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			return http.StatusBadRequest, nil, errors.New(msg)

		case err.Error() == "http: request body too large":
			msg := "request body must not be larger than 1MB"
			return http.StatusBadRequest, nil, errors.New(msg)

		default:
			return http.StatusBadRequest, nil, err
		}
	}

	fieldRef := getFieldRef(string(dst))
	err = json.Unmarshal(dst, &result)
	if err != nil {
		//TODO: set http error
		resultHttpStatus = http.StatusNotAcceptable
		return
	}

	query = &db.Query{}

	for i, field := range fieldRef {
		value, ok := result[field]
		if ok {
			query.ItemLabel = field
			query.ItemValue = string(dst)
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
				Qry.ItemValue = string(data)
				if field == "$WHERE" {
					Qry.Rules, _ = getRule(string(data))
					Qry.Children = nil
				} else {
					qry, err := getUsersStuctFields(value, nxtRef)
					if err != nil {
						return query, err
					}
					Qry.Children = qry
				}
			} else {
				Qry.ItemValue = ""
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

func parseRules(str string) (rule db.Rule) {
	// rule := db.RuleInfo{}
	//   `(?m)(\(|)([A-z]*[.][A-z]*)( [A-Z]{2,4} )([A-z]*[.][A-z]*)(\)|)`
	// `(?m)((\(|)([A-z]*[.][A-z]*)( [A-Z]{2,4} )([A-z]*[.][A-z]*)(\)|))|(\(|)([A-z]*[.][A-z]*)( )([<>=]{1,2})( )(\w*)(\)|)`
	// "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
	var re = regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)|( [<>=A-Z]{1,4} )|(\w[\w]*)`)
	tbl := regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)`)
	optr := regexp.MustCompile(`(?m)([<>=A-Z]{1,4})`)
	oprnd := regexp.MustCompile(`(?m)(\w[\w]*)`)

	for i, match := range re.FindAllString(str, -1) {
		if tbl.FindString(match) != "" {
			if i == 0 {
				rule.LeftRule = getTableInfo(match)
			} else if i == 2 {
				rule.RightRule = getTableInfo(match)
			}
		} else if optr.FindString(match) != "" {
			rule.Operator = getOpr(match)
		} else if oprnd.FindString(match) != "" {
			if i == 0 {
				rule.LeftOperand = match
			} else if i == 2 {
				rule.RightOperand = match
			}
		}
	}
	return
}
func getOpr(str string) db.QueryOp {
	operators := map[string]db.QueryOp{
		"IS":  db.OpIS,
		"OR":  db.OpOR,
		"AND": db.OpAND,
		">":   db.OpGT,
		">=":  db.OpGTEQ,
		"<":   db.OpLT,
		"<=":  db.OpLTEQ,
	}
	var opre = regexp.MustCompile(`(?m)([<>=A-Z]{1,4})`)
	opr := opre.FindString(str)
	if v, ok := operators[opr]; ok {
		return v
	}
	return db.OpNone
}
func getTableInfo(str string) *db.RuleFieldInfo {
	tblInfo := &db.RuleFieldInfo{}
	arr := strings.Split(str, ".")
	tblInfo.TableName = arr[0]
	tblInfo.FieldName = arr[1]
	return tblInfo
}

func getRule(str string) (rules db.Rules, err error) {
	// rule := db.Rule{}
	frstGrpRule := regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)( [A-Z]{2,4} )([A-z]*[.][A-z]*)`)
	scndGrpRule := regexp.MustCompile(`([A-z]*[.][A-z]*)( )([<>=]{1,2})( )(\w*)`)
	bracegrp := regexp.MustCompile(`(?m)(\()(\w[\s.<>=\w]*)(\))`)
	optr := regexp.MustCompile(`(?m)([A-Z]{2,4})`)
	re := regexp.MustCompile(`(?m)(\()(\w[\s.<>=\w]*)(\))|([A-z]*[.][A-z]*)( [A-Z]{2,4} )([A-z]*[.][A-z]*)|(([A-z]*[.][A-z]*)( )([<>=]{1,2})( )(\w*))|([A-Z]{2,4})`)
	for _, match := range re.FindAllString(str, -1) {
		if b := bracegrp.FindStringIndex(match); b != nil {
			child, _ := getBrcGrp(match)
			rules = append(rules, child)
		} else if b := frstGrpRule.FindStringIndex(match); b != nil {
			rule := parseRules(match)
			rules = append(rules, rule)
		} else if b := scndGrpRule.FindStringIndex(match); b != nil {
			rule := parseRules(match)
			rules = append(rules, rule)
		} else if b := optr.FindStringIndex(match); b != nil {
			rules = append(rules, getOpr(match))
		}
	}
	return
}
func getBrcGrp(str string) (rules db.Rules, err error) {
	bracegrp := regexp.MustCompile(`(?m)(\w[\s.<>=\w]*)`)
	if match := bracegrp.FindString(str); match != "" {
		return getRule(match)
	}
	return
}
func getFirstGroup(str string) (rule db.Rule, err error) {
	frstGrpRule := regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)( [A-Z]{2,4} )([A-z]*[.][A-z]*)`)
	if b := frstGrpRule.FindStringIndex(str); b != nil {
		rule = parseRules(str)
	}

	return
}

// func getscndGrp(str string) (rule db.Rule, err error) {
// 	scndGrpRule := regexp.MustCompile(`([A-z]*[.][A-z]*)|([<>=]{1,2})|(\w*)`)
// 	tbl := regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)`)
// 	for _, v := range scndGrpRule.FindAllString(str, -1) {
// 		if byt := tbl.FindString(v); byt != "" {
// 			rule.MainTable = getTableInfo(v)
// 		}

// 	}
// 	return
// }
