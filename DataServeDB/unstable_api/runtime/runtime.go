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
	"fmt"

	"DataServeDB/dbsystem"
	"DataServeDB/unstable_api/dbrouter"
)

//TODO: change to get functions

var syscasing = dbsystem.SystemCasingHandler.Convert
var isInitalized = false

func IsInitalized() bool {
	return isInitalized
}

func Start(disableHttpServer bool) error {
	//TODO: check if this needs go process to independently initalize db server; there could be hanging issue?
	fmt.Println("Starting DataServeDB server ...")

	//TODO: list/log all the databases being mounted.
	//TODO: refactor

	//TODO: error handling
	loadDatabases()

	//routing
	dbrouter.Register("{DB_NAME}/tables/{TBL_NAME}/{1}.*", TableRestPathHandler)
	dbrouter.Register("{DB_NAME}/tables/{TBL_NAME}", TableRestPathHandler)
	dbrouter.Register("{DB_NAME}/tables", TableRestPathHandler)
	// dbrouter.Register("{DB_NAME}/files/{FIL_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files/{DIR_NAME}/{DIR_NAME}/{DIR_NAME}/{FIL_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files/{DIR_NAME}/{DIR_NAME}/{FIL_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files/{DIR_NAME}/{DIR_NAME}/{DIR_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files/{DIR_NAME}/{FIL_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files/{DIR_NAME}/{DIR_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files/{DIR_NAME}", FileRestPathHandler)
	dbrouter.Register("{DB_NAME}/files", FileRestPathHandler)

	if !disableHttpServer {
		//http server and rest api routing
		StartHttpServer()
	}

	cliProcessor()

	//log.Println("Closing DataServeDB server ...")

	isInitalized = true

	return nil
}
