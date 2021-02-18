package db

import (
	"DataServeDB/commtypes"
	"DataServeDB/dbtable"
	"encoding/json"
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
	result := []dbtable.TableRow{}
	table, err := t.GetTable(query.ItemLabel)
	if err != nil {
		resultErr = err
		return
	}
	for i, value := range query.Children {
		if _, found := table.TblMain.TableFieldsMetaData.IsField(value.ItemLabel); found {
			/* TODO make itemtype as macro */
			query.Children[i].ItemType = "field"
		} else {
			query.Children[i].ItemType = "table"
		}
	}
	// sort.SliceStable(query.Children, func(i, j int) bool { return query.Children[i].ItemType == "field" })
	i := sort.Search(len(query.Children), func(i int) bool { return query.Children[i].ItemValue != nil })
	// println(i)
	if i < len(query.Children) {
		println(query.Children[i].ItemLabel)
	}
	rows, resultErr := table.GetTableRows()
	/*********************************************/
	/* TO DO :- filter out the specific row if the userId mentioned,
	example {
		"$type": "application/json",
		"Users": {
	        "UserId": 1,
	        "*": {}
	    }
	}*/
	for _, row := range rows {
		res := dbtable.TableRow{}
		for _, child := range query.Children {
			if child.ItemType == "field" {
				res[child.ItemLabel] = row[child.ItemLabel]
			} else if child.ItemType == "table" {
				tbl, err := t.GetTable(child.ItemLabel)
				if err != nil {
					resultErr = err
					return
				}
				/* TO DO :- Filter Properties for user  */
				props, err := tbl.GetRows()
				if err != nil {
					resultErr = err
					return
				}
				res[child.ItemLabel] = props
			}
		}
		result = append(result, res)
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		resultErr = err
		//TODO: make error result more user friendly.
		return
	}
	return http.StatusOK, []byte(jsonBytes), resultErr

	//return 0, nil, nil
}

// func checkField(tableMain *, field string) {

// }
