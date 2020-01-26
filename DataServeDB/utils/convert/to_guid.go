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
	"DataServeDB/dbsystem/systypes/guid"
	"fmt"
)

//public section

func ToGuid(v interface{}, conversionType ConversionClass) (guid.Guid, error) {
	const toTypeName = "guid"

	//strict conversions

	switch t := v.(type) {
	case guid.Guid:
		return t, nil
	case *guid.Guid:
		//TODO: if *t is nil?
		return *t, nil
	}

	if e := strictError(conversionType, v, toTypeName); e != nil {
		return guid.Guid{}, e ////zero datetime here
	}

	//lossless conversions

	//TODO: what about nil?

	switch t := v.(type) {
	case string:
		g, e := guid.ParseString(t)
		if e == nil {
			return *g, nil
		}
	}

	if e := losslessError(conversionType, v, toTypeName); e != nil {
		return guid.Guid{}, e //zero datetime here
	}

	//weak conversions

	//NOTE: no weak conversions for guid, as it has specific format

	return guid.Guid{}, fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v'", v, v, toTypeName)
}
