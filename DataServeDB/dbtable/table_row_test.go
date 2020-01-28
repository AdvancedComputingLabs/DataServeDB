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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"DataServeDB/dbsystem"
	"DataServeDB/dbtypes"
)

//Field meta data should have public access?
//

/*

External Operations:
1) Create table with fields and their properties.
2) Insert new row.
3) Get table row.
4) Update table row.

*/

func TestAddData_Temp(t *testing.T) {

	//dbtypes.

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

	{ //direct TableRow creation
		row01 := TableRow{
			"Id":     1,
			"IsTrue": true,
		}

		row01_internal, e := fromLabeledByFieldNames(row01, tb01, dbsystem.SystemCasingHandler)

		if e == nil {
			fmt.Printf("%v\n", row01_internal)
		}
	}

	{ // from json creation, test: non-existent field
		row01Json := `{
			"Id" : 1,
			"WrongName" : true
		}`

		var row01 TableRow

		json.Unmarshal([]byte(row01Json), &row01)

		_, e := fromLabeledByFieldNames(row01, tb01, dbsystem.SystemCasingHandler)

		if e != nil {
			fmt.Printf("%v\n", e.Error())
		} else {
			t.Error(errors.New("Should fail."))
		}
	}

	{ // from json creation, test: with one field left out.

		//NOTE: this works at the moment unless there are non-nullable without default value.

		row01Json := `{
			"Id" : 1
		}`

		var row01 TableRow

		json.Unmarshal([]byte(row01Json), &row01)

		row01_internal, e := fromLabeledByFieldNames(row01, tb01, dbsystem.SystemCasingHandler)

		if e == nil {
			fmt.Printf("%v\n", row01_internal)
		} else {
			t.Error(errors.New("Should fail."))
		}
	}

	{ //test: from json creation, with different casing.

		row01Json := `{
			"Id" : 1,
			"Istrue" : true
		}`

		var row01 TableRow

		json.Unmarshal([]byte(row01Json), &row01)

		row01_internal, e := fromLabeledByFieldNames(row01, tb01, dbsystem.SystemCasingHandler)

		if e == nil {
			fmt.Printf("%v\n", row01_internal)
		} else {
			t.Error(e)
		}

	}

	//TODO: test if json to native type conversions were correct e.g. number value has converted to int64 and not int32

	{ //test: from json creation, one field has wrong type

		row01Json := `{
			"Id" : 1,
			"IsTrue" : true
		}`

		var row01 TableRow

		json.Unmarshal([]byte(row01Json), &row01)

		row01_internal, e := fromLabeledByFieldNames(row01, tb01, dbsystem.SystemCasingHandler)

		if e == nil {
			fmt.Printf("%v\n", row01_internal)
			typ := reflect.TypeOf(row01_internal[1])
			fmt.Println("Converted type:", typ)
			if typ.Name() != "bool" {
				t.Error("type should be bool")
			}
		} else {
			t.Error(e)
		}

	}

	//TODO: test different case.

	var row01_internal_g tableRowByInternalIds

	{ //from json creation

		row01Json := `{
			"Id" : 1,
			"IsTrue" : true
		}`

		var row01 TableRow

		json.Unmarshal([]byte(row01Json), &row01)

		row01_internal, e := fromLabeledByFieldNames(row01, tb01, dbsystem.SystemCasingHandler)

		if e == nil {
			fmt.Printf("%v\n", row01_internal)
		} else {
			t.Error(e)
		}

		row01_internal_g = row01_internal
	}

	{ //from tableRowByInternalIds to TableRow
		row01, e := toLabeledByFieldNames(row01_internal_g, tb01)

		if e == nil {
			fmt.Printf("%v\n", row01)
		} else {
			t.Error(e)
		}

		row01JsonBytes, e := json.Marshal(row01)

		if e != nil {
			t.Error(e)
		}

		row01Json := string(row01JsonBytes)

		fmt.Println(row01Json)
	}

}
