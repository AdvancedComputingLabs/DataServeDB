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

package dbtypes

import (
	"DataServeDB/dbtypes/dbtype_props"
)

/*
	Description:
	Iso8601Utc is the standard datetime format in this db server. See: https://www.w3.org/TR/NOTE-datetime
*/

// section: declarations

type dbTypeDateTime struct {
	//private to package
	dbTypeBase
}

type DbTypeDateTimeProperties struct {
	dbtype_props.Conversion
	dbtype_props.Nullable
}

// public

var DateTime = dbTypeDateTime{
	dbTypeBase{
		DbTypeId:    dbDateTime,
		DisplayName: "datetime",
	},
}

func (t dbTypeDateTime) ConvertValue(v interface{}, dbTypeProperties interface{}) (interface{}, error) {
	_ = DateTime //to fix non used error

	//TODO:
	return nil, nil
}

func (t dbTypeDateTime) GetDbTypeDisplayName() string {
	return t.DisplayName
}

func (t dbTypeDateTime) GetDbTypeId() int {
	return t.DbTypeId
}

// public DbTypeDateTimeProperties

func (t dbTypeDateTime) defaultDbTypeProperties() DbTypePropertiesI {
	//TODO:
	return nil
}

func (t *dbTypeDateTime) IsPrimaryKey() bool {
	return false
}

// private
// Section: dbTypeDateTime

func (t dbTypeDateTime) onCreateValidateFieldProperties(fieldProperties interface{}) error {
	//!NotImplemented
	return nil
}

// Section: DbTypeBoolProperties

func defaultDbTypeDateTimeProperties() *DbTypeDateTimeProperties {
	return &DbTypeDateTimeProperties{
		Nullable: dbtype_props.Nullable{State: dbtype_props.NullableFalseDefault},
	}
}

func getDbTypeDateTimePropertiesIndirect(p interface{}) *DbTypeDateTimeProperties {

	switch p := p.(type) {
	case DbTypeDateTimeProperties:
		return &p
	case *DbTypeDateTimeProperties:
		return p
	}

	//TODO: log
	//TODO: panic with code location.
	panic("Coding error, this should not happen!")
}