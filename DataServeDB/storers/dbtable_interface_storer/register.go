package dbtable_interface_storer

import "errors"

var registeredStoresBasic = map[string]NewStoreBasic{}

// GetStoreBasic If not found returns nil
func GetStoreBasic(key string) NewStoreBasic {
	if storeCreator, ok := registeredStoresBasic[key]; ok {
		return storeCreator
	}
	return nil
}

func RegisterStoreBasic(key string, storeCreator NewStoreBasic) error {
	if _, exists := registeredStoresBasic[key]; exists {
		//can only have one key per storeCreator.
		return errors.New("'storeCreator' key already exits")
	}

	registeredStoresBasic[key] = storeCreator

	return nil
}
