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

package dbtype_props

import (
	"errors"
	"strconv"

	"DataServeDB/parsers"
)

// Section: Interfaces & General Types

type DbTypePropertyWithParser interface {
	Parse([]parsers.Token, int) (int, error) // first int is current token index, second is return token index (note: return i will be next token)
}

type DbTypePropertyParserItem struct {
	DbTypeProperty   interface{}
	MustPreviousWord string
	MustNextWord     string
}

// Section: General

func PrimaryKeyConstraintCheck(n Nullable, i Indexing) error {
	if n.True() {
		return errors.New("primary key cannot be Nullable")
	}
	if i.IndexType == IndexingNone {
		return errors.New("primary key must be indexed")
	}
	return nil
}

// private
func GetDbTypePropertyWithParser(i interface{}) DbTypePropertyWithParser {
	if t, ok := i.(DbTypePropertyWithParser); ok {
		return t
	}
	return nil
}

func SetIndexingType(t interface{}, indexType indexingType) error {
	if t, ok := t.(*Indexing); ok {
		if !(t.Supports&indexType == indexType) {
			//TODO:
			return errors.New("")
		}
		t.IndexType = indexType
		return nil
	}
	return errors.New("coding error shouldn't happen (please report)")
}

func SetNullableFlag(t interface{}, f nullableFlag) error {
	if t, ok := t.(*Nullable); ok {
		t.State = f
		return nil
	}
	return errors.New("coding error shouldn't happen (please report)")
}

// Section: Indexing
// Indexing: IndexingNone(default) | UniqueIndex | SequentialUniqueIndex

type indexingType int

const (
	IndexingNone          indexingType = 2
	UniqueIndex                        = 4
	SequentialUniqueIndex              = 8
)

type Indexing struct {
	IndexType indexingType
	Supports  indexingType //this is compile time so can be private
}

// Section: TypeLength

type TypeLength struct {
	minLimit uint64
	maxLimit uint64
	Min      uint64
	Max      uint64
}

func NewTypeLength(minLimit uint64, maxLimit uint64) TypeLength {
	return TypeLength{minLimit: minLimit, maxLimit: maxLimit, Min: minLimit, Max: maxLimit}
}

func (t *TypeLength) Parse(tokens []parsers.Token, i int) (int, error) {
	//supports: 0..10, 10 (min for the type and max), ..10 same as before, 2.. (min 2 and max till max for the type).
	//defaults: min and max are based on db type, which is supposed to be set during db type initialization.

	//pre checks:
	// 1) tokens nil; panic, must be code error
	// 2) tokens len is 0; exit
	// 3) tokens len is less than i; panic, must be code error
	// 4) i is at the end, exit

	//code checks:
	// 1) min and max are according to the limits

	//post checks:
	// 1) min <= max

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

	//0 min, 1 max
	atMinOrMax := 0
	colonCount := 0
	firstNum := false
	failReturnPos := i

	for ; i < l; i++ {
		if tokens[i].Word == ":" && colonCount == 0 {
			// if ':' is last character then this is error
			if i+1 == l {
				//TODO: make this error more user friendly?
				return i, errors.New("error during parsing Length limits") // : is not for next character so not i-1
			}
			colonCount++
			failReturnPos = i
			continue
		}

		failReturnPos = i - 1 //one keyword back for returning

		// !WARNING: if lexing includes '.' then following code will break.
		if tokens[i].Word == ".." && atMinOrMax == 0 {
			atMinOrMax++
			continue
		}

		num, e := strconv.ParseUint(tokens[i].Word, 10, 64)
		if e != nil {
			if firstNum {
				//then is could be Range: 5 [Next Keyword]
				i--
				goto FINISH
			}
			//TODO: make this error more user friendly?
			return failReturnPos, errors.New("error during parsing Length limits") // i-1 because this item might be needed by parent parser.
		}

		if atMinOrMax == 0 {
			if num < t.minLimit {
				//TODO: central error handling with multi language support.
				//TODO: make this error more user friendly?
				return failReturnPos, errors.New("minimum length value is less than allowed minimum length")
			}
			t.Min = num
			firstNum = true //safer to check this way as t.Min can be 0
		} else if atMinOrMax == 1 {
			if num > t.maxLimit {
				//TODO: central error handling with multi language support.
				//TODO: make this error more user friendly?
				return failReturnPos, errors.New("maximum length value is more than allowed maximum length")
			}
			t.Max = num
			break
		}
	}

FINISH:
	failReturnPos = i //i-1 is not correct for the location, hence, updated here. i would have been correct but did this here for keeping with the semantics

	if firstNum && atMinOrMax == 0 {
		t.Max = t.Min
	}

	if t.Max < t.Min { // min max constraint check
		return failReturnPos, errors.New("max cannot be less than min")
	}

	return i, nil
}

// Section: Nullable
// Nullable: NullableDefault | Nullable | NotNullable

type nullableFlag int

const (
	NullableFalseDefault nullableFlag = iota
	NullableTrue
	NullableFalse
)

type Nullable struct {
	State nullableFlag
}

func (t *Nullable) Negate() {
	t.State = NullableFalse
}

func (t *Nullable) True() bool {
	return t.State == NullableTrue
}

// Section: PrimaryKeyable

type PrimaryKeyable struct {
	IsPrimarykey bool
}

// Section: NumberRange

type NumberRange struct {
	minLimit int64
	maxLimit int64
	Min      int64
	Max      int64
}

func NewNumberRange(minLimit int64, maxLimit int64) NumberRange {
	return NumberRange{minLimit: minLimit, maxLimit: maxLimit, Min: minLimit, Max: maxLimit}
}

func (t *NumberRange) Parse(tokens []parsers.Token, i int) (int, error) {
	return i, nil
}
