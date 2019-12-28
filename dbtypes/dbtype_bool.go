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

package dbtypes

import (
	"DataServeDB/utils/convert"
)

type dbTypeBool struct {
	//private to package
	dbTypeBase
}

func (t dbTypeBool) ConvertValue(v interface{}, weakConversion bool) (interface{}, error) {
	return convert.ToBool(v, weakConversionFlagToRule(weakConversion))
}

func (t dbTypeBool) GetDbTypeDisplayName() string {
	return t.DisplayName
}

func (t dbTypeBool) GetDbTypeId() int {
	return t.DbTypeId
}

var Bool = dbTypeBool{
	dbTypeBase{
		DbTypeId:    dbBool,
		DisplayName: "bool",
	},
}



