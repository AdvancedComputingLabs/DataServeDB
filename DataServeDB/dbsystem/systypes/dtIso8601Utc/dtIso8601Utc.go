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
	"errors"
	"fmt"
	"time"
)

// ## Delcarations

const iso8601UtcForm = `2006-01-02T15:04:05.9999999Z`
const iso8601UtcFormInStrQuotes = `"2006-01-02T15:04:05.9999999Z"`

// no need to use date time in the Utc is the give away that is datetime only
type Iso8601Utc time.Time

// ## Static Functions

// Must be in standard Iso8601Utc format '2020-01-19T00:00:00Z'
// Supports precision for time up to 7 decimal places like ''2020-01-19T00:00:00.0000001Z'
func Iso8601UtcFromString(s string) (Iso8601Utc, error) {
	t, err := time.Parse(iso8601UtcForm, s)

	if err != nil {
		return Iso8601Utc(t), err
	}

	return Iso8601Utc(t), nil
}

func Iso8601UtcNow() Iso8601Utc {
	return Iso8601Utc(time.Now().UTC())
}

// ## Methods

//TODO: test performance for Iso8601Utc vs *Iso8601Utc

func (t Iso8601Utc) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t Iso8601Utc) String() string {
	return time.Time(t).Format(iso8601UtcForm)
}

func (t Iso8601Utc) MarshalJSON() ([]byte, error) {
	//Note: Doesn't need to be pointer since it is not changing its state. -hy

	var s string
	//var e error //for later use

	dt_gonative := time.Time(t)
	s = fmt.Sprintf(`"%s"`, dt_gonative.Format(iso8601UtcForm))

	return []byte(s), nil
}

func (t *Iso8601Utc) UnmarshalJSON(data []byte) error {
	//Note: changes state hence used pointer to itself. -hy

	var s string

	s = string(data)

	//for json, date must be quoted as string
	dt_gonative, e := time.Parse(iso8601UtcFormInStrQuotes, s)
	if e != nil {
		//TODO: log error? is it needed?
		fmt.Println(e)
		return errors.New("date/time is not in ISO 8601 UTC format")
	}

	*t = Iso8601Utc(dt_gonative)
	return nil
}
