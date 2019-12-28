// Copyright (c) 2019 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package dbtable

import (
	"fmt"

	"DataServeDB/dbsystem"
	"DataServeDB/dbtypes"
)

var syscasing = dbsystem.SystemCasingHandler.Convert

var dbtypes_map = map[string]dbtypes.DbTypeInterface {}

func addDbTypeToMap(dbtype dbtypes.DbTypeInterface) {
	cased_type_name := syscasing(dbtype.GetDbTypeDisplayName())
	dbtypes_map[cased_type_name] = dbtype
}

func getDbType(dbtype_name string) (dbtypes.DbTypeInterface, error) {
	cased_type_name := syscasing(dbtype_name)
	if dt, ok := dbtypes_map[cased_type_name]; ok {
		return dt, nil
	}
	return nil, fmt.Errorf("variable type '%s' doesn't exist", dbtype_name)
}

func init() {
	addDbTypeToMap(dbtypes.Bool)
	addDbTypeToMap(dbtypes.Int32)
	addDbTypeToMap(dbtypes.String)
}