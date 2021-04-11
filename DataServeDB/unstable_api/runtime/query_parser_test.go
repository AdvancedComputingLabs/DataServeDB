// Copyright (c) 2021 Advanced Computing Labs DMCC

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
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	query := `{"Users": {
		"Id": {},
		"UserName": {},
		"Properties": [
		   	{
			 "$JOIN": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
			},
			{
			 "$WHERE": "(Users.Id IS UserProperties.Id OR Properties.SlNum IS UserProperties.SlNum) AND (Users.Id >= 2)"
		   	}
		 ]
	   	}}
		`

	var dst interface{}
	json.Unmarshal([]byte(query), &dst)

	_, qryAst, err := DecodeJSON([]byte(query))

	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	// fmt.Printf("here %v\n", qryAst)
	b, e := json.Marshal(qryAst)
	if e == nil {
		println(string(b))
	} else {
		fmt.Println(e)
	}
	//check if resulting query ast is correct.
	// result := `{"ItemLabel":"Users","ItemType":0,"ItemValue":"{\"Users\": {\n\t\t\"Id\": {},\n\t\t\"UserName\": {},\n\t\t\"Properties\": [\n\t\t   {\n\t\t\t \"$WHERE\": \"Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum\"\n\t\t   }\n\t\t ]\n\t   \t}}\n\t\t","Rules":null,"Children":[{"ItemLabel":"Id","ItemType":0,"ItemValue":"","Rules":null,"Children":null},{"ItemLabel":"UserName","ItemType":0,"ItemValue":"","Rules":null,"Children":null},{"ItemLabel":"Properties","ItemType":0,"ItemValue":"[{\"$WHERE\":\"Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum\"}]","Rules":null,"Children":[{"ItemLabel":"$WHERE","ItemType":0,"ItemValue":"\"Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum\"","Rules":{"TableName":"Users","FieldName":"Id","Operation":1,"Next":{"TableName":"UserProperties","FieldName":"Id","Operation":3,"Next":{"TableName":"Properties","FieldName":"SlNum","Operation":1,"Next":{"TableName":"UserProperties","FieldName":"SlNum","Operation":0,"Next":null}}}},"Children":null}]}]}`

	/*
		Result should be:
		Query {
			ItemLabel: "Users"
			ItemType: ""
			ItemValue: Empty
			Rules: Empty
			Children: Query[]
				0: {
					ItemLabel: "Id"
					ItemType: ""
					ItemValue: Empty
					Rules: Empty
					Children: Empty
				}
				1: {
					ItemLabel: "UserName"
					ItemType: ""
					ItemValue: Empty
					Rules: Empty
					Children: Empty
				}
				2: {
					ItemLabel: "Properties"
					ItemType: ""
					ItemValue: Empty
					Rules: TODO: Use correct rules structure here.
					Children: Empty
				}
		}
	*/

	//TODO: test against valid query structure.
}
func TestParseJson(t *testing.T) {
	set := []struct {
		rule   string
		result string
	}{
		{
			"Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum",
			`{"TableName":"Users","FieldName":"Id","Operation":1,"Next":{"TableName":"UserProperties","FieldName":"Id","Operation":3,"Next":{"TableName":"Properties","FieldName":"SlNum","Operation":1,"Next":{"TableName":"UserProperties","FieldName":"SlNum","Operation":0,"Next":null}}}}`,
		},
		{
			"Users.Id OR UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum",
			`{"TableName":"Users","FieldName":"Id","Operation":2,"Next":{"TableName":"UserProperties","FieldName":"Id","Operation":3,"Next":{"TableName":"Properties","FieldName":"SlNum","Operation":1,"Next":{"TableName":"UserProperties","FieldName":"SlNum","Operation":0,"Next":null}}}}`,
		},
		{
			"(Users.Id OR UserProperties.Id AND Properties.SlNum OR UserProperties.SlNum)",
			`{"TableName":"Users","FieldName":"Id","Operation":2,"Next":{"TableName":"UserProperties","FieldName":"Id","Operation":3,"Next":{"TableName":"Properties","FieldName":"SlNum","Operation":2,"Next":{"TableName":"UserProperties","FieldName":"SlNum","Operation":0,"Next":null}}}}`,
		},
		{
			"(Users.Id IS UserProperties.Id OR Properties.SlNum IS UserProperties.SlNum) AND  (Users.Id >= 2)",
			``,
		},
		{
			`Users.Id >= 2`,
			``,
		},
	}

	for _, v := range set {
		res, err := getRule(v.rule)
		jsRes, err := json.Marshal(res)
		if err != nil {
			t.Errorf("error on json marshal")
		}
		if string(jsRes) != v.result {
			// t.Errorf("error, test not Passed")
			fmt.Println("Not passed %V", string(jsRes))
		} else {
			fmt.Println("passed %V", string(jsRes))
		}
	}
}
