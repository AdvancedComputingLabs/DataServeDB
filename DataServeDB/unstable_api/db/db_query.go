package db

import (
	"DataServeDB/commtypes"
	"DataServeDB/dbtable"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	// "DataServeDB/unstable_api/runtime"
	"net/http"
)

type ruleInfo struct {
	tableName string
	fieldName string
	operation int
	next      *ruleInfo
}
type ruleTabel struct {
	pTable   ruleInfo
	lnkTable ruleInfo
	// pkField   string
	// linkField string
}

type Query struct {
	ItemLabel string
	ItemType  string
	ItemValue []byte // json Converted
	Rules     *ruleInfo
	Children  []Query
}

func (t *DB) TablesQueryGet(dbReqCtx *commtypes.DbReqContext, query Query) (resultHttpStatus int, resultContent []byte, resultErr error) {
	// result := []dbtable.TableRow{}
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
	// sort.SliceStable(query.Children, func(i, j int) bool { return query.Children[i].ItemType == "field" })
	/*********************************************/
	/* TO DO :- filter out the specific row if the userId mentioned,
	}*/

	result, err := t.processQuery(query)
	if err != nil {
		resultErr = err
		return
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		resultErr = err
		//TODO: make error result more user friendly.
		return
	}
	return http.StatusOK, jsonBytes, resultErr

	//return 0, nil, nil
}

func (t *DB) verifyQuery(query Query, table *dbtable.DbTable) (Query, error) {
	for i, value := range query.Children {
		if _, found := table.TblMain.TableFieldsMetaData.IsField(value.ItemLabel); found {
			/* TODO make itemtype as macro */
			e := found
			_ = e
			query.Children[i].ItemType = "field"
		} else if tbl, e := t.GetTable(value.ItemLabel); e == nil {
			child, err := t.verifyQuery(query.Children[i], tbl)
			if err != nil {
				return Query{}, err
			}
			query.Children[i] = child
			query.Children[i].ItemType = "table"
		} else if value.ItemLabel == "$WHERE" {
			// process string
			rules := getRules(value.ItemValue)
			// fmt.Println(*rules)
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
			if child.ItemType == "field" {
				res[child.ItemLabel] = row[child.ItemLabel]
			} else if child.ItemType == "table" {
				var tempRow dbtable.TableRow
				temp := []dbtable.TableRow{}
				tbres := []dbtable.TableRow{}
				// if child.Rules[0].pTable.tableName != query.ItemLabel {
				// 	return nil, fmt.Errorf("primary table not fount in rule")
				// }
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
				//   "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
				for rule := child.Rules; rule != nil; rule = rule.next {
					if rule.operation == 3 {
						continue
					}
					if rule.tableName == query.ItemLabel {
						tempRow = row
						if rule.operation == 1 {
							for _, rv := range ArrRows[rule.next.tableName] {
								if tempRow[rule.fieldName] == rv[rule.fieldName] {
									temp = append(temp, rv)
								}
							}
						}
					}
					if rule.tableName == child.ItemLabel {
						if rule.operation == 1 {
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
		if v.ItemType == "field" {
			if v.ItemValue != nil {
				return i
			}
		}
	}
	return -1
}

func getRules(b []byte) (rulse *ruleInfo) {
	rule := ruleInfo{}
	// tblRule := []ruleTabel{}
	// "Properties": [
	//   {
	//     "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
	//   }
	// ]

	// var re = regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)( IS )([A-z]*[.][A-z]*)`)

	// for _, match := range re.FindAllString(str, -1) {
	// 	rule := ruleTabel{}
	// 	var reg = regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)`)
	// 	for i, match1 := range reg.FindAllString(match, -1) {
	// 		arr := strings.Split(match1, ".")
	// 		if i == 0 {
	// 			rule.pTable.tableName = arr[0]
	// 			rule.pTable.fieldName = arr[1]
	// 		} else if i == 1 {
	// 			rule.lnkTable.tableName = arr[0]
	// 			rule.lnkTable.fieldName = arr[1]
	// 		}
	// 	}
	// 	rulse = append(rulse, rule)
	// }

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

func getOpr(str []byte) int {
	operators := map[string]int{
		"IS":  1,
		"OR":  2,
		"AND": 3,
	}
	var opre = regexp.MustCompile(`(?m)([A-Z]{2,5})`)
	opr := opre.Find(str)
	if v, ok := operators[string(opr)]; ok {
		return v
	}
	return 0
}
