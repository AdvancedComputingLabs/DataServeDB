// Copyright 2009 The Go Authors. All rights reserved.

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

package dbstrcmp_unicodesimplefoldci

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

/*
	Description: Simple fold serves as a decent case insensitive and culture invariant comparision and folding
		for Unicode characters. But it cannot cover all case insensitive cases, it use in the db is by design for fast Unicode case insensitive comparer.
		For issues with Unicode casing and case insensitive comparision see: https://www.w3.org/International/wiki/Case_folding
		-- HY (02-Dec-2019)
 */

type DbStrCmpUnicodeSimplefoldCI struct {}

func (t DbStrCmpUnicodeSimplefoldCI) AreEqual(s1, s2 string) bool {
	return EqualFold(s1, s2)
}

func (t DbStrCmpUnicodeSimplefoldCI) Compare(s1, s2 string) int {
	panic("unimplemented method!")
	return -2
}

func (t DbStrCmpUnicodeSimplefoldCI) Convert(input string) string {
 	return ToEqualFold(input)
}

func EqualFold(s, t string) bool {
	for s != "" && t != "" {
		// Extract first rune from each string.
		var sr, tr rune
		if s[0] < utf8.RuneSelf {
			sr, s = rune(s[0]), s[1:]
		} else {
			r, size := utf8.DecodeRuneInString(s)
			sr, s = r, s[size:]
		}
		if t[0] < utf8.RuneSelf {
			tr, t = rune(t[0]), t[1:]
		} else {
			r, size := utf8.DecodeRuneInString(t)
			tr, t = r, t[size:]
		}

		// If they match, keep going; if not, return false.

		// Easy case.
		if tr == sr {
			continue
		}

		// Make sr < tr to simplify what follows.
		if tr < sr {
			tr, sr = sr, tr
		}
		// Fast check for ASCII.
		if tr < utf8.RuneSelf && 'A' <= sr && sr <= 'Z' {
			// ASCII, and sr is upper case.  tr must be lower case.
			if tr == sr+'a'-'A' {
				continue
			}
			return false
		}

		// General case. SimpleFold(x) returns the next equivalent rune > x
		// or wraps around to smaller values.
		r := unicode.SimpleFold(sr)
		for r != sr && r < tr {
			r = unicode.SimpleFold(r)
		}
		if r == tr {
			continue
		}
		return false
	}

	// One string is empty. Are both?
	return s == t
}

func ToEqualFold(s string) string {
	var rs bytes.Buffer
	//Note: range here decodes codepoint correctly.
	for _, r := range s {
		switch {

		case r <= 'z':
			rs.WriteRune(unicode.ToUpper(r))

		default:
			rs1 := unicode.SimpleFold(r)
			rs2 := unicode.SimpleFold(rs1)

			if rs1 != rs2 && (rs2 != r || rs1 != r) {
				r = unicode.ToUpper(r)
			} else {
				r = rs1
			}

			rs.WriteRune(r)
		}

	}

	return rs.String()
}
