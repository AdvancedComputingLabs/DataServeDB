package dbtable

import "errors"

type MemStoreV1 struct {
}

func NewMemStoreV1() (*MemStoreV1, error) {
	return &MemStoreV1{}, nil
}

func (d MemStoreV1) Implemented(feature TableStorerFeatures) bool {

	switch feature {
	case TableStorerFeatureInsert:
		return true
	}

	return false
}

func (d MemStoreV1) Delete(indexName string, key string) (int, error) {
	return -1, errors.New("NotImplemented")
}

func (d MemStoreV1) Get(indexName string, key string) (int, any, error) {
	return -1, nil, errors.New("NotImplemented")
}

func (d MemStoreV1) Insert(rowWithProps TableRowWithFieldProperties, data any) (int, error) {
	return -1, errors.New("UndefiedErrorOccurred")
}

func (d MemStoreV1) Update(rowWithProps TableRowWithFieldProperties, data any) (int, error) {
	return -1, errors.New("NotImplemented")
}
