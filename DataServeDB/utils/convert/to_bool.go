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

package convert

import (
	"fmt"
	"strconv"
)

//public section

func ToBool(i interface{}, conversionType ConversionClass) (bool, error) {
	const toTypeName = "bool"

	//strict
	switch t := i.(type) {
	case bool:
		return t, nil
	}

	if e := strictError(conversionType, i, toTypeName); e != nil {
		return false, e
	}

	//lossless
	switch i.(type) {
	case nil:
		return false, nil
	}

	//Rationale: remainder supported type conversions could result is loss of data, e.g. int 5 could be true but true cannot be converted back to 5. -HY

	if e := losslessError(conversionType, i, toTypeName); e != nil {
		return false, e
	}

	//weak
	var r = false

	switch t := i.(type) {
	case int:
		r = t != 0
	case uint:
		r = t != 0
	case int8:
		r = t != 0
	case uint8:
		r = t != 0
	case int16:
		r = t != 0
	case uint16:
		r = t != 0
	case int32:
		r = t != 0
	case uint32:
		r = t != 0
	case int64:
		r = t != 0
	case uint64:
		r = t != 0
	case float32:
		r = t != 0
	case float64:
		r = t != 0
	case string:
		if r, e := strconv.ParseBool(t); e == nil {
			return r, nil
		}
		goto DEFAULT

	default:
		goto DEFAULT
	}

	return r, nil

DEFAULT:
	return false, fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v'", i, i, toTypeName)
}