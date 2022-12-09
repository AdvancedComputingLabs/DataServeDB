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
		TODO:
			- db name and target place holders are hard coded, it would be better if they are provided during registeration or configuration?
			- resource type ids are also hard coded.
			- match path should be case sensitive?
*/

package dbrouter

import (
	"DataServeDB/utils/rest"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"DataServeDB/dbsystem/constants"
	rules "DataServeDB/dbsystem/rules"
)

//NOTE: placeholders can be here since they don't change.

// keywords, identifiers, and placeholders
// {DB_NAME} is special name, which is place holder for any db name provided in the path. For example: re_db/tables/users
// db_name must be validated with db naming rule.
const (
	dbNamePlaceHolder   = "{DB_NAME}"
	tblNamePlaceHolder  = "{TBL_NAME}"
	fileNamePlaceHolder = "{FIL_NAME}"
	dirNamePlaceHolder  = "{DIR_NAME}"
)

type PathLevel struct {
	PathItem       string
	PathItemTypeId constants.DbResTypes
}

// HttpRestApiHandlerFnI function interface for rest api handler
// NOTE: db is must, hence, db name is extracted by the system before matching.
type HttpRestApiHandlerFnI = func(w http.ResponseWriter, r *http.Request, httpMethod,
	resPath, /* path remaining after db name */
	matchedPath, dbName string, pathLevels []PathLevel)

// private as it doesn't need to be exposed.
type reqPathToHandler struct {
	MatchPath      string
	matchPathRegEx *regexp.Regexp
	HandlerFn      HttpRestApiHandlerFnI
}

// Array, first match returns. User have to make sure a mapping doesn't override other mappings
var pathsToHandlers []reqPathToHandler

// TODO: save path to handler mappings upon changes

func init() {
	//TODO: load path to handler mappings
}

func NewPathLevel(pathItem string, pathItemTypeId constants.DbResTypes) PathLevel {
	return PathLevel{PathItem: pathItem, PathItemTypeId: pathItemTypeId}
}

func Register(matchPath string, handlerFn HttpRestApiHandlerFnI) error {

	p2h := reqPathToHandler{}

	p2h.MatchPath = strings.ToUpper(matchPath)

	sMatchPathForRegEx := strings.Replace(matchPath, dbNamePlaceHolder, rules.DbNameValidatorRuleReStrBasic, 1)
	sMatchPathForRegEx = strings.Replace(sMatchPathForRegEx, tblNamePlaceHolder, rules.TableNameValidatorRuleReStrBasic, 1)
	sMatchPathForRegEx = strings.Replace(sMatchPathForRegEx, fileNamePlaceHolder, rules.FileNameValidator, 1)
	sMatchPathForRegEx = strings.Replace(sMatchPathForRegEx, dirNamePlaceHolder, rules.DirNameValidator, 3)
	p2h.matchPathRegEx = regexp.MustCompile("(?i)" + sMatchPathForRegEx) // (?i) makes it case-insensitive

	p2h.HandlerFn = handlerFn

	pathsToHandlers = append(pathsToHandlers, p2h)

	//fmt.Println(sMatchPathForRegEx) //[A-Za-z][_0-9A-Za-z]{2,49}/tables/[A-Za-z][0-9A-Za-z]{2,49}
	//fmt.Printf("%v\n%v\n", pathsToHandlers, p2h.MatchPath) //for debugging
	return nil
}

func MatchPathAndCallHandler(w http.ResponseWriter, r *http.Request, reqPath string, httpMethod string) {

	if pathsToHandlers == nil {
		rest.ResponseWriteHelper(w, http.StatusTeapot, nil, errors.New("CodingError: pathsToHandlers is nil"))
		return
	}

	var dbName string

	//find matching path, first match returns
	for _, m := range pathsToHandlers {
		if path := m.matchPathRegEx.FindString(reqPath); path != "" {

			// placeholders can be in any order, so following loop and matching is used.

			matchPathSplit := strings.Split(m.MatchPath, "/")
			pathSplit := strings.Split(path, "/")

			// path levels
			var pathLevels []PathLevel

			for i := 0; i < len(matchPathSplit) && i < len(pathSplit); i++ {
				switch matchPathSplit[i] {
				case dbNamePlaceHolder:
					dbName = pathSplit[i]
				case "TABLES":
					pathLevels = append(pathLevels, NewPathLevel(pathSplit[i], constants.DbResTypeTablesNamespace))
				case "FILES":
					pathLevels = append(pathLevels, NewPathLevel(pathSplit[i], constants.DbResTypeFileNamespace))
				case tblNamePlaceHolder:
					pathLevels = append(pathLevels, NewPathLevel(pathSplit[i], constants.DbResTypeTable))
				case fileNamePlaceHolder:
					pathLevels = append(pathLevels, NewPathLevel(pathSplit[i], constants.DbResTypeFile))
				case dirNamePlaceHolder:
					pathLevels = append(pathLevels, NewPathLevel(pathSplit[i], constants.DbResTypeDirName))
				default: // suppose to lower level like row.
					pathLevels = append(pathLevels, NewPathLevel(pathSplit[i], constants.DbResTypeUndefined))
				}
			}

			m.HandlerFn(w, r, httpMethod, reqPath, path, dbName, pathLevels)
			return
		}
	}

	rest.ResponseWriteHelper(w, http.StatusNotFound, nil, errors.New("no matching rest path found"))
}
