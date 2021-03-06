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
	"DataServeDB/dbtypes/dbtype_props"
	"DataServeDB/utils/convert"
)

// section: declarations

type dbTypeBool struct {
	//private to package
	dbTypeBase
}

type DbTypeBoolProperties struct {
	dbtype_props.Conversion
	dbtype_props.Nullable
}

// public

var Bool = dbTypeBool{
	dbTypeBase{
		DbTypeId:    dbBool,
		DisplayName: "bool",
	},
}

func (t dbTypeBool) ConvertValue(v interface{}, dbTypeProperties interface{}) (interface{}, error) {
	p := getDbTypeBoolPropertiesIndirect(dbTypeProperties)

	if p == nil {
		//TODO: log.
		//TODO: update with location of the code.
		panic("coding error!")
	}

	return convert.ToBool(v, p.ToSystemConversionClass())
}

func (t dbTypeBool) GetDbTypeDisplayName() string {
	return t.DisplayName
}

func (t dbTypeBool) GetDbTypeId() int {
	return t.DbTypeId
}

// public DbTypeBoolProperties

func (t dbTypeBool) defaultDbTypeProperties() DbTypePropertiesI {
	return defaultDbTypeBoolProperties()
}

func (t *DbTypeBoolProperties) IsPrimaryKey() bool {
	return false
}

// private
// Section: DbTypeBool

func (t dbTypeBool) onCreateValidateFieldProperties(fieldProperties interface{}) error {
	//!NotImplemented
	return nil
}

// Section: DbTypeBoolProperties

func defaultDbTypeBoolProperties() *DbTypeBoolProperties {
	return &DbTypeBoolProperties{
		Nullable: dbtype_props.Nullable{State: dbtype_props.NullableFalseDefault},
	}
}

func getDbTypeBoolPropertiesIndirect(p interface{}) *DbTypeBoolProperties {

	switch p := p.(type) {
	case DbTypeBoolProperties:
		return &p
	case *DbTypeBoolProperties:
		return p
	}

	//TODO: log
	//TODO: panic with code location.
	panic("Coding error, this should not happen!")
}

// Section: <dbTypeBoolValueOrFun>

// </dbTypeBoolValueOrFun>

