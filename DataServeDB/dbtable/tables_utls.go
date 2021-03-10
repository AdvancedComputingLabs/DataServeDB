package dbtable

import (
	storage "DataServeDB/dbsystem/dbstorage"
	"DataServeDB/paths"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"strings"
)

func parseKeyValue(resPath string) (key string, value string, err error) {
	pos := strings.LastIndex(resPath, "/") + 1

	if pos == 0 {
		err = errors.New("key path is in wrong format")
		return
	}

	if pos >= len(resPath) {
		err = errors.New("key or value is not provided")
		return
	}

	splitted := strings.SplitN(resPath[pos:], ":", 2)

	if len(splitted) == 1 {
		value = splitted[0]
	} else {
		key = splitted[0]
		value = splitted[1]
	}

	return
}

func saveToDiskUtil(t *DbTable) error {

	//TODO: path for table needs its own function?
	fileName := fmt.Sprintf("table_%d.dat", t.TblMain.TableId)
	path := paths.Combine(t.createTableStructure._dbPtr.DbPath(), tablesDataPathRelative, fileName)

	//TODO: refector this into own function with binary or json option?
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t.TblData)
	if err != nil {
		//TODO: better handling needed
		log.Fatal("encode error:", err)
	}

	//TODO: disk error handling is needed.
	storage.SaveToDisk(buf.Bytes(), path)

	return nil
}

// deleting the row of table by slicing method,  and returning the value which want to updated on map
func deleteRowUnordered(Rows []map[int]interface{}, rowNum int64) ([]map[int]interface{}, error) {
	if int(rowNum) >= len(Rows) {
		return nil, fmt.Errorf("there is a miss match in row index")
	}
	lastIndxVal := Rows[len(Rows)-1]
	Rows[rowNum] = lastIndxVal
	Rows[len(Rows)-1] = nil
	Rows = Rows[:len(Rows)-1]

	return Rows, nil
}
