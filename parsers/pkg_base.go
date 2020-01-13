package parsers

//NOTE: main parser(s) does not contain any specific information about tables, their fields, or db types.

type ParserWordKind int

const (
	ParserWordUnknown = iota
	ParserWordKeyword
	ParserWordOperator
	ParserWordPostfix
)

type ParserWordsMapItem struct {
	ExactName string
	Kind ParserWordKind
	MustPreviousWord string
	MustNextWord string
}

type Token struct {
	Word     string
	StartPos int
	EndPos   int
}

// public
func NewParserWordsMapItem(exactName string, kind ParserWordKind) ParserWordsMapItem {
	return ParserWordsMapItem{exactName, kind, "", ""}
}
