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
			 "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
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

	fmt.Printf("here %v\n", qryAst)

	//check if resulting query ast is correct.

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
