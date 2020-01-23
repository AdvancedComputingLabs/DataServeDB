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
	"errors"

	"DataServeDB/dbsystem/systypes/dtIso8601Utc"
	"DataServeDB/dbtypes/dbtype_props"
	"DataServeDB/parsers"
	"DataServeDB/utils/convert"
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

type dbTypeDateTimeFp = func() dtIso8601Utc.Iso8601Utc

type dbTypeDateTimeValueOrFun struct {
	// Go doesn't have discriminated union so this will do.
	number *dtIso8601Utc.Iso8601Utc
	fun    dbTypeInt32Fp
}

type DbTypeDateTimeProperties struct {
	dbtype_props.Conversion
	dbtype_props.Nullable
	Auto    dbTypeDateTimeValueOrFun
	Default dbTypeDateTimeValueOrFun
}

// public

var DateTime = dbTypeDateTime{
	dbTypeBase{
		DbTypeId:    dbDateTime,
		DisplayName: "datetime",
	},
}

func (t dbTypeDateTime) ConvertValue(v interface{}, dbTypeProperties interface{}) (interface{}, error) {

	//NOTE#1: auto and default are validated during property setting, but function returned values are not
	//	guaranteed so these result values are still validated here.

	var dt dtIso8601Utc.Iso8601Utc
	var e error

	if v == nil {
		v = DbNull{}
	}

	p := getDbTypeDateTimePropertiesIndirect(dbTypeProperties)

	if p.Auto.NotNil() {
		// See note#1
		dt = p.Auto.Result()
		goto POST_ToIso8601Utc
	}

	switch v.(type) {
	case DbNull:
		if p.Default.NotNil() { //db null should also execute default function
			// See note#1
			dt = p.Default.Result()
			goto POST_ToIso8601Utc
		}

		if !p.Nullable.True() {
			//TODO: change %s to named tag
			return nil, errors.New("value for this field %s cannot be null")
		} else {
			return v, nil
		}
	}

	dt, e = convert.ToIso8601Utc(v, p.ToSystemConversionClass())
	if e != nil {
		return nil, e
	}

POST_ToIso8601Utc:
	e = validateDbDateTimeByProperties(dt, p)
	if e != nil {
		return nil, e
	}

	return dt, nil
}

func (t dbTypeDateTime) GetDbTypeDisplayName() string {
	return t.DisplayName
}

func (t dbTypeDateTime) GetDbTypeId() int {
	return t.DbTypeId
}

// public DbTypeDateTimeProperties

func (t *DbTypeDateTimeProperties) IsPrimaryKey() bool {
	return false
}

// private
// Section: dbTypeDateTime

func (t dbTypeDateTime) defaultDbTypeProperties() DbTypePropertiesI {
	return defaultDbTypeDateTimeProperties()
}

func (t dbTypeDateTime) onCreateValidateFieldProperties(fieldProperties interface{}) error {
	//no primarykey at the moment so no validation code.
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

func validateDbDateTimeByProperties(dt dtIso8601Utc.Iso8601Utc, p *DbTypeDateTimeProperties) error {
	return nil
}

// Section: <dbTypeDateTimeValueOrFun>

func (t *dbTypeDateTimeValueOrFun) NotNil() bool {
	if t.number != nil || t.fun != nil {
		return true
	}
	return false
}

func (t *dbTypeDateTimeValueOrFun) Parse(tokens []parsers.Token, i int) (int, error) {
	return 0, nil
}

func (t *dbTypeDateTimeValueOrFun) Result() dtIso8601Utc.Iso8601Utc {
	return dtIso8601Utc.Iso8601Utc{}
}