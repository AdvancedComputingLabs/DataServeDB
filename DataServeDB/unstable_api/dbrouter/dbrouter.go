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

package dbrouter

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	rules "DataServeDB/dbsystem/rules"
)

//function interface for rest api handler
//NOTE: db is must, hence, db name is extracted by the system before matching.
type HttpRestApiHandlerFnI = func(w http.ResponseWriter, r *http.Request, httpMethod string, dbName string, resPath string /* path remaining after db name */) (resultHttpStatus int, resultContent []byte, resultErr error)

// private as it doesn't need to be exposed.
type reqPathToHandler struct {
	MatchPath string
	matchPathRegEx *regexp.Regexp
	HandlerFn HttpRestApiHandlerFnI
}

// keywords, identifiers, and placeholders
// {db_name} is special name, which is place holder for any db name provided in the path. For example: re_db/tables/users
// db_name must be validated with db naming rule.
const dbNamePlaceHolder = "{db_name}"
const tblNamePlaceHolder = "{tbl_name}"

//Array, first match returns. User have to make sure a mapping doesn't override other mappings
var pathsToHandlers []reqPathToHandler

// TODO: save path to handler mappings upon changes

func init() {
	//TODO: load path to handler mappings
}

func Register(matchPath string, handlerFn HttpRestApiHandlerFnI) error {

	p2h := reqPathToHandler{}

	p2h.MatchPath = matchPath

	sMathPathForRegEx := strings.Replace(matchPath, dbNamePlaceHolder, rules.DbNameValidatorRuleReStrBasic, 1)
	sMathPathForRegEx = strings.Replace(sMathPathForRegEx, tblNamePlaceHolder, rules.TableNameValidatorRuleReStrBasic, 1)
	p2h.matchPathRegEx = regexp.MustCompile(sMathPathForRegEx)

	p2h.HandlerFn = handlerFn

	pathsToHandlers = append(pathsToHandlers, p2h)

	return nil
}

func MatchPathAndCallHandler(w http.ResponseWriter, r *http.Request, reqPath string, httpMethod string) (resultHttpStatus int, resultContent []byte, resultErr error) {
	if pathsToHandlers == nil {
		return 0, nil, errors.New("no match path exits")
	}

	for _, m := range pathsToHandlers {
		if m.matchPathRegEx.MatchString(reqPath) {
			//TODO: extract db name and check
			//TODO: extract correct path for the handler
			//TODO: permissions check for db access?
			//TODO: add auth
			return m.HandlerFn(w, r, httpMethod, "todo", reqPath)
		}
	}

	return 0, nil, errors.New("resource in the path not found")
}
