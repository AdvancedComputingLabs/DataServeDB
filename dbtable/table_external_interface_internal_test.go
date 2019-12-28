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
	"fmt"
	"reflect"
	"testing"
)

func TestCreateTableToAndFromJsonConversions(t *testing.T) {

	//from obj to json: start
	createTableObj := &createTableExternalInterface{
		TableName: "Tbl01",
		PrimaryKeyName: "Id",
		TableFields: []createTableExternalInterfaceFieldInfo{
			{FieldName: "Id", FieldType: "int32"},
			{FieldName: "UserName", FieldType: "string"},
		},
	}

	createTableObjJSON, _ := json.Marshal(createTableObj)
	//fmt.Println(string(createTableObjJSON))

	if string(createTableObjJSON) !=
		`{"TableName":"Tbl01","PrimaryKeyName":"Id","TableFields":[{"FieldName":"Id","FieldType":"int32"},{"FieldName":"UserName","FieldType":"string"}]}` {
		t.Errorf("Marshalling to json has error(s)\n")
	}

	//from obj to json: end

	//from json to obj: start
	//NOTE: also testing if field names in json have different casing.
	createTableObj2JSON := `{
  "tableName": "Tbl02",
  "PrimaryKeyName":"Id",
  "tableFields": [
    {
      "fieldName": "Id",
      "fieldType": "int32",
      "primaryKey": true
    },
    {
      "FieldName": "UserName",
      "FieldType": "string",
      "PrimaryKey": false
    }
  ]
}`

	var createTableObj2 createTableExternalInterface
	if err := json.Unmarshal([]byte(createTableObj2JSON), &createTableObj2); err != nil {
		t.Error(err)
	}
	fmt.Println(createTableObj2)

	createTableObj2Matcher := createTableExternalInterface{
		TableName: "Tbl02",
		PrimaryKeyName: "Id",
		TableFields: []createTableExternalInterfaceFieldInfo{
			{FieldName: "Id", FieldType: "int32"},
			{FieldName: "UserName", FieldType: "string"},
		},
	}

	if reflect.DeepEqual(createTableObj2, createTableObj2Matcher) == false  {
		t.Errorf("Unmarshalling from json has error(s)\n")
	}

	//from json to obj: end
}
