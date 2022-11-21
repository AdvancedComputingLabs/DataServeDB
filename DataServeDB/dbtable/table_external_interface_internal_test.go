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
	"testing"
)

func TestCreateTableToAndFromJsonConversions(t *testing.T) {

	//from obj to json: start
	//TODO: code obsolete, need to remake these tests

	//from obj to json: end

	//from json to obj: start
	//NOTE: also testing if field names in json have different casing.
	//TODO: code obsolete, need to remake these tests
	//	createTableObj2JSON := `{
	//  "tableName": "Tbl02",
	//  "PrimaryKeyName":"Id",
	//  "tableFields": [
	//    {
	//      "fieldName": "Id",
	//      "fieldType": "int32",
	//      "primaryKey": true
	//    },
	//    {
	//      "FieldName": "UserName",
	//      "FieldType": "string",
	//      "primaryKey": false
	//    }
	//  ]
	//}`
	//
	//	var createTableObj2 createTableExternalStruct
	//	if err := json.Unmarshal([]byte(createTableObj2JSON), &createTableObj2); err != nil {
	//		t.Error(err)
	//	}
	//	fmt.Println(createTableObj2)
	//
	//	createTableObj2Matcher := createTableExternalStruct{
	//		TableName: "Tbl02",
	//		PrimaryKeyName: "Id",
	//		TableColumns: []createTableExternalInterfaceFieldInfo{
	//			{FieldName: "Id", FieldType: "int32"},
	//			{FieldName: "UserName", FieldType: "string"},
	//		},
	//	}

	//if reflect.DeepEqual(createTableObj2, createTableObj2Matcher) == false  {
	//	t.Errorf("Unmarshalling from json has error(s)\n")
	//}

	//from json to obj: end
}
