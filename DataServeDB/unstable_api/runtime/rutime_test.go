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
	"fmt"
	"testing"
)

func TestGens(t *testing.T) {
	loadDatabases()
}

//putting it here in case, parsing of the path belongs in db package.
//update: runtime was giving cyclic imports, hence, moved here.

/*
	Example:
		resPath: /re_db/tables/users/Id:1
		matchedPath: re_db/tables/users
		dbName: re_db
		targetName: users
		targetDbResTypeId: 1
*/

func TestTableGet(t *testing.T) {

	errMountingDb := mountDb("re_db", "./../../Databases")
	if errMountingDb != nil {
		t.Fatal(errMountingDb)
	}

	db, e := GetDb("re_db")
	if e != nil {
		//at the moment it only returns database not found.
		//fmt.Print(e)
		t.Fatal(e)
	}

	{
		var dbReqCtx *commtypes.DbReqContext
		dbReqCtx = commtypes.NewDbReqContext(
			"GET", "/re_db/tables/Tbl01/1", "re_db/tables/users",
			"re_db", db, "Tbl01", constants.DbResTypeTable)

		resultHttpStatus, resultContent, resultErr := db.TablesGet(dbReqCtx)

		if resultErr != nil {
			t.Fatal(resultErr)
		}
		fmt.Println(resultHttpStatus)
		fmt.Println(string(resultContent))
	}

	{
		var dbReqCtx *commtypes.DbReqContext
		dbReqCtx = commtypes.NewDbReqContext(
			"GET", "/re_db/tables/Tbl01/Id:1", "re_db/tables/users",
			"re_db", db, "Tbl01", constants.DbResTypeTable)

		resultHttpStatus, resultContent, resultErr := db.TablesGet(dbReqCtx)

		if resultErr != nil {
			t.Fatal(resultErr)
		}
		fmt.Println(resultHttpStatus)
		fmt.Println(string(resultContent))
	}

}