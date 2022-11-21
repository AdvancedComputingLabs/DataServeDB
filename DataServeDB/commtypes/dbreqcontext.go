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

package commtypes

import (
	"DataServeDB/comminterfaces"
	"DataServeDB/dbsystem/constants"
	"DataServeDB/unstable_api/dbrouter"
)

// fields are public, easier.

/*
	Example:
		resPath: /re_db/tables/users/Id:1
		matchedPath: re_db/tables/users
		dbName: re_db
		targetName: users
		targetDbResTypeId: 1
*/

type DbReqContext struct {
	RestMethod   string
	RestMethodId constants.RestMethods
	ResPath      string
	MatchedPath  string
	DbName       string
	Dbi          comminterfaces.DbPtrI
	PathLevels   []dbrouter.PathLevel
}

func NewDbReqContext(restMethod, resPath, matchedPath, dbName string,
	dbi comminterfaces.DbPtrI, pathLevels []dbrouter.PathLevel) *DbReqContext {

	dbreqCtx := DbReqContext{
		RestMethod:  restMethod,
		ResPath:     resPath,
		MatchedPath: matchedPath,
		DbName:      dbName,
		Dbi:         dbi,
		PathLevels:  pathLevels,
	}

	return &dbreqCtx
}
