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
	"DataServeDB/dbsystem/builtinfuns"
	"DataServeDB/dbtypes/dbtype_props"
	"DataServeDB/parsers"
	"DataServeDB/utils/convert"
	"errors"
	"math"
)

// section: declarations

type dbTypeInt32 struct {
	//private to package
	dbTypeBase
}

type dbTypeInt32Fp = func() int32

type dbTypeInt32ValueOrFun struct {
	// Go doesn't have discriminated union so this will do.
	number *int32
	fun    dbTypeInt32Fp
}

type dbTypeInt32Properties struct {
	dbtype_props.Conversion
	dbtype_props.PrimaryKeyable
	dbtype_props.Nullable
	dbtype_props.Indexing
	dbtype_props.NumberRange
	Auto    dbTypeInt32ValueOrFun
	Default dbTypeInt32ValueOrFun
}

// public

var Int32 = dbTypeInt32{
	dbTypeBase{
		DbTypeId:    dbInt32,
		DisplayName: "int32",
	},
}

func (t dbTypeInt32) ConvertValue(v interface{}, dbTypeProperties interface{}) (interface{}, error) {

	//NOTE#1: auto and default are validated during property setting, but function returned values are not
	//	guaranteed so these result values are still validated here.

	var intval int32 //used intval because i is normally used for index
	var e error

	if v == nil {
		v = DbNull{}
	}

	p := getDbTypeInt32PropertiesIndirect(dbTypeProperties)

	if p == nil {
		//TODO: log.
		//TODO: update with location of the code.
		panic("coding error!")
	}

	if p.Auto.NotNil() {
		// See note#1
		intval = p.Auto.Result()
		goto POST_ToInt32
	}

	switch v.(type) {
	case DbNull:
		if p.Default.NotNil() { //db null should also execute default function
			// See note#1
			intval = p.Default.Result()
			goto POST_ToInt32
		}

		if !p.Nullable.True() {
			//TODO: change %s to named tag
			return nil, errors.New("value for this field %s cannot be null")
		} else {
			return v, nil
		}
	}

	intval, e = convert.ToInt32(v, p.ToSystemConversionClass())
	if e != nil {
		return nil, e
	}

POST_ToInt32:
	e = validateDbInt32ByProperties(intval, p)
	if e != nil {
		return nil, e
	}

	return intval, nil
}

func (t dbTypeInt32) GetDbTypeDisplayName() string {
	return t.DisplayName
}

func (t dbTypeInt32) GetDbTypeId() int {
	return t.DbTypeId
}

// public dbTypeInt32Properties

func (t *dbTypeInt32Properties) IsPrimaryKey() bool {
	return t.IsPrimarykey
}

// private
// Section: DbTypeInt32

func (t dbTypeInt32) defaultDbTypeProperties() DbTypePropertiesI {
	return defaultDbTypeInt32Properties()
}

