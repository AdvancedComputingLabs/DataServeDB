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

/*
DbType Properties Concept (and rationale for keeping type properties under this package):
- Keeps property tied to the db type rather than the field. What keywords a table field can support depends on its db type.
- PrimaryKeyable keeps true value if field is primary key or not, it is for convenience.
 */

const (
	dbNullInternalId = iota //if ever used, DbNull is special type
	dbBool
	dbDateTime
	dbInt32
	dbString
)

//DbNull is special type, it does not use iterfaces required by other types. In essence, it is not really a type
// but represents non-existance of a value of other types. DbNull is used to differentiate from null in programming.
// although programming null/nil values from inputs can be (or will be) converted to dbnull in appropriate cases.
type DbNull struct{}

func (t DbNull) String() string {
	return "DbNull"
}

//private to package
type dbTypeBase struct {
	DisplayName string
	DbTypeId    int
}

type DbTypeI interface {
	ConvertValue(value interface{}, dbTypeProperties interface{}) (interface{}, error) //Note: named for clarity
	GetDbTypeDisplayName() string
	GetDbTypeId() int

	//following methods are internal to the package
	defaultDbTypeProperties() DbTypePropertiesI
	onCreateValidateFieldProperties(fieldProperties interface{}) error // convenient to validate some of them here.
}

type DbTypePropertiesI interface {
	IsPrimaryKey() bool
}
