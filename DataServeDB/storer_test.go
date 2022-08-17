package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"DataServeDB/unstable_api/runtime"
)

func TestRollback(t *testing.T) {
	db, e := runtime.GetDb("re_db")
	if e != nil {
		t.Errorf("%v\n", e)
		return
	}
	tbl, e := db.GetTable("Tbl01")
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Println("Rows --> ", tbl.TblData.Rows, "\nEnd Rows")
	//tableLength := len(tbl.TblData.Rows)

	items := []string{"Batman", "Superman", "Flash", "wonder women", "Shazaam"}
	length := tbl.GetLength()

	for i, item := range items {
		row01 := row{
			Id:       i + length,
			UserName: item,
		}

		row01Json, err := json.Marshal(row01)
		if err != nil {
			t.Error("error converting")
		} else {
			if e := tbl.InsertRowJSON(string(row01Json)); e == nil {
				// if tableLength != len(tbl.TblData.Rows) {
				// 	fmt.Println("The RollBack Test Successful")
				// }

				// TODO :- Should make an error on storage function to get succesfull test case
				if row, _ := tbl.GetRowByPrimaryKey(row01.Id); row == nil {
					fmt.Println("The RollBack Test Successful")
				} else {
					//fmt.Println(row01.Id)
					t.Errorf("The roollBack Test has failed!!!")
				}

			} else {
				t.Errorf("%v\n", e)
			}
		}
	}

	fmt.Println("Final Rows --> \n", tbl.TblData.Rows)
}
