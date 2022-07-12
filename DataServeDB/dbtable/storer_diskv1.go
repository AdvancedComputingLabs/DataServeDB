package dbtable

import (
	storage "DataServeDB/dbsystem/dbstorage"
	"bytes"
	"encoding/gob"
	"errors"
	"log"
)

type DiskStoreV1 struct {
	fileName string
	path     string
}

func NewDiskStoreV1(fileName, path string) (*DiskStoreV1, error) {
	return &DiskStoreV1{fileName: fileName, path: path}, nil
}

func (d DiskStoreV1) Implemented(feature TableStorerFeatures) bool {

	switch feature {
	case TableStorerFeatureInsert:
		return true
	}

	return false
}

func (d DiskStoreV1) Delete(indexName string, key string) (int, error) {
	return -1, errors.New("NotImplemented")
}

func (d DiskStoreV1) Get(indexName string, key string) (int, any, error) {
	return -1, nil, errors.New("NotImplemented")
}

func (d DiskStoreV1) Insert(rowWithProps TableRowWithFieldProperties, data any) (int, error) {

	//TODO: refector this into own function with binary or json option.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data) //NOTE: currently this saves all data and index structure to the disk.
	if err != nil {
		//TODO: should return error
		println("error ")
		log.Fatal("encode error:", err)
	}

	//TODO: disk error handling is needed.
	storage.SaveToDisk(buf.Bytes(), d.path)

	//TODO: 1 is correct?
	return 1, nil
}

func (d DiskStoreV1) Update(rowWithProps TableRowWithFieldProperties, data any) (int, error) {
	return -1, errors.New("NotImplemented")
}
