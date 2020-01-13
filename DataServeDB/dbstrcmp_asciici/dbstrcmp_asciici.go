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

package dbstrcmp_asciici

import "strings"

type DbStrCmpAsciiCI struct {}

func (t DbStrCmpAsciiCI) AreEqual(s1, s2 string) bool {
	return strings.ToUpper(s1) == strings.ToUpper(s2)
}

func (t DbStrCmpAsciiCI) Compare(s1, s2 string) int {
	return strings.Compare(strings.ToUpper(s1), strings.ToUpper(s2))
}

//NOTE: it converts to upper case, I found it to be better than using lower case.
//But there might be edge cases in which it may not be suitable then build your own. --HY
func (t DbStrCmpAsciiCI) Convert(input string) string {
	return strings.ToUpper(input)
}
