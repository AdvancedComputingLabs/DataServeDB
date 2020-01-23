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

package dbstrcmp_asciics

import "strings"

// NOTE: For case sensitive using go's own string functions are enough.

type DbStrCmpAsciiCS struct {}

func (t DbStrCmpAsciiCS) AreEqual(s1, s2 string) bool {
	return strings.Compare(s1, s2) == 0
}

func (t DbStrCmpAsciiCS) Compare(s1, s2 string) int {
	return strings.Compare(s1, s2)
}

func (t DbStrCmpAsciiCS) Convert(input string) string {
	// just return as it is, no need for any change.
	return input
}