func (t dbTypeInt32) onCreateValidateFieldProperties(fieldProperties interface{}) error {
	if fp, ok := fieldProperties.(*dbTypeInt32Properties); ok {
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

// Section: dbTypeInt32Properties

func defaultDbTypeInt32Properties() *dbTypeInt32Properties {
	return &dbTypeInt32Properties{
		Conversion:     dbtype_props.NewConversion(),
		PrimaryKeyable: dbtype_props.PrimaryKeyable{IsPrimarykey: false},
		Nullable:       dbtype_props.Nullable{State: dbtype_props.NullableFalseDefault},
		Indexing:       dbtype_props.Indexing{IndexType: dbtype_props.IndexingNone, Supports: dbtype_props.UniqueIndex | dbtype_props.SequentialUniqueIndex},
		NumberRange:    dbtype_props.NewNumberRange(math.MinInt32, math.MaxInt32),
		Auto:           dbTypeInt32ValueOrFun{},
		Default:        dbTypeInt32ValueOrFun{},
	}
}

func getDbTypeInt32PropertiesIndirect(p interface{}) *dbTypeInt32Properties {

	switch p := p.(type) {
	case dbTypeInt32Properties:
		return &p
	case *dbTypeInt32Properties:
		return p
	}

	//TODO: panic with code location.
	panic("Coding error, this should not happen!")
}

func validateDbInt32ByProperties(intval int32, p *dbTypeInt32Properties) error {

	if intval < int32(p.NumberRange.Min) {
		return errors.New("value is less than minimum limit")
	}

	if intval > int32(p.NumberRange.Max) {
		return errors.New("value is more than maximum limit")
	}

	return nil
}

// Section: <dbTypeInt32ValueOrFun>

func getDbInt32Fun(funName string, params ...int32) dbTypeInt32Fp {

	// function names are case sensitive
	switch funName {

	case "Increment":
		if len(params) == 2 {
			return builtinfuns.IncrementInt32(params[0], params[1])
		}

	}

	return nil
}

func (t *dbTypeInt32ValueOrFun) NotNil() bool {
	if t.number != nil || t.fun != nil {
		return true
	}
	return false
}

func (t *dbTypeInt32ValueOrFun) Parse(tokens []parsers.Token, i int) (int, error) {

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

		//could be number or function, function doesn't start with digit.
		num, e := convert.ToInt32(tokens[i].Word, convert.Lossless)
		if e == nil {
			// it is number
			t.number = &num
			break
		}

		//conversion failed, but it could be fun

		//store i-1 to return in case of failure.
		failReturnPos = i - 1

		if l <= i+2 {
			//TODO: make this error more user friendly?
			return failReturnPos, errors.New("bad fun name") //-1 due sending last token back to parent parser.
		}

		if tokens[i+1].Word != "(" {
			//TODO: make this error more user friendly?
			return failReturnPos, errors.New("bad fun name")
		}

		funName := tokens[i].Word
		i += 2 //skipping '(', which already checked above

		var intParams []int32 //NOTE: at the moment only supports int32 params

		for ; i < l; i++ {
			// next could be number, comma, or comma+number, or closing paren.
			if tokens[i].Word == ")" {
				goto FUN_ENDED
			}

			if tokens[i].Word == "," {
				//TODO: should be only one after number.
				continue
			}

			//this is function parameter, must be number. so on error fail
			//TODO: should only accept number, lossless can also accept string which converts to number properly like "1"
			//NOTE (for above): tested it is not accepting '1' or "1", which is good but need to see why as lossless will accept this.
			num, e := convert.ToInt32(tokens[i].Word, convert.Lossless)
			if e == nil {
				// it is number
				// TODO: check how append is giving type error? Note: don't remember why i wrote this but keeping it here.
				intParams = append(intParams, num)
			} else {
				// error in int fun call
				//TODO: make this error more user friendly?
				return failReturnPos, errors.New("bad fun code")
			}
		}

		if tokens[i-1].Word != ")" {
			//TODO: make this error more user friendly?
			return failReturnPos, errors.New("bad function name")
		}

	FUN_ENDED:
		fun := getDbInt32Fun(funName, intParams...)
		if fun == nil {
			//TODO: better error message if function name exists but params are different.
			return failReturnPos, errors.New("function not found") //NOTE: i-1 is not correct here; function not found, then move to next item. i+2 for ()
		}
		t.fun = fun
		// no need to i++ as parent parser will move to next keyword
		break
	}

	//both of them cannot be nil
	if t.number == nil && t.fun == nil {
		//TODO: better error message
		return failReturnPos, errors.New("malformed Default/Auto property") //failReturnPos returns position of function name token for next parsing word, which is correct.
	}

	return i, nil
}

func (t *dbTypeInt32ValueOrFun) Result() int32 {
	//Note: return can be string.

	if t.number != nil {
		return *t.number
	}

	if t.fun != nil {
		return t.fun()
	}

	//NOTE: both shouldn't be nil, something is wrong in the code.
	//TODO: panic with code location.
	panic("Coding error, this should not be reached!")
}

// </dbTypeInt32ValueOrFun>
