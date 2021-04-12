package db

import (
	"DataServeDB/commtypes"
	"DataServeDB/dbtable"
	"encoding/json"
	"fmt"

	"net/http"
)

type QueryOp int
type ItemType int

const (
	OpNone = iota // starting with 0
	OpIS
	OpOR
	OpAND
	OpGT
	OpGTEQ
	OpLT
	OpLTEQ
)

const (
	field = iota
	table
	// rule
)

// type RuleInfo struct {
// 	TableName string
// 	FieldName string
// 	Operation QueryOp
// 	Next      *RuleInfo
// }

type RuleFieldInfo struct {
	TableName string
	FieldName string
}

type RuleFeild struct {
	LeftRule     *RuleFieldInfo
	RightRule    *RuleFieldInfo
	LeftOperand  string
	RightOperand string

	Operator QueryOp
}

type Rule []interface{}
type Rules struct {
	Label string
	Rule  Rule
}

type Query struct {
	ItemLabel string
	ItemType  ItemType
	ItemValue string // json Converted
	Rules     []Rules
	Children  []Query
}
type setRule struct {
	res      *RuleFieldInfo // result table name
	ref      *RuleFieldInfo // reference tble name
	relation string         // relation table name
	rule     Rules          // rules
}

type Row map[string]interface{}

var Operators = map[string]QueryOp{
	"IS":  OpIS,
	"OR":  OpOR,
	"AND": OpAND,
	">":   OpGT,
	">=":  OpGTEQ,
	"<":   OpLT,
	"<=":  OpLTEQ,
}

