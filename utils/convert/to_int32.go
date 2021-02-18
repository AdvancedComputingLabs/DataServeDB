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
	"math"
	"strconv"
	"strings"
	"unsafe"
)

// public section

func ToInt32(i interface{}, conversionType ConversionClass) (int32, error) {
	//ATTENTION: Approach taken here is to cap at min and max of int32. For weakConversion = false, it returns error.

	//NOTE: See the follow article for integer overflow issue: https://en.wikipedia.org/wiki/Integer_overflow

	//int32 range of values: -2,147,483,648 to 2,147,483,647.
	//WARNING: uint32 value can be greater: 0 to 4,294,967,295.

	const toTypeName = "int32"

	//strict conversions
	switch t := i.(type) {
	case int32:
		return t, nil
	case int:
		if unsafe.Sizeof(t) == 4 {
			return int32(t), nil
		}
	}

	if e := strictError(conversionType, i, toTypeName); e != nil {
		return 0, e //zero is default
	}

	//lossless conversions
	switch t := i.(type) {
	case int: //NOTE: int can be greater than int32, see: https://golang.org/pkg/builtin/#int
		if t >= math.MinInt32 && t <= math.MaxInt32 {
			return int32(t), nil
		}

	case uint: //NOTE: uint can be greater than int32, see: https://golang.org/pkg/builtin/#uint
		if t >= 0 && t <= math.MaxInt32 {
			return int32(t), nil
		}

	case int8:
		return int32(t), nil
	case uint8:
		return int32(t), nil
	case int16:
		return int32(t), nil
	case uint16:
		return int32(t), nil
	}

	// it is not used before because json conversion gives float64, but i may not always be float64 and above type checks in those cases are needed.
	if result, ok := isConversionToInt32Lossless(i); ok {
		return result, nil
	}

	if e := losslessError(conversionType, i, toTypeName); e != nil {
		return 0, e //zero is default
	}

	//weak conversions
	switch t := i.(type) {
	case nil:
		return int32(0), nil //QUESTION: are nil/dbnull lossless?
	case int:
		return capInt32(float64(t)), nil
	case uint:
		return capInt32(float64(t)), nil
	case uint32:
		return capInt32(float64(t)), nil
	case int64:
		return capInt32(float64(t)), nil
	case float32:
		return capInt32(float64(t)), nil
	case float64:
		return capInt32(t), nil
	case string:
		if r, e := strconv.ParseFloat(strings.TrimSpace(t), 64); e == nil {
			return capInt32(r), nil
		}
		goto DEFAULT

	default:
		goto DEFAULT
	}

DEFAULT:
	return 0, fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v'", i, i, toTypeName)
}

// private section

func capInt32(f float64) int32 {
	return int32(capByLimits(f, math.MinInt32, math.MaxInt32))
}

func isConversionToInt32Lossless(v interface{}) (int32, bool) {
	//TODO: post question, not sure this is exactly lossless
	// reason for doing this way is golang json converts every number to float64. But if user sends bigger number it will error but if
	// number is within int32 it will convert. But is this the desired behavor for lossless?
	// Float64 number even if within the int32 range is not lossless imo. I could be wrong on this though.

	//json conversion comes in float64 for any number in golang

	//TODO: check if it has edge cases
	//TODO: check if first conversion to float64 can be used to cover all cases of lossless.

	switch vCasted := v.(type) {

	case float64:
		i32 := int32(vCasted)
		if vCasted == float64(i32) {
			return i32, true
		}

	case string:
		if f64, e := strconv.ParseFloat(vCasted, 64); e == nil {
			i32 := int32(f64)
			if f64 == float64(i32) {
				return i32, true
			}
		}
	}

	return 0, false
}

