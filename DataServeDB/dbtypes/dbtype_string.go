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
	"errors"

	"DataServeDB/dbtypes/dbtype_props"
	"DataServeDB/parsers"
	"DataServeDB/utils/convert"
)

// section: declarations

type dbTypeString struct {
	//private to package
	dbTypeBase
}

type dbTypeStringFp = func() string

type dbTypeStringValueOrFun struct {
	// Go doesn't have discriminated union so this will do.
	text *string
	fun  dbTypeStringFp
}

type dbTypeStringProperties struct {
	dbtype_props.Conversion
	dbtype_props.PrimaryKeyable
	dbtype_props.Nullable
	dbtype_props.Indexing
	dbtype_props.TypeLength
	Auto    dbTypeStringValueOrFun
	Default dbTypeStringValueOrFun
}

// public

var String = dbTypeString{
	dbTypeBase{
		DbTypeId:    dbString,
		DisplayName: "string",
	},
}

func (t dbTypeString) ConvertValue(v interface{}, dbTypeProperties interface{}) (interface{}, error) {

	//NOTE#1: auto and default are validated during property setting, but function returned values are not
	//	guaranteed so these result values are still validated here.

	var s string
	var e error

	if v == nil {
		v = DbNull{}
	}

	//TODO: panic if properties are nil, coding error.
	p := getDbTypeStringPropertiesIndirect(dbTypeProperties)

	if p.Auto.NotNil() {
		// See note#1
		s = p.Auto.Result()
		goto POSTTOSTRING
	}

	switch v.(type) {
	case DbNull:
		if p.Default.NotNil() { //db null should also execute default function
			// See note#1
			s = p.Default.Result()
			goto POSTTOSTRING
		}

		if !p.Nullable.True() {
			return nil, errors.New("value for column '%s' cannot be null")
		} else {
			return v, nil
		}
	}

	s, e = convert.ToString(v, p.ToSystemConversionClass())
	if e != nil {
		return nil, e
	}

POSTTOSTRING:
	if e = validateDbTypeStringProperties(s, p); e != nil {
		return nil, e
	}

	return s, nil
}

func (t dbTypeString) GetDbTypeId() int {
	return t.DbTypeId
}

func (t dbTypeString) GetDbTypeDisplayName() string {
	return t.DisplayName
}

// public dbTypeStringProperties

func (t *dbTypeStringProperties) IsNullable() bool {
	return t.Nullable.True()
}

func (t *dbTypeStringProperties) IsPrimaryKey() bool {
	return t.IsPrimarykey
}

// private
// Section: DbTypeString

func (t dbTypeString) defaultDbTypeProperties() DbTypePropertiesI {
	return defaultDbTypeStringProperties()
}

func (t dbTypeString) onCreateValidateFieldProperties(fieldProperties interface{}) error {
	if fp, ok := fieldProperties.(*dbTypeStringProperties); ok {
		if fp.IsPrimarykey {
			if fp.IndexType == dbtype_props.IndexingNone {
				fp.IndexType = dbtype_props.UniqueIndex
			}
			if e := dbtype_props.PrimaryKeyConstraintCheck(fp.Nullable, fp.Indexing); e != nil {
				return e
			}
		}
	} else {
		//TODO: code error
	}
	return nil
}

// Section: dbTypeStringProperties

func defaultDbTypeStringProperties() *dbTypeStringProperties {
	return &dbTypeStringProperties{
		Conversion:     dbtype_props.NewConversion(),
		PrimaryKeyable: dbtype_props.PrimaryKeyable{IsPrimarykey: false},
		Nullable:       dbtype_props.Nullable{State: dbtype_props.NullableFalseDefault},
		Indexing:       dbtype_props.Indexing{IndexType: dbtype_props.IndexingNone, Supports: dbtype_props.UniqueIndex | dbtype_props.SequentialUniqueIndex},
		TypeLength:     dbtype_props.NewTypeLength(0, 4000),
		Auto:           dbTypeStringValueOrFun{},
		Default:        dbTypeStringValueOrFun{},
	}
}

func getDbTypeStringPropertiesIndirect(p interface{}) *dbTypeStringProperties {

	switch p := p.(type) {
	case dbTypeStringProperties:
		return &p
	case *dbTypeStringProperties:
		return p
	}

	//TODO: panic with code location.
	panic("Coding error, this should not happen!")
}

func validateDbTypeStringProperties(s string, p *dbTypeStringProperties) error {
	//TODO: make errors more user friendly

	sLen := uint64(len(s))

	if sLen < p.TypeLength.Min {
		return errors.New("length is less than minimum size")
	}

	if sLen > p.TypeLength.Max {
		return errors.New("length is more than maximum size")
	}

	return nil
}

// Section: <dbTypeStringValueOrFun>

func getDbStringFun(funName string) dbTypeStringFp {
	// function names are case sensitive
	switch funName {
	case "HelloWorld":
		//return FnHelloWorld //TODO: later
	}

	return nil
}

func (t *dbTypeStringValueOrFun) NotNil() bool {
	if t.text != nil || t.fun != nil {
		return true
	}
	return false
}

func (t *dbTypeStringValueOrFun) Parse(tokens []parsers.Token, i int) (int, error) {

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

	failReturnPos := i

	for ; i < l; i++ {
		if tokens[i].Word == ":" {
			continue
		}

		//token could be text or function.

		failReturnPos = i - 1 //one keyword back for returning

		if tokens[i].Word[0] == '"' || tokens[i].Word[0] == '\'' {
			//token is text
			textLen := len(tokens[i].Word)
			if textLen < 2 || tokens[i].Word[textLen-1] != tokens[i].Word[0] {
				return failReturnPos, errors.New("text must have proper start and end string quotes, see help docs") //-1 due sending last token back to parent parser.
			}
			if textLen == 2 {
				t.text = new(string)
				break
			} else {
				text := tokens[i].Word[1 : textLen-1]
				t.text = &text //TODO: check is this copy? references for this token shouldn't stay in memory.
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

			fun := getDbStringFun(tokens[i].Word)
			if fun == nil {
				//TODO: test return position is correct.
				//TODO: make this error more user friendly?
				return i + 2, errors.New("fun not found") //NOTE: i-1 is not correct here; function not found, then move to next item. i+2 for ()
			}

			t.fun = fun

			//fmt.Println("#", t.fun())
			i = i + 2
			break
		}
	}

	//TODO: return from word before possible fun name word position.
	//both of them cannot be nil
	if t.text == nil && t.fun == nil {
		//TODO: better error message
		return i, errors.New("malformed Default property")
	}

	return i, nil
}

func (t *dbTypeStringValueOrFun) Result() string {
	//Note: return can be string.

	if t.text != nil {
		return *t.text
	}

	if t.fun != nil {
		return t.fun()
	}

	//NOTE: both shouldn't be nil, something is wrong in the code.
	//TODO: panic with code location.
	panic("Coding error, this should not be reached!")
}

// </dbTypeStringValueOrFun>
