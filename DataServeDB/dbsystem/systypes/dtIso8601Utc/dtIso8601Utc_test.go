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

package dtIso8601Utc

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
	Cases:
		1) datetime in string
		2) datetime in go's native datetime object
*/

func TestIso8601UtcFromString(t *testing.T) {
	dt, e := Iso8601UtcFromString("2020-01-19T00:00:00Z")
	if e == nil {
		fmt.Println("Iso8601Utc:", dt)
	} else {
		t.Error(e)
	}
}

func TestIso8601UtcNow(t *testing.T) {
	dt := Iso8601UtcNow()
	fmt.Println("Iso8601Utc:", dt)
}

func TestIso8601UseCases(t *testing.T) {

	type DtTester struct {
		Title      string
		MyDateTime Iso8601Utc
	}

	type DtTesterPointerVer struct {
		Title      string
		MyDateTime *Iso8601Utc
	}

	// uninitallized/empty/zero
	{
		dt := Iso8601Utc{}
		fmt.Println("Iso8601Utc:", dt)

		if !dt.IsZero() {
			t.Errorf("coding error! IsZero should be true\n")
		}
		fmt.Println("is zero?", dt.IsZero())
	}

	// from string zero cases
	{
		// this should fail since Iso8601 only handles proper cases. weak conversions should be handled in conversions utility package
		{
			_, e := Iso8601UtcFromString("0")
			if e == nil {
				t.Errorf("coding error! this should fail\n")
			}
		}

		// this should pass
		{
			dt, e := Iso8601UtcFromString("0001-01-01T00:00:00Z")
			if e == nil {
				fmt.Println("Iso8601Utc:", dt)
			} else {
				t.Error(e)
			}
			if !dt.IsZero() {
				t.Errorf("coding error! IsZero should be true\n")
			}
			fmt.Println("is zero?", dt.IsZero())
		}

		// nil case Non Pointer
		{
			tester := DtTester{
				Title: "Nil Case Non Pointer",
			}

			testerJson, e := json.Marshal(tester)
			if e == nil {
				fmt.Println(string(testerJson)) //gives zero value
			}
		}

		// nil case field with pointer
		{
			tester := DtTesterPointerVer{
				Title: "nil case field with pointer",
			}

			testerJson, e := json.Marshal(tester)
			if e == nil {
				fmt.Println(string(testerJson)) //gives null; it never reaches marshaljson method even with 't *Iso8601Utc'
			}
		}

		{ // nil unmarshal case
			sJson := `{"Title":"unmarshal case #0", "MyDateTime":"0900-01-01T00:00:00Z"}`

			var tester DtTesterPointerVer

			e := json.Unmarshal([]byte(sJson), &tester)
			if e == nil {
				fmt.Println("nil unmarshal case:", tester)
			} else {
				fmt.Println(e)
			}
		}

	}

	// extreme datetime ranges
	{
		{ // Year 900 AD
			{
				dt, e := Iso8601UtcFromString("0900-01-01T00:00:00.0000Z")
				if e == nil {
					fmt.Println("Iso8601Utc:", dt)
				} else {
					t.Error(e)
				}
			}

			{
				tester := DtTester{
					Title: "Nil Case Non Pointer",
				}
				tester.MyDateTime, _ = Iso8601UtcFromString("0900-01-01T00:00:00.0000Z")

				testerJson, e := json.Marshal(tester)
				if e == nil {
					fmt.Println(string(testerJson)) //gives zero value
				}
			}

		}
	}

}
