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
	"bytes"
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

func decodeJSONBody(w http.ResponseWriter, r *http.Request) (resultHttpStatus int, query []db.Query, err error) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "content-Type header is not application/json"
			return http.StatusBadRequest, nil, &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	resultHttpStatus, query, err = DecodeJSON(data)
	if err != nil {
		return
	}
	return
}
func DecodeJSON(dst []byte) (resultHttpStatus int, query []db.Query, err error) {
	var result map[string]interface{}
	fmt.Println(string(dst))

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

	err = json.Unmarshal(dst, &result)
	if err != nil {
		resultHttpStatus = http.StatusNotAcceptable
		return
	}

	fieldRef := getFieldRef(dst)
	query, err = getUsersStuctFields(result, fieldRef)
	if err != nil {
		resultHttpStatus = http.StatusNotAcceptable
		return
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
	// re := regexp.MustCompile(`(?m)(\[(\s*|)\{(\s*|))([\D\d]*)((\s*|)\}(\s*|)\])`)
	refmap := map[string]bool{}
	for i, field := range fieldRef {
		nxtRef := fieldRef[i+1:]
		if _, ok := refmap[field]; !ok {
			if value, ok := dst[field]; ok {
				var Qry db.Query = db.Query{}
				Qry.ItemLabel = field

				data, err := json.Marshal(value)
				if err != nil {
					return query, err
				}
				if string(data) != "{}" && string(data) != "[{}]" {
					Qry.ItemValue = string(data)
					if nxtRef[0] == "$WHERE" || nxtRef[0] == "$JOIN" {
						Qry.Rules, _ = getRules(value)
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
				refmap[field] = true
			}
		}
	}

	return query, err
}
func getArrayStruct(dst []interface{}, fieldRef []string) (query []db.Query, err error) {
	return getUsersStuctFields(dst[0], fieldRef)
}

func getFieldRef(dst []byte) (fieldRef []string) {
	// ref := []string{}
	// data := dst
	dec := json.NewDecoder(bytes.NewReader(dst))
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
	// var result map[string]interface{}
	// fmt.Println("data", string(data), "ref-> ", ref)
	// err := json.Unmarshal(data, &result)
	// if err != nil {
	// 	return
	// }
	// return getRef(result, ref)
	return
}

// func getRef(result interface{}, ref []string) (fieldRef []string) {
// 	refmap := map[string]bool{}
// 	// for i, r := range ref {
// 	for i, r := range ref {
// 		if res, ok := result.(map[string]interface{}); ok {
// 			if _, ok := refmap[r]; !ok {
// 				if v, ok := res[r]; ok {
// 					// TODO :- delete r from ref
// 					ref = utils.DeleteArrayElement(ref, r)
// 					fieldRef = append(fieldRef, r)
// 					if _, ok := v.(map[string]interface{}); ok {
// 						fieldRef = append(fieldRef, getRef(v, ref)...)
// 					}
// 					refmap[r] = true
// 				}
// 			}
// 		} else if _, ok := result.([]interface{}); ok {
// 			ref = utils.DeleteArrayElement(ref, r)
// 			fieldRef = append(fieldRef, r)
// 			fmt.Println("arr", i, r, fieldRef)
// 		}
// 	}
// 	return
// }

func parseRules(str string) (rule db.RuleFeild) {
	re := regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)|( [<>=A-Z]{1,4} )|(\w[\w]*)`)
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
	var opre = regexp.MustCompile(`(?m)([<>=A-Z]{1,4})`)
	opr := opre.FindString(str)
	if v, ok := db.Operators[opr]; ok {
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
func getRules(dst interface{}) (rules []db.Rules, err error) {
	rule := db.Rules{}
	if rl, ok := dst.([]interface{}); ok {
		for _, v := range rl {
			if v1, ok := v.(map[string]interface{}); ok {
				// Unmarshal or Decode the JSON to the user struct.
				for field, value := range v1 {
					rule.Label = field
					str, er := getRuleStr(value)
					if er != nil {
						return
					}
					rule.Rule, err = getRule(str)
					if err != nil {
						return
					}
					rules = append(rules, rule)
				}
			}
		}
	}
	return
}
func getRuleStr(value interface{}) (string, error) {
	if str, ok := value.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("rule string not found")
}

func getRule(str string) (rules db.Rule, err error) {
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
func getBrcGrp(str string) (rules db.Rule, err error) {
	bracegrp := regexp.MustCompile(`(?m)(\w[\s.<>=\w]*)`)
	if match := bracegrp.FindString(str); match != "" {
		return getRule(match)
	}
	return
}
func getFirstGroup(str string) (rule db.RuleFeild, err error) {
	frstGrpRule := regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)( [A-Z]{2,4} )([A-z]*[.][A-z]*)`)
	if b := frstGrpRule.FindStringIndex(str); b != nil {
		rule = parseRules(str)
	}

	return
}
