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
	// NOTE: see bool file.
	//private to package
	dbTypeBase
}

type dbTypeDateTimeFp = func() dtIso8601Utc.Iso8601Utc

type dbTypeDateTimeValueOrFun struct {
	// Go doesn't have discriminated union so this will do.
	datetime *dtIso8601Utc.Iso8601Utc
	fun      dbTypeDateTimeFp
}

type DbTypeDateTimeProperties struct {
	// NOTE: see bool file.
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
	// NOTE: see bool file.
	return t.DisplayName
}

func (t dbTypeDateTime) GetDbTypeId() int {
	// NOTE: see bool file.
	return t.DbTypeId
}

// public DbTypeDateTimeProperties

func (t *DbTypeDateTimeProperties) IsNullable() bool {
	// NOTE: see bool file.
	return t.Nullable.True()
}

func (t *DbTypeDateTimeProperties) IsPrimaryKey() bool {
	// NOTE: see bool file.
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

func getDbDateTimeFun(funName string) dbTypeDateTimeFp {

	// function names are case sensitive
	switch funName {

	case "Now":
		return dtIso8601Utc.Iso8601UtcNow

	}

	return nil
}

func (t *dbTypeDateTimeValueOrFun) NotNil() bool {
	if t.datetime != nil || t.fun != nil {
		return true
	}
	return false
}

func (t *dbTypeDateTimeValueOrFun) Parse(tokens []parsers.Token, i int) (int, error) {
	if tokens == nil {
		//TODO: log.
		//TODO: update with location of the code.
		panic("coding error, shouldn't happen")
	}

	l := len(tokens)

	if l == 0 {
		//QUESTION: empty token, should be error?
		return i, nil
	}

	if i > l {
		//TODO: log.
		//TODO: update with location of the code.
		panic("coding error, shouldn't happen")
	}

	colonCount := 0
	failReturnPos := i

	for ; i < l; i++ {
		if tokens[i].Word == ":" && colonCount == 0 {
			colonCount++
			failReturnPos = i
			continue
		}

		failReturnPos = i - 1 //one keyword back for returning

		if tokens[i].Word[0] == '"' || tokens[i].Word[0] == '\'' {
			//token is text
			textLen := len(tokens[i].Word)
			if textLen < 2 || tokens[i].Word[textLen-1] != tokens[i].Word[0] {
				return failReturnPos, errors.New("datetime must have proper start and end string quotes, see help docs") //-1 due sending last token back to parent parser.
			}
			if textLen == 2 {
				t.datetime = &dtIso8601Utc.Iso8601Utc{}
				break
			} else {
				text := tokens[i].Word[1 : textLen-1]
				dt, e := dtIso8601Utc.Iso8601UtcFromString(text)
				if e != nil {
					return failReturnPos, e
				}
				t.datetime = &dt
				break
			}
		} else {
			//TODO: check if function name is valid.

			if l <= i+2 {
				//TODO: make this error more user friendly?
				return failReturnPos, errors.New("bad fun name") //-1 due sending last token back to parent parser.
			}

			if !(tokens[i+1].Word == "(" || tokens[i+2].Word == ")") {
				//TODO: make this error more user friendly?
				return failReturnPos, errors.New("bad fun name")
			}

			fun := getDbDateTimeFun(tokens[i].Word)
			if fun == nil {
				//TODO: test return position is correct.
				//TODO: make this error more user friendly?
				return i + 2, errors.New("fun not found") //NOTE: i-1 is not correct here; function not found, then move to next item. i+2 for ()
			}

			t.fun = fun

			i = i + 2
			break
		}
	}

	//TODO: return from word before possible fun name word position.
	//both of them cannot be nil
	if t.datetime == nil && t.fun == nil {
		//TODO: better error message
		return i, errors.New("malformed Default property")
	}

	return i, nil
}

func (t *dbTypeDateTimeValueOrFun) Result() dtIso8601Utc.Iso8601Utc {
	if t.datetime != nil {
		return *t.datetime
	}

	if t.fun != nil {
		return t.fun()
	}

	//NOTE: both shouldn't be nil, something is wrong in the code.
	//TODO: panic with code location.
	panic("Coding error, this should not be reached!")
}
