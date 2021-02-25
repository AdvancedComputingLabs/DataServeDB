package db

import (
	"DataServeDB/commtypes"
	"DataServeDB/dbtable"
	"encoding/json"
	"fmt"
	"sort"

	// "DataServeDB/unstable_api/runtime"
	"net/http"
)

type Query struct {
	ItemLabel string
	ItemType  string
	ItemValue []byte // json Converted
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
	sort.SliceStable(query.Children, func(i, j int) bool { return query.Children[i].ItemType == "field" })
	// i := sort.Search(len(query.Children), func(i int) bool { return query.Children[i].ItemValue != nil })
	// // println(i)
	// if i < len(query.Children) {
	// 	println(query.Children[i].ItemLabel)
	// }
	/*********************************************/
	/* TO DO :- filter out the specific row if the userId mentioned,
	example {
		"$type": "application/json",
		"Users": {
	        "UserId": 1,
	        "*": {}
	    }
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
			query.Children[i].ItemType = "field"
		} else if tbl, e := t.GetTable(value.ItemLabel); e == nil {
			child, err := t.verifyQuery(query.Children[i], tbl)
			if err != nil {
				return Query{}, err
			}
			query.Children[i] = child
			query.Children[i].ItemType = "table"
		} else {
			return Query{}, fmt.Errorf("Query field '%s' does not exit in database", value.ItemLabel)
		}
	}
	return query, nil
}
func (t *DB) processQuery(query Query) (result []dbtable.TableRow, err error) {
	rows := []dbtable.TableRow{}
	table, err := t.GetTable(query.ItemLabel)
	if err != nil {
		// resultErr = err
		return
	}
	spec := getSpec(query)
	if spec != -1 {
		rows, err = table.GetTableRows(string(query.Children[spec].ItemValue))

	} else {
		rows, err = table.GetTableRows()

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
				refRows := []dbtable.TableRow{}
				if child.ItemLabel == "Properties" {
					table, er := t.GetTable("UserProperties")
					if er != nil {
						err = er
						return
					}
					refRows, err = table.GetRows()
					if err != nil {
						return
					}
				}
				r, err := t.processQuery(child)
				if err != nil {
					return nil, err
				}
				if refRows != nil {
					ref := []dbtable.TableRow{}
					for _, v := range refRows {
						if row["Id"] == v["Id"] {
							for _, rr := range r {
								if rr["SlNum"] == v["SlNum"] {
									ref = append(ref, rr)
								}
							}
						}
					}
					res[child.ItemLabel] = ref
				} else {

					res[child.ItemLabel] = r
				}
			}
		}
		result = append(result, res)
	}
	return
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