func (t *DB) TablesQueryGet(dbReqCtx *commtypes.DbReqContext, query Query) (resultHttpStatus int, resultContent []byte, resultErr error) {
	tbl, err := t.GetTable(query.ItemLabel)
	if err != nil {
		resultErr = err
		return
	}
	query, err = t.verifyQuery(query, tbl)
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

func (t *DB) verifyQuery(query Query, tbl *dbtable.DbTable) (Query, error) {
	for i, value := range query.Children {
		if _, found := tbl.TblMain.TableFieldsMetaData.IsField(value.ItemLabel); found {
			/* TODO make itemtype as macro */
			query.Children[i].ItemType = field
		} else if tbl, e := t.GetTable(value.ItemLabel); e == nil {
			child, err := t.verifyQuery(query.Children[i], tbl)
			if err != nil {
				return Query{}, err
			}
			query.Children[i] = child
			query.Children[i].ItemType = table
		} else if value.ItemLabel == "$WHERE" {
			// query.ItemType = Rule
			// fmt.Println(value.Rules)
			// value.Rules
		} else {
			return Query{}, fmt.Errorf("Query field '%s' does not exist in database", value.ItemLabel)
		}
	}
	return query, nil
}
func (t *DB) processQuery(query Query) (result []dbtable.TableRow, err error) {
	rows := []dbtable.TableRow{}
	tbl, err := t.GetTable(query.ItemLabel)
	if err != nil {
		return
	}
	spec := getSpec(query)
	// get the value if specific item mentioned
	if spec != -1 {
		println(" spec")
		rows, err = tbl.GetTableRows(string(query.Children[spec].ItemValue))
		if err != nil {
			return
		}
	} else {
		println("no spec")
		rows, err = tbl.GetTableRows()
		if err != nil {
			return
		}
	}

	// if there is no sub bracnches
	if query.Children == nil {
		result = rows
		return
	}
	for _, row := range rows {
		res := dbtable.TableRow{}

		for _, child := range query.Children {
			ArrRows := map[string][]dbtable.TableRow{}
			if child.ItemType == field {
				res[child.ItemLabel] = row[child.ItemLabel]
			} else if child.ItemType == table {
				temp := []dbtable.TableRow{}
				tbres := []dbtable.TableRow{}
				if _, ok := ArrRows[query.ItemLabel]; !ok {
					ArrRows[query.ItemLabel] = append(temp, row)
				}
				for _, rules := range child.Rules {
					if rules.Label == "$JOIN" {
						set := setRules(rules, query.ItemLabel, child.ItemLabel)
						if err = t.getTablesFromRules(ArrRows, rules.Rule); err != nil {
							return nil, err
						}
						tbres, err = processJoin(ArrRows, set)
						if err != nil {
							return nil, err
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
		if v.ItemType == field {
			if v.ItemValue != "" {
				return i
			}
		}
	}
	return -1
}

func (t *DB) getRowsBytableName(tableName string) (rows []dbtable.TableRow, err error) {
	tbl, err := t.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	rows, err = tbl.GetTableRows()
	if err != nil {
		return nil, err
	}
	return
}

// getTablesFromRules get the table rows from the rule structure
func (t *DB) getTablesFromRules(ArrRows map[string][]dbtable.TableRow, ruleAst Rule) error {
	rul := []*RuleFieldInfo{}
	for _, rule := range ruleAst {
		if rl, ok := rule.(RuleFeild); ok {
			rul := append(rul, rl.LeftRule, rl.RightRule)
			for _, v := range rul {
				if v != nil {
					if _, ok := ArrRows[v.TableName]; !ok {
						row, err := t.getRowsBytableName(v.TableName)
						if err != nil {
							return err
						}
						ArrRows[v.TableName] = row
					}
				}
			}
		}
	}
	return nil
}
func setRules(rule Rules, ref, res string) (set setRule) {
	// set := setRule{}
	// ruleAst := rules.Rule
	// rul := []*RuleFieldInfo{}
	set.rule = rule
	set.res = &RuleFieldInfo{res, ""}
	set.ref = &RuleFieldInfo{ref, ""}
	for _, rule := range rule.Rule {
		if rl, ok := rule.(RuleFeild); ok {
			if set.ref.FieldName == "" {
				if set.ref.TableName == rl.LeftRule.TableName {
					set.ref = rl.LeftRule
				} else if set.ref.TableName == rl.RightRule.TableName {
					set.ref = rl.RightRule
				}
			}
			if set.res.FieldName == "" {
				if set.res.TableName == rl.LeftRule.TableName {
					set.res = rl.LeftRule
				} else if set.ref.TableName == rl.RightRule.TableName {
					set.res = rl.RightRule
				}
			}
			if set.relation == "" {
				if rl.LeftRule.TableName != set.ref.TableName && rl.LeftRule.TableName != set.res.TableName {
					set.relation = rl.LeftRule.TableName
				} else if rl.RightRule.TableName != set.ref.TableName && rl.RightRule.TableName != set.res.TableName {
					set.relation = rl.RightRule.TableName
				}
			}
		}
	}
	return
}

func processJoin(ArrRows map[string][]dbtable.TableRow, set setRule) ([]dbtable.TableRow, error) {
	res := []dbtable.TableRow{}
	if set.relation != "" {
		if set.ref != nil {
			res = filteRows(ArrRows[set.ref.TableName], ArrRows[set.relation], set.ref.FieldName)
			fmt.Println("rela, ", res)
		} else {
			return nil, fmt.Errorf("no reference rule found")
		}
		if set.res != nil {
			return filteRows(res, ArrRows[set.res.TableName], set.res.FieldName), nil
		} else {
			return nil, fmt.Errorf("no reference rule found")
		}
	} else {
		if set.ref != nil && set.res != nil {
			return filteRows(ArrRows[set.ref.TableName], ArrRows[set.res.TableName], set.res.FieldName), nil
		}
		return nil, fmt.Errorf("no reference rule found")
	}
}
func filteRows(ref, res []dbtable.TableRow, field string) (result []dbtable.TableRow) {
	for _, refRow := range ref {
		for _, relRow := range res {
			if refRow[field] == relRow[field] {
				result = append(result, relRow)
			}
		}
	}
	return
}

// Process $WHERE

// func (t *DB) processRules(rules Rule, ArrRows map[string][]dbtable.TableRow) ([]dbtable.TableRow, error) {
// 	tbres := []dbtable.TableRow{}

// 	var crntOpr QueryOp = OpNone
// 	temp := []dbtable.TableRow{}
// 	// prev := []dbtable.TableRow{}
// 	for i, rule := range rules {
// 		if rl, ok := rule.(RuleFeild); ok {
// 			// var left, right interface{}
// 			if rl.LeftRule != nil {
// 				for _, lv := range ArrRows[rl.LeftRule.TableName] {
// 					if rl.RightRule != nil {
// 						for _, rv := range ArrRows[rl.RightRule.TableName] {
// 							if rl.Operator == OpIS {
// 								if rv[rl.RightRule.FieldName] == lv[rl.LeftRule.FieldName] {
// 									temp = append(temp, lv)
// 								}
// 							}
// 						}
// 					} else if rl.RightOperand != "" {
// 						if rl.Operator == OpIS {
// 							if lv[rl.LeftRule.FieldName] == rl.RightOperand {
// 								temp = append(temp, lv)
// 							}
// 						}
// 					}
// 				}
// 			} else if rl.LeftOperand != "" {
// 				if rl.RightRule != nil {
// 					for _, rv := range ArrRows[rl.RightRule.TableName] {
// 						if rl.Operator == OpIS {
// 							if rv[rl.RightRule.FieldName] == rl.LeftOperand {
// 								temp = append(temp, rv)
// 							}
// 						}
// 					}
// 				} else if rl.RightOperand != "" {
// 					if rl.LeftOperand == rl.RightOperand {
// 						// error handling
// 					}
// 				}
// 			}
// 		} else if rl, ok := rule.([]interface{}); ok {
// 			temp, err := t.processRules(rl, ArrRows)
// 			if err != nil {
// 				return nil, err
// 			}
// 			_ = temp
// 		}

// 		if opr, ok := rule.(QueryOp); ok {
// 			crntOpr = opr
// 		} else if i != 0 {
// 			// todo operations on stackwise
// 			// process temp with prev
// 			if crntOpr == OpIS {
// 				//
// 			} else if crntOpr == OpAND {

// 			} else if crntOpr == OpOR {

// 			} else if crntOpr == OpGT {

// 			} else if crntOpr == OpGTEQ {

// 			} else if crntOpr == OpLTEQ {

// 			} else if crntOpr == OpLT {

// 			}

// 		} else {

// 			// prev = temp
// 		}
// 	}

// 	return tbres, nil
// }

// func isFunc(left, right interface{}) {

// }

// func doOp(left, right interface{}, opr QueryOp) {
// 	if opr == OpIS {
// 		isFunc(left, right)
// 	}

// }
