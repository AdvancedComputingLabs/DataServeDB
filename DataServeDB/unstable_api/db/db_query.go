package db

import (
	"DataServeDB/commtypes"
	"DataServeDB/dbtable"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"net/http"
)

type QueryOp int

const (
	OpNone = iota // starting with 0
	OpIS
	OpOR
	OpAND
)

const (
	Field = "field"
	Table = "table"
)

type ruleInfo struct {
	tableName string
	fieldName string
	operation QueryOp
	next      *ruleInfo
}

type Query struct {
	ItemLabel string
	ItemType  string
	ItemValue []byte // json Converted
	Rules     *ruleInfo
	Children  []Query
}

func (t *DB) TablesQueryGet(dbReqCtx *commtypes.DbReqContext, query Query) (resultHttpStatus int, resultContent []byte, resultErr error) {
	table, err := t.GetTable(query.ItemLabel)
	if err != nil {
		resultErr = err
		return
	}
	query, err = t.verifyQuery(query, table)
	if err != nil {
		resultErr = err
		return
	}
	result, err := t.processQuery(query)
	if err != nil {
		resultErr = err
		return
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		resultErr = fmt.Errorf("error decoding result: %V", err)
		//TODO: make error result more user friendly.
		return
	}
	return http.StatusOK, jsonBytes, resultErr
}

func (t *DB) verifyQuery(query Query, table *dbtable.DbTable) (Query, error) {
	for i, value := range query.Children {
		if _, found := table.TblMain.TableFieldsMetaData.IsField(value.ItemLabel); found {
			/* TODO make itemtype as macro */
			query.Children[i].ItemType = Field
		} else if tbl, e := t.GetTable(value.ItemLabel); e == nil {
			child, err := t.verifyQuery(query.Children[i], tbl)
			if err != nil {
				return Query{}, err
			}
			query.Children[i] = child
			query.Children[i].ItemType = Table
		} else if value.ItemLabel == "$WHERE" {
			rules := getRules(value.ItemValue)
			query.Rules = rules
			return query, nil
		} else {
			return Query{}, fmt.Errorf("Query field '%s' does not exist in database", value.ItemLabel)
		}
	}
	return query, nil
}
func (t *DB) processQuery(query Query) (result []dbtable.TableRow, err error) {
	rows := []dbtable.TableRow{}
	ArrRows := map[string][]dbtable.TableRow{}
	table, err := t.GetTable(query.ItemLabel)
	if err != nil {
		return
	}
	spec := getSpec(query)
	if spec != -1 {
		rows, err = table.GetTableRows(string(query.Children[spec].ItemValue))
		if err != nil {
			return
		}
	} else {
		rows, err = table.GetTableRows()
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	for _, row := range rows {
		res := dbtable.TableRow{}
		if query.Children == nil {
			result = append(result, row)
			continue
		}
		for _, child := range query.Children {
			if child.ItemType == Field {
				res[child.ItemLabel] = row[child.ItemLabel]
			} else if child.ItemType == Table {
				var tempRow dbtable.TableRow
				temp := []dbtable.TableRow{}
				tbres := []dbtable.TableRow{}
				if _, ok := ArrRows[query.ItemLabel]; !ok {
					ArrRows[query.ItemLabel] = rows
				}
				for rule := child.Rules; rule != nil; rule = rule.next {
					if _, ok := ArrRows[rule.tableName]; !ok {
						rows1, err := t.getRowsBytableName(rule.tableName)
						if err != nil {
							return nil, err
						}
						ArrRows[rule.tableName] = rows1
					}
				}
				for rule := child.Rules; rule != nil; rule = rule.next {
					if rule.operation == OpAND {
						continue
					}
					if rule.tableName == query.ItemLabel {
						tempRow = row
						if rule.operation == OpIS {
							for _, rv := range ArrRows[rule.next.tableName] {
								if tempRow[rule.fieldName] == rv[rule.fieldName] {
									temp = append(temp, rv)
								}
							}
						}
					}
					if rule.tableName == child.ItemLabel {
						if rule.operation == OpIS {
							for _, prv := range ArrRows[rule.tableName] {
								for _, rv := range temp {
									if rv[rule.fieldName] == prv[rule.fieldName] {
										tbres = append(tbres, prv)
									}
								}
							}
						}
					}
				}
				res[child.ItemLabel] = tbres
			}
		}
		result = append(result, res)
	}
	return result, nil
}
func getSpec(query Query) int {
	for i, v := range query.Children {
		if v.ItemType == Field {
			if v.ItemValue != nil {
				return i
			}
		}
	}
	return -1
}

func getRules(b []byte) (rulse *ruleInfo) {
	rule := ruleInfo{}
	// "Properties": [
	//   {
	//     "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
	//   }
	// ]
	var re = regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)`)
	if byt := re.Find(b); byt != nil {
		arr := strings.Split(string(byt), ".")
		rule.tableName = arr[0]
		rule.fieldName = arr[1]
		rule.operation = getOpr(b[len(byt)+1:])
		rule.next = getRules(b[(len(byt) + 3):])
		return &rule
	}
	return nil

}
func (t *DB) getRowsBytableName(tableName string) (rows []dbtable.TableRow, err error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	rows, err = table.GetTableRows()
	if err != nil {
		return nil, err
	}
	return
}

func getOpr(str []byte) QueryOp {
	operators := map[string]QueryOp{
		"IS":  OpIS,
		"OR":  OpOR,
		"AND": OpAND,
	}
	var opre = regexp.MustCompile(`(?m)([A-Z]{2,5})`)
	opr := opre.Find(str)
	if v, ok := operators[string(opr)]; ok {
		return v
	}
	return OpNone
}
