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

package guid

import (
	"encoding/gob"
	"github.com/beevik/guid"
)

type Guid guid.Guid

func IsValidGuid(s string) bool {
	return guid.IsGuid(s)
}

func NewGuid() *Guid {
	g := Guid(*guid.New())
	return &g
}

func NewGuidString () string {
	return guid.NewString()
}

func ParseString(s string) (*Guid, error) {
	g, e := guid.ParseString(s)
	if e != nil {
		return nil, e
	}
	g2 := Guid(*g)
	return &g2, nil
}

func (t Guid) MarshalBinary() ([]byte, error) {
	g := guid.Guid(t)
	return []byte(g.String()), nil
}

func (t Guid) MarshalJSON() ([]byte, error) {
	//Note: Doesn't need to be pointer since it is not changing its state. -hy

	var s string
	g := guid.Guid(t)
	s = `"` + g.String() + `"`
	return []byte(s), nil
}

func (t *Guid) UnmarshalBinary(data []byte) error {
	s := string(data)
	g, e := ParseString(s)
	if e != nil {
		return e
	}
	*t = *g
	return nil
}

func (t *Guid) UnmarshalJSON(data []byte) error {
	//Note: changes state hence used pointer to itself. -hy

	s := string(data)

	if len(s) > 2 {
		s = s[1:len(s)-1]
	}

	g, e := ParseString(s)
	if e != nil {
		return e
	}

	*t = *g

	return nil
}

//private
func init() {
	gob.Register(Guid{})
}



