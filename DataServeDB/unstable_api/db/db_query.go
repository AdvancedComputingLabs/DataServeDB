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
	OpISNOT
	OpOR
	OpAND
	OpBETWEEN
	OpGT
	OpGTEQ
	OpLT
	OpLTEQ
	OpNTEQ
)

const (
	field = iota
	table
)

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
type relInfo struct {
	relation  *RuleFieldInfo
	parOrChld *RuleFieldInfo
}
type relation struct {
	relToPar  relInfo
	relToChld relInfo
}
type setRule struct {
	res      *RuleFieldInfo // result table name
	ref      *RuleFieldInfo // reference tble name
	relation string         // relation table name
	rule     Rule           // rules
}

type Row map[string]interface{}

var Operators = map[string]QueryOp{
	"IS":      OpIS,
	"IS NOT":  OpISNOT,
	"OR":      OpOR,
	"AND":     OpAND,
	"BETWEEN": OpBETWEEN,
	">":       OpGT,
	">=":      OpGTEQ,
	"<":       OpLT,
	"<=":      OpLTEQ,
	"!=":      OpNTEQ,
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

// verifyQuery has the function of veryfing the query and marks the field and table
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
			if child.ItemType == field {
				res[child.ItemLabel] = row[child.ItemLabel]
			} else if child.ItemType == table {
				tbres := []dbtable.TableRow{}
				joinInfo := getJoinInfo(child.Rules)
				joinInfoArr := setRelation(joinInfo, query.ItemLabel, child.ItemLabel)
				parentRowInfo := row
				tabRows, err := t.getRowsBytableName(child.ItemLabel)
				if err != nil {
					return nil, err
				}
				for _, currentRow := range tabRows {
					// IF checkAgainstJoinRelations(joinInfoArr, parentRowInfo, currentRowInfo) == FALSE: continue;
					if !t.checkAgainstJoinRelations(joinInfoArr, parentRowInfo, currentRow) {
						continue
					}
					//	IF checkAgainstWhereClauses(whereInfoArr, parentRowInfo, currentRowInfo) == FALSE: continue;
					// make as function (Rules, ArrRows)
					tbres = append(tbres, currentRow)
				}
				res[child.ItemLabel] = tbres
			}
		}
		result = append(result, res)
	}
	return result, nil
}

// this function theck the query has any specific values to any field
// it returns the index of field if has any other wise return -1
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

func (t *DB) checkAgainstJoinRelations(joinInfo relation, parentRow, currentRow dbtable.TableRow) bool {
	rel, err := t.getRowsBytableName(joinInfo.relToPar.relation.TableName)
	if err != nil {
		return false
	}
	for _, relRow := range rel {
		if relRow[joinInfo.relToPar.relation.FieldName] == parentRow[joinInfo.relToPar.parOrChld.FieldName] && relRow[joinInfo.relToChld.relation.FieldName] == currentRow[joinInfo.relToChld.parOrChld.FieldName] {
			return true
		}
	}
	return false
}
func getJoinInfo(rules []Rules) Rule {
	for _, rule := range rules {
		if rule.Label == "$JOIN" {
			return rule.Rule
		}
	}
	return nil
}
func getWhereInfo(rules []Rules) Rule {
	for _, rule := range rules {
		if rule.Label == "$WHERE" {
			return rule.Rule
		}
	}
	return nil
}

func setRelation(rule Rule, par, chld string) (rel relation) {
	for _, rule := range rule {
		if rf, ok := rule.(RuleFeild); ok {
			// set parent relation
			if par == rf.LeftRule.TableName {
				rel.relToPar.parOrChld = rf.LeftRule
				rel.relToPar.relation = rf.RightRule
			} else if par == rf.RightRule.TableName {
				rel.relToPar.parOrChld = rf.RightRule
				rel.relToPar.relation = rf.LeftRule
			}
			// set child relation
			if chld == rf.LeftRule.TableName {
				rel.relToChld.parOrChld = rf.LeftRule
				rel.relToChld.relation = rf.RightRule
			} else if chld == rf.RightRule.TableName {
				rel.relToChld.parOrChld = rf.RightRule
				rel.relToChld.relation = rf.LeftRule
			}
		}
	}
	return
}

func (t *DB) checkAgainstWhereClauses(joinInfo relation, parentRow, currentRow dbtable.TableRow) bool {

	return false
}

// setRules sets the ruels set inorder to process query clouses, it find the relation table, w.r.t ref and res
// func setRules(rule Rule, ref, res string) (set setRule) {
// 	set.rule = rule
// 	set.res = &RuleFieldInfo{res, ""}
// 	set.ref = &RuleFieldInfo{ref, ""}
// 	for _, rule := range rule {
// 		if rl, ok := rule.(RuleFeild); ok {
// 			if set.ref.FieldName == "" {
// 				if set.ref.TableName == rl.LeftRule.TableName {
// 					set.ref = rl.LeftRule
// 				} else if set.ref.TableName == rl.RightRule.TableName {
// 					set.ref = rl.RightRule
// 				}
// 			}
// 			if set.res.FieldName == "" {
// 				if set.res.TableName == rl.LeftRule.TableName {
// 					set.res = rl.LeftRule
// 				} else if set.ref.TableName == rl.RightRule.TableName {
// 					set.res = rl.RightRule
// 				}
// 			}
// 			if set.relation == "" {
// 				if rl.LeftRule.TableName != set.ref.TableName && rl.LeftRule.TableName != set.res.TableName {
// 					set.relation = rl.LeftRule.TableName
// 				} else if rl.RightRule.TableName != set.ref.TableName && rl.RightRule.TableName != set.res.TableName {
// 					set.relation = rl.RightRule.TableName
// 				}
// 			}
// 		}
// 	}
// 	return
// }

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
