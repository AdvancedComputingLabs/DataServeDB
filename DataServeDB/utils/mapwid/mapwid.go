package mapwid

import "errors"

/*
	Description: Map with Ids allows changable name with an id that stays same.

	Important Points:
		1) Id range is from 0 to 1_000_000. One million should be enough for its use cases.
		1) Ids and names must be unique.
		2) Operations are not thread safe, it has to be handled prior to the operations.
		3) Name mapping is not case insensitve, it has to be handled prior to the operations.
		4) Id does not change to keep the operations simple. Remove is allowed but use it with care to avoid unexpected bugs.

	Dev Issues:
		1) It is using map, which avoids some problems compared to list but it is more expensive.
		2) There is no get by id because it would need id to name mapping. It maybe added later if needed.
*/

//NOTE: Might run into struct export issues, hence, kept the field public. But don't use them outside of the package directly.
//TODO: Make them private and test.

type MapWithId struct {
	IdMap       map[int]interface{}
	NameToIdMap map[string]int
	LastId      int
}

func New() *MapWithId {
	return &MapWithId{
		IdMap:       make(map[int]interface{}),
		NameToIdMap: make(map[string]int),
		LastId:      -1,
	}
}

func (t *MapWithId) AddUnsync(id int, nameCaseSen string, object interface{}) error {

	//TODO: error message format needs standardization.

	if id < 0 {
		return errors.New("id cannot be negative")
	}

	if id > 1_000_000 {
		return errors.New("id cannot be greater than 1,000,000")
	}

	if _, exists := t.IdMap[id]; exists {
		return errors.New("id already exists") //TODO: make it more user friendly
	}

	if _, exists := t.NameToIdMap[nameCaseSen]; exists {
		return errors.New("name already exists") //TODO: make it more user friendly
	}

	t.IdMap[id] = object
	t.NameToIdMap[nameCaseSen] = id

	if t.LastId < id {
		t.LastId = id
	}

	return nil
}

// GetByName get by name
func (t *MapWithId) GetByNameUnsync(nameCaseSen string) (int, interface{}, error) {
	var id int
	var exists bool
	var object interface{}

	// for key := range t.NameToIdMap {
	// 	println("Key:", key)
	// }

	if id, exists = t.NameToIdMap[nameCaseSen]; !exists {
		return -1, nil, errors.New("name does not exist") //TODO: make it more user friendly
	}

	if object, exists = t.IdMap[id]; !exists {
		return -1, nil, errors.New("id does not exist") //TODO: make it more user friendly
	}

	return id, object, nil
}

func (t *MapWithId) GetIdUnsync(nameCaseSen string) (int, error) {
	var id int
	var exists bool

	if id, exists = t.NameToIdMap[nameCaseSen]; !exists {
		return -1, errors.New("name does not exist") //TODO: make it more user friendly
	}

	return id, nil
}

func (t *MapWithId) GetItemsUnsync() map[int]interface{} {
	return t.IdMap
}

func (t *MapWithId) RemoveByNameUnsync(nameCaseSen string) (int, interface{}, error) {
	var id int
	var exists bool
	var object interface{}

	if id, exists = t.NameToIdMap[nameCaseSen]; !exists {
		return -1, nil, errors.New("name does not exist") //TODO: make it more user friendly
	}

	if object, exists = t.IdMap[id]; !exists {
		return -1, nil, errors.New("id does not exist") //TODO: make it more user friendly
	}

	//NOTE: deletion from IdMap is first for a reason.
	//If there goes something wrong while deleting from IdMap or NameToIdMap easier to correct by name,
	//but if NameToIdMap is delete first and IdMap id is not then it is more difficult find and correct.
	delete(t.IdMap, id)
	delete(t.NameToIdMap, nameCaseSen)

	return id, object, nil
}

func (t *MapWithId) UpdateUnsync(nameCurrentCaseSen, nameNewCaseSen string) error {
	//COMMENT: do caller need to send id for checking?
	// It can be checked before calling this through through GetId first.

	var id int
	var exists bool

	if id, exists = t.NameToIdMap[nameCurrentCaseSen]; !exists {
		return errors.New("name does not exist") //TODO: make it more user friendly
	}

	if _, exists = t.IdMap[id]; !exists {
		return errors.New("id does not exist") //TODO: make it more user friendly
	}

	t.NameToIdMap[nameNewCaseSen] = id

	delete(t.NameToIdMap, nameCurrentCaseSen)

	return nil
}

// HasName finds wether the new entry name has alredy exist or not
func (t *MapWithId) HasNameUnsync(nameCaseSen string) bool {
	_, exists := t.NameToIdMap[nameCaseSen]
	return exists //TODO: make it more user friendly
}

func (t *MapWithId) GetLastIdUnsync() int {
	return  t.LastId
}
