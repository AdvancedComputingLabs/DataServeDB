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

// public section

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