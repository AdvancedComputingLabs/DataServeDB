package db

import (
	"DataServeDB/commtypes"
	"DataServeDB/dbtable"
	"encoding/json"
	"fmt"
	"regexp"

	"net/http"
)

type QueryOp int
type ItemType int

const (
	OpNone = iota // starting with 0
	OpIS
	OpOR
	OpAND
)

const (
	Field = iota
	Table
	Rule
)

type RuleInfo struct {
	TableName string
	FieldName string
	Operation QueryOp
	Next      *RuleInfo
}

type Query struct {
	ItemLabel string
	ItemType  ItemType
	ItemValue []byte // json Converted
	Rules     *RuleInfo
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
			// query.ItemType = Rule
			fmt.Println(value.Rules)
			// value.Rules
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
				for rule := child.Children[0].Rules; rule != nil; rule = rule.Next {
					if _, ok := ArrRows[rule.TableName]; !ok {
						rows1, err := t.getRowsBytableName(rule.TableName)
						if err != nil {
							return nil, err
						}
						ArrRows[rule.TableName] = rows1
					}
				}
				for rule := child.Children[0].Rules; rule != nil; rule = rule.Next {
					if rule.Operation == OpAND {
						continue
					}
					if rule.TableName == query.ItemLabel {
						tempRow = row
						if rule.Operation == OpIS {
							for _, rv := range ArrRows[rule.Next.TableName] {
								if tempRow[rule.FieldName] == rv[rule.FieldName] {
									temp = append(temp, rv)
								}
							}
						}
					}
					if rule.TableName == child.ItemLabel {
						if rule.Operation == OpIS {
							for _, prv := range ArrRows[rule.TableName] {
								for _, rv := range temp {
									if rv[rule.FieldName] == prv[rule.FieldName] {
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
