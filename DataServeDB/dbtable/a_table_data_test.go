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
	"testing"

	"DataServeDB/dbsystem"
	"DataServeDB/dbtypes"
)

func TestAddData_Normal(t *testing.T) {

	dbFieldRule := tableFieldStruct{
		FieldInternalId: 0,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: 1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	/*
		Id: 1, IsTrue: true
	*/

	row001 := tableRowByInternalIds{
		0: 1,
		1: true,
	}

	tdata := tableDataContainer{}
	tdata.Rows = append(tdata.Rows, row001)

	fmt.Printf("%v\n", tdata.Rows)
}
