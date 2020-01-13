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

//private type decl

//public type decl
type ConversionClass int

const (
	Strict ConversionClass = iota 	//Only if type is same.

	Lossless //Only if conversion doesn't loose data. e.g. int to and from string will not lose data, hence, it
				// will not error under Lossless conversion type.

	Weak	//It will convert in most cases but to and from type may not result is same value.
				// For example float might lose data converting to and from string.
)

//private section

func capByLimits(f float64, min float64, max float64) float64 {
	if f >= min && f <= max {
		return f
	} else if f >= min {
		return max
	} else {
		return min
	}
}

func capInt32(f float64) int32 {
	return int32(capByLimits(f, math.MinInt32, math.MaxInt32))
}

func strictError(cc ConversionClass, i interface{}, toTypeName string) error {
	if cc == Strict {
		return fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v' due to strict type conversion rule", i, i, toTypeName)
	}
	return nil
}

func losslessError(cc ConversionClass, i interface{}, toTypeName string) error {
	if cc == Lossless {
		return fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v' due to lossless type conversion rule", i, i, toTypeName)
	}
	return nil
}

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

	if e := losslessError(conversionType, i, toTypeName); e != nil {
		return 0, e //zero is default
	}

	//weak conversions
	switch t := i.(type) {
	case nil:
		return int32(0), nil
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

func ToString(i interface{}, conversionType ConversionClass) (string, error) {
	const toTypeName = "string"

	//strict
	switch t := i.(type) {
	case string:
		return t, nil
	}

	if e := strictError(conversionType, i, toTypeName); e != nil {
		return "", e //empty string is default
	}

	//lossless
	switch t := i.(type) {
	case int:
		return strconv.Itoa(t), nil
	case int8:
		return strconv.Itoa(int(t)), nil
	case int16:
		return strconv.Itoa(int(t)), nil
	case int32:
		return strconv.Itoa(int(t)), nil
	case int64:
		return strconv.Itoa(int(t)), nil
	case uint:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(t), 10), nil
	case uint64:
		return strconv.FormatUint(t, 10), nil
	}

	if e := losslessError(conversionType, i, toTypeName); e != nil {
		return "", e //empty string is default
	}

	//weak
	switch t := i.(type) {
	case nil:
		return "", nil //this case is confusing. To keep it under strong conversion or weak conversion. Resolution: I'm keeping it under weak - HY 07-Dec-2019.
	case float32:
		return strconv.FormatFloat(float64(t),'E', -1, 32), nil
	case float64:
		return strconv.FormatFloat(t,'E', -1, 64), nil
	default:
		goto DEFAULT
	}

DEFAULT:
	return "", fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v'", i, i, toTypeName)
}
