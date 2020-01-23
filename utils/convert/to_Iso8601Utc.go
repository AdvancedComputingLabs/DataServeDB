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

package convert

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"DataServeDB/dbsystem/systypes/dtIso8601Utc"
)

//public section

func ToIso8601Utc(v interface{}, conversionType ConversionClass) (dtIso8601Utc.Iso8601Utc, error) {
	const toTypeName = "Iso8601UtcDateTime"

	//strict conversions

	switch t := v.(type) {
	case dtIso8601Utc.Iso8601Utc:
		return t, nil
	case *dtIso8601Utc.Iso8601Utc:
		return *t, nil
	}

	if e := strictError(conversionType, v, toTypeName); e != nil {
		return dtIso8601Utc.Iso8601Utc{}, e ////zero datetime here
	}

	//lossless conversions

	//NOTE: couldn't check with time.Time if it was set in UTC. Lossing UTC is not lossless, hence, time.Time is in weak conversion.

	switch t := v.(type) {
	case string:
		dt, e := dtIso8601Utc.Iso8601UtcFromString(t)
		if e == nil {
			return dt, nil
		}
	}

	if e := losslessError(conversionType, v, toTypeName); e != nil {
		return dtIso8601Utc.Iso8601Utc{}, e //zero datetime here
	}

	//weak conversions

	switch t := v.(type) {
	case time.Time:
		return dtIso8601Utc.Iso8601Utc(t.UTC()), nil

	case string:
		return iso8601UtcWeakStringConversions(t)
	}

	return dtIso8601Utc.Iso8601Utc{}, fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v'", v, v, toTypeName)
}

// !WARNING: unstable api
func iso8601UtcWeakStringConversions(s string) (dtIso8601Utc.Iso8601Utc, error) {
	//TODO: Only year should be valid date? It will be missing month and day.

	re := regexp.MustCompile(`(?:(\d{4})(?:[-](\d{2})(?:[-](\d{2}))*)*(?:[Tt](\d{2})(?:[:](\d{2})(?:[:](\d{2}[.0-9]{0,9}))*)*)*)`)

	matches := re.FindAllStringSubmatch(s, -1)

	if len(matches) < 1 {
		return dtIso8601Utc.Iso8601Utc{}, errors.New("error in parsing Iso8601Utc datetime")
	}

	if len(matches[0]) < 2 && matches[0][1] != "" {
		return dtIso8601Utc.Iso8601Utc{}, errors.New("error in parsing Iso8601Utc datetime")
	}

	var yr int
	var mon time.Month = 0
	var day int = 0

	var hr int = 0
	var min int = 0
	var sec int = 0

	//year
	if tmp, err := ToInt32(matches[0][1], Lossless); err != nil {
		return dtIso8601Utc.Iso8601Utc{}, err
	} else {
		yr = int(tmp)
	}

	//month
	if len(matches[0]) > 2 && matches[0][2] != "" {
		if tmp, err := ToInt32(matches[0][2], Lossless); err != nil {
			return dtIso8601Utc.Iso8601Utc{}, err
		} else {
			mon = time.Month(tmp)
		}
	}

	//day
	if len(matches[0]) > 3 && matches[0][3] != "" {
		if tmp, err := ToInt32(matches[0][3], Lossless); err != nil {
			return dtIso8601Utc.Iso8601Utc{}, err
		} else {
			day = int(tmp)
		}
	}

	//hour
	if len(matches[0]) > 4 && matches[0][4] != "" {
		if tmp, err := ToInt32(matches[0][4], Lossless); err != nil {
			return dtIso8601Utc.Iso8601Utc{}, err
		} else {
			hr = int(tmp)
		}
	}

	//min
	if len(matches[0]) > 5 && matches[0][5] != "" {
		if tmp, err := ToInt32(matches[0][5], Lossless); err != nil {
			return dtIso8601Utc.Iso8601Utc{}, err
		} else {
			min = int(tmp)
		}
	}

	//second
	if len(matches[0]) > 6 && matches[0][6] != "" {
		if tmp, err := ToInt32(matches[0][6], Lossless); err != nil {
			return dtIso8601Utc.Iso8601Utc{}, err
		} else {
			sec = int(tmp)
		}
	}

	dt := time.Date(yr, mon, day, hr, min, sec, 0, time.UTC)

	return dtIso8601Utc.Iso8601Utc(dt), nil
}
