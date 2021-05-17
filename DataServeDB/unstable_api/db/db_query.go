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
type QueryOprnd interface{}

const (
	OpNone = iota // starting with 0
	OpIS
	OpISNOT
	OpOR
	OpAND
	OpBETWEEN
	OpGT
	OpGTEQ
	opEQ
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
	relToPar  *relInfo
	relToChld *relInfo
}
type setRule struct {
	res      *RuleFieldInfo // result table name
	ref      *RuleFieldInfo // reference tble name
	relation string         // relation table name
	rule     Rule           // rules
}
type rowInfo struct {
	name string
	row  dbtable.TableRow
}

type Row map[string]interface{}

var Operators = map[string]QueryOp{
	"IS":      OpIS,
	"IS NOT":  OpISNOT,
	"OR":      OpOR,
	"AND":     OpAND,
	"BETWEEN": OpBETWEEN,
	"==":      opEQ,
	">":       OpGT,
	">=":      OpGTEQ,
	"<":       OpLT,
	"<=":      OpLTEQ,
	"!=":      OpNTEQ,
}

func (t *DB) TablesQueryGet(dbReqCtx *commtypes.DbReqContext, query []Query) (resultHttpStatus int, resultContent []byte, resultErr error) {
	// verify query one by one
	for i, qry := range query {
		// getting table by query item label for verify
		tbl, err := t.GetTable(qry.ItemLabel)
		if err != nil {
			resultErr = err
			return
		}
		qry, err = t.verifyQuery(qry, tbl)
		if err != nil {
			resultErr = err
			return
		}
		query[i] = qry
	}
	result, err := t.processRoot(query)
	if err != nil {
		resultHttpStatus = http.StatusNoContent
		return resultHttpStatus, nil, err
	}
	// convert the result to json format, and store as bytes
	resultContent, resultErr = json.Marshal(result)
	if resultErr != nil {
		resultHttpStatus = http.StatusNoContent
		return resultHttpStatus, nil, resultErr
	}
	resultHttpStatus = http.StatusOK
	return
}

// verifyQuery has the function of verifying the query and marks the field and table
func (t *DB) verifyQuery(query Query, tbl *dbtable.DbTable) (Query, error) {
	if tb, e := t.GetTable(query.ItemLabel); e == nil {
		query.ItemType = table
		for i, qry := range query.Children {
			child, err := t.verifyQuery(qry, tb)
			if err != nil {
				return Query{}, err
			}
			query.Children[i] = child
		}
	} else if tbl != nil {
		if _, found := tbl.TblMain.TableFieldsMetaData.IsField(query.ItemLabel); found {
			query.ItemType = field
		} else {
			return Query{}, fmt.Errorf("Query field '%s' does not exist in database", query.ItemLabel)
		}
	} else {
		return Query{}, fmt.Errorf("Query field '%s' does not exist in database", query.ItemLabel)
	}

	return query, nil
}

// processRoot process the array of query
// initial process of  query array, loops throug the queries and collects the results
func (t *DB) processRoot(queries []Query) (dbtable.TableRow, error) {
	result := dbtable.TableRow{}
	for _, query := range queries {
		if query.ItemType != table {
			return nil, fmt.Errorf("ivalid item type or orphen")
		}
		res, err := t.processQuery(query, rowInfo{})
		if err != nil {
			return nil, err
		}
		result[query.ItemLabel] = res
	}

	return result, nil
}

// process child is to loop the children and process each one to processQuery function and returns the result
func (t *DB) processChild(children []Query, parent rowInfo) (result dbtable.TableRow, err error) {
	res := dbtable.TableRow{}
	for _, child := range children {
		res[child.ItemLabel], err = t.processQuery(child, parent)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
func (t *DB) processQuery(query Query, parent rowInfo) (result interface{}, err error) {
	// if get the field item
	// else if get the table item, processing sub tables in the query
	if query.ItemType == field {
		if parent.row != nil {
			return parent.row[query.ItemLabel], nil
		}
		return nil, fmt.Errorf("ivalid item type or orphen")
	} else if query.ItemType == table {
		tabres := []dbtable.TableRow{}
		rows := []dbtable.TableRow{}

		// collecting the rules Info (JOIN and WHERE)
		joinInfo := getJoinInfo(query.Rules)
		joinInfoArr := setJoinRelation(joinInfo, parent.name, query.ItemLabel)
		whereInfo := getWhereInfo(query.Rules)

		tbl, err := t.GetTable(query.ItemLabel)
		if err != nil {
			return nil, err
		}
		spec := getSpec(query)
		// get the value if specific value mentioned for any children for filtering, then those roes oly will return on passing the index or value
		if spec != -1 {
			rows, err = tbl.GetTableRows(string(query.Children[spec].ItemValue))
			if err != nil {
				return nil, nil
			}
		} else {
			rows, err = tbl.GetTableRows()
			if err != nil {
				return nil, nil
			}
		}
		for _, row := range rows {
			res := dbtable.TableRow{}
			if query.Rules != nil {
				// IF checkAgainstJoinRelations(joinInfoArr, parentRowInfo, currentRowInfo) == FALSE: continue;
				if joinInfo != nil {
					if !t.checkIsChild(joinInfoArr, parent.row, row) {
						continue
					}
				}

				//	IF checkAgainstWhereClauses(whereInfoArr, parentRowInfo, currentRowInfo) == FALSE: continue;
				if whereInfo != nil {
					ok := whereClouse(whereInfo, parent, rowInfo{query.ItemLabel, row})
					if !ok {
						continue
					}
				}

			}
			// will process the chldren after passing JOIN and WHERE clouse,
			if query.Children != nil {
				res, err = t.processChild(query.Children, rowInfo{query.ItemLabel, row})
				if err != nil {
					return nil, err
				}
			} else {
				// get whole row as the result if there is no rules and children to process
				res = row
			}
			tabres = append(tabres, res)
		}
		return tabres, nil
	}
	return
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

func (t *DB) getSingleRowBytableName(tableName string, fieldName string, value interface{}) (row dbtable.TableRow, err error) {
	tbl, err := t.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	row, err = tbl.GetTableSingleRow(fieldName, value)
	if err != nil {
		return nil, err
	}
	return
}

// check join relation of parent and child rows
func (t *DB) checkIsChild(joinInfo relation, parentRow, currentRow dbtable.TableRow) bool {
	if joinInfo.relToChld != nil {
		relRow, err := t.getSingleRowBytableName(joinInfo.relToChld.relation.TableName, joinInfo.relToChld.relation.FieldName, currentRow[joinInfo.relToChld.parOrChld.FieldName])
		if err != nil {
			return false
		} else {
			if relRow[joinInfo.relToPar.relation.FieldName] == parentRow[joinInfo.relToPar.parOrChld.FieldName] && relRow[joinInfo.relToChld.relation.FieldName] == currentRow[joinInfo.relToChld.parOrChld.FieldName] {
				return true
			}
		}
	} else {
		if parentRow[joinInfo.relToPar.parOrChld.FieldName] == currentRow[joinInfo.relToPar.relation.FieldName] {
			return true
		}
	}
	return false
}
