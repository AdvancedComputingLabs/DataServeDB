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

	"DataServeDB/dbsystem/systypes/guid"
	"DataServeDB/dbtypes/dbtype_props"
	"DataServeDB/parsers"
	"DataServeDB/utils/convert"
)

// section: declarations

type dbTypeGuid struct {
	//private to package
	dbTypeBase
}

type dbTypeGuidFp = func() guid.Guid

type dbTypeGuidValueOrFun struct {
	// Go doesn't have discriminated union so this will do.
	guid *guid.Guid
	fun    dbTypeGuidFp
}

type DbTypeGuidProperties struct {
	dbtype_props.PrimaryKeyable
	dbtype_props.Indexing
	dbtype_props.Conversion
	dbtype_props.Nullable
	Auto    dbTypeGuidValueOrFun
	Default dbTypeGuidValueOrFun
}

// public

var Guid = dbTypeGuid{
	dbTypeBase{
		DbTypeId:    dbGuid,
		DisplayName: "guid",
	},
}

// public

func (t dbTypeGuid) ConvertValue(v interface{}, dbTypeProperties interface{}) (interface{}, error) {

	//NOTE#1: auto and default are validated during property setting, but function returned values are not
	//	guaranteed so these result values are still validated here.

	var g guid.Guid
	var e error

	if v == nil {
		v = DbNull{}
	}

	p := getDbTypeGuidPropertiesIndirect(dbTypeProperties)

	if p.Auto.NotNil() {
		// See note#1
		g = p.Auto.Result()
		goto POST_ToGuid
	}

	switch v.(type) {
	case DbNull:
		if p.Default.NotNil() { //db null should also execute default function
			// See note#1
			g = p.Default.Result()
			goto POST_ToGuid
		}

		if !p.Nullable.True() {
			//TODO: change %s to named tag
			return nil, errors.New("value for this field %s cannot be null")
		} else {
			return v, nil
		}
	}

	g, e = convert.ToGuid(v, p.ToSystemConversionClass())
	if e != nil {
		return nil, e
	}

POST_ToGuid:
	e = validateDbGuidByProperties(g, p)
	if e != nil {
		return nil, e
	}

	return g, nil
}

func (t dbTypeGuid) GetDbTypeDisplayName() string {
	return t.DisplayName
}

func (t dbTypeGuid) GetDbTypeId() int {
	return t.DbTypeId
}

// public DbTypeGuidProperties

func (t *DbTypeGuidProperties) IsPrimaryKey() bool {
	return t.IsPrimarykey
}

// private
// Section: DbTypeGuid

func (t dbTypeGuid) defaultDbTypeProperties() DbTypePropertiesI {
	return defaultDbTypeGuidProperties()
}

func (t dbTypeGuid) onCreateValidateFieldProperties(fieldProperties interface{}) error {
	if fp, ok := fieldProperties.(*DbTypeGuidProperties); ok {
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

// Section: DbTypeGuidProperties

func defaultDbTypeGuidProperties() *DbTypeGuidProperties {
	return &DbTypeGuidProperties{
		PrimaryKeyable: dbtype_props.PrimaryKeyable{IsPrimarykey: false},
		Indexing:       dbtype_props.Indexing{IndexType: dbtype_props.IndexingNone, Supports: dbtype_props.UniqueIndex},
		Nullable: dbtype_props.Nullable{State: dbtype_props.NullableFalseDefault},
	}
}

func getDbTypeGuidPropertiesIndirect(p interface{}) *DbTypeGuidProperties {

	switch p := p.(type) {
	case DbTypeGuidProperties:
		return &p
	case *DbTypeGuidProperties:
		return p
	}

	//TODO: log
	//TODO: panic with code location.
	panic("Coding error, this should not happen!")
}

func validateDbGuidByProperties(g guid.Guid, p *DbTypeGuidProperties) error {

	return nil
}

// Section: <dbTypeGuidValueOrFun>

func NewGuid() guid.Guid {
	return *guid.NewGuid()
}

func getDbGuidFun(funName string) dbTypeGuidFp {

	// function names are case sensitive
	switch funName {

	case "NewGuid":
		return NewGuid

	}

	return nil
}

func (t *dbTypeGuidValueOrFun) NotNil() bool {
	if t.guid != nil || t.fun != nil {
		return true
	}
	return false
}

func (t *dbTypeGuidValueOrFun) Parse(tokens []parsers.Token, i int) (int, error) {
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
				t.guid = guid.NewGuid()
				break
			} else {
				text := tokens[i].Word[1 : textLen-1]
				g, e := guid.ParseString(text)
				if e != nil {
					return failReturnPos, e
				}
				t.guid = g
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

			fun := getDbGuidFun(tokens[i].Word)
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
	if t.guid == nil && t.fun == nil {
		//TODO: better error message
		return i, errors.New("malformed Default property")
	}

	return i, nil
}

func (t *dbTypeGuidValueOrFun) Result() guid.Guid {
	if t.guid != nil {
		return *t.guid
	}

	if t.fun != nil {
		return t.fun()
	}

	//NOTE: both shouldn't be nil, something is wrong in the code.
	//TODO: panic with code location.
	panic("Coding error, this should not be reached!")
}