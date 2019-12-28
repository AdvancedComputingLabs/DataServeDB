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

const (
	dbBool = iota
	dbInt32
	dbString
)

//private to package
type dbTypeBase struct {
	DisplayName string
	DbTypeId    int
}

type DbTypeInterface interface {
	ConvertValue(interface{}, bool) (interface{}, error) //Note: bool is weakConversion
	GetDbTypeDisplayName() string
	GetDbTypeId() int
}

//Rationale: DB server only needs weak or lossless, making it strict makes it too complicated.
//However, utils/convert package needs all three rules as it is general purpose package. -HY
func weakConversionFlagToRule(weakConversion bool) convert.ConversionClass {
	if weakConversion {
		return convert.Weak
	}
	return convert.Lossless
}
