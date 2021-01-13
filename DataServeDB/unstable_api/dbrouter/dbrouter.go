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

/*
	Package Name: dbrouter

	Package Description:
		Provides handling of routing paths. Purpose are as follows:
		i)  gives organized way to manage paths to db resources;
		ii) makes it easier to add more paths for new db resources.

		Important to note: it does not provide if validation of the target resource, it only checks the path exits.
			There is handling function, it should validate the target resource.

	Rules:
		1) Only contains maximum two placeholders at the moment. One should be database name place holder
			and another target resource like a table, file, or other similar database resource target.
		2) database placeholder doesn't have to be there always, if no database is matched empty string is valid in this case.
			Handling function handle empty string(s) according to its needs.
		3) First matching path returns.

	Valid Usage Patterns:
		TODO
*/

package dbrouter

import (
	"DataServeDB/dbsystem/constants"
	"errors"
	"net/http"
	"regexp"
	"strings"

	rules "DataServeDB/dbsystem/rules"
)

//function interface for rest api handler
//NOTE: db is must, hence, db name is extracted by the system before matching.
type HttpRestApiHandlerFnI = func(w http.ResponseWriter, r *http.Request, httpMethod,
	resPath /* path remaining after db name */,
	matchedPath, dbName, targetName string, targetDbResTypeId constants.DbResTypes) (resultHttpStatus int, resultContent []byte, resultErr error)

// private as it doesn't need to be exposed.
type reqPathToHandler struct {
	MatchPath      string
	matchPathRegEx *regexp.Regexp
	HandlerFn      HttpRestApiHandlerFnI
}

//NOTE: placeholders can be here since they don't change.

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

	sMatchPathForRegEx := strings.Replace(matchPath, dbNamePlaceHolder, rules.DbNameValidatorRuleReStrBasic, 1)
	sMatchPathForRegEx = strings.Replace(sMatchPathForRegEx, tblNamePlaceHolder, rules.TableNameValidatorRuleReStrBasic, 1)
	p2h.matchPathRegEx = regexp.MustCompile(sMatchPathForRegEx)

	p2h.HandlerFn = handlerFn

	pathsToHandlers = append(pathsToHandlers, p2h)

	//fmt.Println(sMatchPathForRegEx) //[A-Za-z][_0-9A-Za-z]{2,49}/tables/[A-Za-z][0-9A-Za-z]{2,49}
	//fmt.Printf("%v\n%v\n", pathsToHandlers, p2h.MatchPath) //for debugging
	return nil
}

func MatchPathAndCallHandler(w http.ResponseWriter, r *http.Request, reqPath string, httpMethod string) (resultHttpStatus int, resultContent []byte, resultErr error) {
	if pathsToHandlers == nil {
		return http.StatusNotFound, nil, errors.New("no matching path found")
	}

	var dbName string
	var targetResName string
	var targetDbResTypeid constants.DbResTypes = constants.DbResTypeNone

	//find matching path, first match returns
	for _, m := range pathsToHandlers {
		if path := m.matchPathRegEx.FindString(reqPath); path != "" {

			// placeholders can be in any order, so following loop and matching is used.

			matchPathSplit := strings.Split(m.MatchPath, "/")
			pathSplit := strings.Split(path, "/")

			for i := 0; i < len(matchPathSplit) && i < len(pathSplit); i++ {
				switch matchPathSplit[i] {
				case dbNamePlaceHolder:
					dbName = pathSplit[i]
				case tblNamePlaceHolder:
					targetDbResTypeid = constants.DbResTypeTable
					targetResName = pathSplit[i]
				}
			}

			return m.HandlerFn(w, r, httpMethod, reqPath, path, dbName, targetResName, targetDbResTypeid)
		}
	}

	return http.StatusNotFound, nil, errors.New("no matching path found")
}
