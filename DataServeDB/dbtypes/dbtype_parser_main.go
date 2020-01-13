package dbtypes

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	db_rules "DataServeDB/dbsystem/rules"
	"DataServeDB/dbtypes/dbtype_props"
	"DataServeDB/parsers"
)

// Section:

var keywordsMap = map[string]parsers.ParserWordsMapItem{
	// NOTE: names don't always match struct name, because internal names are more descriptive of the nature of the struct.

	strings.ToUpper("PrimaryKey"):            parsers.NewParserWordsMapItem("PrimaryKeyable", parsers.ParserWordKeyword),
	strings.ToUpper("Nullable"):              parsers.NewParserWordsMapItem("Nullable", parsers.ParserWordKeyword),
	strings.ToUpper("UniqueIndex"):           parsers.NewParserWordsMapItem("Indexing", parsers.ParserWordKeyword),
	strings.ToUpper("SequentialUniqueIndex"): parsers.NewParserWordsMapItem("Indexing", parsers.ParserWordKeyword),
	strings.ToUpper("EmptyString"):           parsers.NewParserWordsMapItem("EmptyString", parsers.ParserWordKeyword),
	strings.ToUpper("Length"):                {"TypeLength", parsers.ParserWordKeyword, "", ":"},
	strings.ToUpper("Default"):               {"Default", parsers.ParserWordKeyword, "", ":"},
	strings.ToUpper("!"):                     parsers.NewParserWordsMapItem("!", parsers.ParserWordOperator),
}

// Section: public

//Keeping it, can be used to debug bugs.
func DebugPrintTokens(tokens []parsers.Token) {
	for _, t := range tokens {
		fmt.Println(t.Word)
	}
}

func ParseFieldProperties(s string) (fieldName string, dbType DbTypeI, dbTypeProperties DbTypePropertiesI, e error) {
	tokens := lex(s)
	//DebugPrintTokens(tokens)
	return parseAndProcess(tokens)
}

// Section: private

func applyTableFieldPropertyOperator(tfProperty interface{}, operator *string) error {
	//UniqueIndex and SequentialUniqueIndex is not technically an operator, but it is checked under operator. In a sense, they could be.

	//NOTE: negation is done wity interface '.Negate()', maybe other should be done same way?

	if operator != nil {
		switch strings.ToUpper(*operator) {

		//Indexing
		case "SEQUENTIALUNIQUEINDEX":
			return dbtype_props.SetIndexingType(tfProperty, dbtype_props.SequentialUniqueIndex)

		case "UNIQUEINDEX":
			return dbtype_props.SetIndexingType(tfProperty, dbtype_props.UniqueIndex)

		//Negation
		case "!":
			if n, e := dbtype_props.GetNegatable(tfProperty); e == nil {
				n.Negate()
				return nil
			} else {
				// TODO: suggestion: replace %s with text tag for db type property name, which will make error handling experience better.
				return errors.New("table field property '%s' does not support negation operator '!'")
			}

		//nullable
		case "NULLABLE":
			return dbtype_props.SetNullableFlag(tfProperty, dbtype_props.NullableTrue)

		//primaryKeyFlag
		case "PRIMARYKEY":
			if t, ok := tfProperty.(*dbtype_props.PrimaryKeyable); ok {
				t.IsPrimarykey = true
			}
			return nil

		} //end switch
	}

	return nil
}

func getTableFieldProperty(fieldProperties interface{}, parserItemPtr *parsers.ParserWordsMapItem) (*dbtype_props.DbTypePropertyParserItem, error) {

	s := reflect.Indirect(reflect.ValueOf(fieldProperties))

	if _, ok := s.Type().FieldByName(parserItemPtr.ExactName); !ok {
		return nil, errors.New("field property is not supported")
	}

	f := s.FieldByName(parserItemPtr.ExactName)

	if !f.IsValid() {
		return nil, errors.New("field property has accessibility problem (coding error please report this error)")
	}

	returnItem := dbtype_props.DbTypePropertyParserItem{
		DbTypeProperty:   f.Addr().Interface(),
		MustPreviousWord: parserItemPtr.MustPreviousWord,
		MustNextWord:     parserItemPtr.MustNextWord,
	}

	return &returnItem, nil
}

func isKeywordDuplicated(kwDupMap map[string]interface{}, kw string) bool {
	if _, exists := kwDupMap[kw]; exists {
		return true
	}
	kwDupMap[kw] = nil
	return false
}

func isValidParserItem(w string) (*parsers.ParserWordsMapItem, bool) {
	w = strings.ToUpper(w)
	if parserItem, exists := keywordsMap[w]; exists {
		return &parserItem, true
	}
	return nil, false
}

