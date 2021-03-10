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
	Rules     []ruleTabel
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
			rules := getRules(string(value.ItemValue))
			query.Rules = rules
			return query, nil
		} else {
			return Query{}, fmt.Errorf("Query field '%s' does not exit in database", value.ItemLabel)
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
				temp := []dbtable.TableRow{}
				tbres := []dbtable.TableRow{}
				//     "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
				if child.Rules[0].pTable.tableName != query.ItemLabel {
					return nil, fmt.Errorf("primary table not fount in rule")
				}
				if _, ok := ArrRows[query.ItemLabel]; !ok {
					ArrRows[query.ItemLabel] = rows
				}
				for i, rule := range child.Rules {
					if _, ok := ArrRows[rule.pTable.tableName]; !ok {
						rows1, err := t.getRowsBytableName(rule.pTable.tableName)
						if err != nil {
							return nil, err
						}
						ArrRows[rule.pTable.tableName] = rows1
					}
					if _, ok := ArrRows[rule.lnkTable.tableName]; !ok {
						rows1, err := t.getRowsBytableName(rule.lnkTable.tableName)
						if err != nil {
							return nil, err
						}
						ArrRows[rule.lnkTable.tableName] = rows1
					}
					if i == 0 {
						for _, rv := range ArrRows[rule.lnkTable.tableName] {
							if row[rule.pTable.fieldName] == rv[rule.pTable.fieldName] {
								temp = append(temp, rv)
							}
						}
					} else if i == 1 {
						for _, prv := range ArrRows[rule.pTable.tableName] {
							for _, rv := range temp {
								if rv[rule.lnkTable.fieldName] == prv[rule.lnkTable.fieldName] {
									tbres = append(tbres, prv)
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

func getRules(str string) (rulse []ruleTabel) {
	// tblRule := []ruleTabel{}
	// "Properties": [
	//   {
	//     "$WHERE": "Users.Id IS UserProperties.Id AND Properties.SlNum IS UserProperties.SlNum"
	//   }
	// ]

	var re = regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)( IS )([A-z]*[.][A-z]*)`)

	for _, match := range re.FindAllString(str, -1) {
		rule := ruleTabel{}
		var reg = regexp.MustCompile(`(?m)([A-z]*[.][A-z]*)`)
		for i, match1 := range reg.FindAllString(match, -1) {
			arr := strings.Split(match1, ".")
			if i == 0 {
				rule.pTable.tableName = arr[0]
				rule.pTable.fieldName = arr[1]
			} else if i == 1 {
				rule.lnkTable.tableName = arr[0]
				rule.lnkTable.fieldName = arr[1]
			}

		}
		rulse = append(rulse, rule)
	}

	return
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