func IsNumeric(s string) bool {
	//TODO: move to utils
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func lex(s string) []parsers.Token {
	lexerRe := regexp.MustCompile(`\.\.|!|:|\(|,|\)|\b(?:string|int|PrimaryKey|UniqueIndex|SequentialUniqueIndex|Nullable|Length|NumberRange|EmptyString|default|auto|[a-zA-Z][a-zA-Z0-9_]+|[0-9_]+)\b|(?:['].*[']|["].*["])|[^ ;]+`)
	matchIndexes := lexerRe.FindAllStringIndex(s, -1)
	l := len(matchIndexes)
	var tokens []parsers.Token

	for i := 0; i < l; i++ {
		match := matchIndexes[i]

		tok := parsers.Token{
			StartPos: match[0],
			EndPos: match[1],
			Word: s[match[0]:match[1]],
		}

		if tok.Word != ".." {
			tokens = append(tokens, tok)
		} else {
			// if .. !number without space(s)
			//NOTE: don't need to check !number .. as it will be handled by parser, but needs testing to make sure.
			//TODO: test for )Nullable and does that needs to be separated by spaces?

			if i < l-1 { // i must be +2 or more of l
				match := matchIndexes[i+1]
				nextWord := s[match[0]:match[1]]

				if !IsNumeric(nextWord) && (tok.EndPos - match[0]) == 0 {
					tok.Word = tok.Word + nextWord
					tok.EndPos = match[1]
					tokens = append(tokens, tok)
					i += 2
					continue
				}
			}
			tokens = append(tokens, tok)
		}
	}

	return tokens
}

//Description: Parses field creation text.
func parseAndProcess(tokens []parsers.Token) (fieldName string, dbType DbTypeI, dbTypeProperties DbTypePropertiesI, e error) {

	//NOTE: parser exists at first error because this is executed on server side and thus it needs to exit if it is not doing any data processing.
	// However, on the client side, like in an editor, all errors should be detected and showed to the user for fixing.

	l := len(tokens)

	if l < 2 { // I believe minimum required: name and dbtype.
		// TODO: add table name to the message. suggestion: add a text tag for table name which can be processed later.
		e = errors.New("field creation code is missing required fields, at minium it needs field name and field type")
		return
	}

	fieldName = tokens[0].Word
	if !db_rules.TableFieldNameRulesCheck(fieldName) {
		e = fmt.Errorf("invalid table field name '%s'", fieldName)
		return
	}

	dbType, e = getDbType(tokens[1].Word)
	if e != nil {
		return
	}

	dbTypeProperties = dbType.defaultDbTypeProperties()
	if dbTypeProperties == nil {
		//TODO: log.
		//TODO: update with location of the code.
		panic("coding error, this cannot be nil, update code")
	}

	var keywordsDupMap = map[string]interface{}{}

	for i := 2; i < l; i++ {
		token := tokens[i]

		parserItemPtr, parserItemOk := isValidParserItem(token.Word)
		if !parserItemOk {
			//TODO: make this error more user friendly like show the wrong word and its location.
			e = errors.New("bad character or keyword")
			return
		}

		if parserItemPtr.Kind == parsers.ParserWordKeyword {
			if isKeywordDuplicated(keywordsDupMap, token.Word) {
				//TODO: make this error more user friendly like show the duplicated keyword and its location.
				e = errors.New("duplicate keyword")
				return
			}
		}

		//NOTE: operators don't do operations on their own either they are prefix or post fix, both are handled by their operated object.
		// So just continue.
		if parserItemPtr.Kind == parsers.ParserWordOperator {
			if l == i+1 {
				//TODO: make this error more user friendly like show the wrong operator and its location.
				e = errors.New("misplaced or mistyped operator")
				return
			}
			continue
		}

		var fieldPropertyParserItem *dbtype_props.DbTypePropertyParserItem
		fieldPropertyParserItem, e = getTableFieldProperty(dbTypeProperties, parserItemPtr)
		if e != nil {
			// removed token.Word, but it could be used in error message.
			//TODO: make this error more user friendly
			return
		}

		if fieldPropertyParserItem.MustNextWord != "" {
			if l <= i+1 || tokens[i+1].Word != fieldPropertyParserItem.MustNextWord { //NOTE: this needs to check if l is equal i+1 which means nothing is after it but last keyword requires operator after it.
				//TODO: make this error more user friendly?
				e = fmt.Errorf("table field property '%s' must have '%s'", token.Word, fieldPropertyParserItem.MustNextWord)
				return
			}
		}

		//NOTE: there is only one operator at the moment which can be check with i-1, otherwise, apply table field property.
		if i > 0 && tokens[i-1].Word == "!" { // Negation operator
			e = applyTableFieldPropertyOperator(fieldPropertyParserItem.DbTypeProperty, &tokens[i-1].Word)
			if e != nil {
				//TODO: make this error more user friendly?
				e = fmt.Errorf(e.Error(), token.Word)
				return
			}
		} else {
			e = applyTableFieldPropertyOperator(fieldPropertyParserItem.DbTypeProperty, &token.Word)
			if e != nil {
				//TODO: make this error more user friendly?
				e = fmt.Errorf(e.Error(), token.Word)
				return
			}
			if i < l-1 && tokens[i+1].Word == ":" { //i must be +2 of l, hence, i < l-1
				parsingType := dbtype_props.GetDbTypePropertyWithParser(fieldPropertyParserItem.DbTypeProperty)
				if parsingType != nil {
					i, e = parsingType.Parse(tokens, i+1)
					if e != nil {
						//TODO: make this error more user friendly?
						e = fmt.Errorf("error in parameters of table field property '%s', error: %s", token.Word, e.Error())
						return
					}
				}
			}
		}
	}

	//check primary key constraints: primary key must be index and not nullable.
	if e = dbType.onCreateValidateFieldProperties(dbTypeProperties); e != nil {
		return
	}

	return
}
