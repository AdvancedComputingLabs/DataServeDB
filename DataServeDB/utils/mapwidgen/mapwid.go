package mapwidgen

import (
	"DataServeDB/constraints"
	"errors"
)

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

type MapWithId[T constraints.HasId] struct {
	IdMap       map[int]T
	NameToIdMap map[string]int
	LastId      int
}

func New[T constraints.HasId]() *MapWithId[T] {
	return &MapWithId[T]{
		IdMap:       make(map[int]T),
		NameToIdMap: make(map[string]int),
		LastId:      -1,
	}
}

func (t *MapWithId[T]) AddUnsync(id int, nameCaseSen string, object T) error {

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

func (t *MapWithId[T]) GetByNameUnsync(nameCaseSen string) (int, T, error) {
	var id int
	var exists bool
	var object T

	// for key := range t.NameToIdMap {
	// 	println("Key:", key)
	// }

	if id, exists = t.NameToIdMap[nameCaseSen]; !exists {
		return -1, object, errors.New("name does not exist") //TODO: make it more user friendly
	}

	if object, exists = t.IdMap[id]; !exists {
		return -1, object, errors.New("id does not exist") //TODO: make it more user friendly
	}

	return id, object, nil
}

func (t *MapWithId[T]) GetIdUnsync(nameCaseSen string) (int, error) {
	var id int
	var exists bool

	if id, exists = t.NameToIdMap[nameCaseSen]; !exists {
		return -1, errors.New("name does not exist") //TODO: make it more user friendly
	}

	return id, nil
}

func (t *MapWithId[T]) GetItemsUnsync() map[int]T {
	return t.IdMap
}

func (t *MapWithId[T]) RemoveByNameUnsync(nameCaseSen string) (int, T, error) {
	var id int
	var exists bool
	var object T

	if id, exists = t.NameToIdMap[nameCaseSen]; !exists {
		return -1, object, errors.New("name does not exist") //TODO: make it more user friendly
	}

	if object, exists = t.IdMap[id]; !exists {
		return -1, object, errors.New("id does not exist") //TODO: make it more user friendly
	}

	//NOTE: deletion from IdMap is first for a reason.
	//If there goes something wrong while deleting from IdMap or NameToIdMap easier to correct by name,
	//but if NameToIdMap is delete first and IdMap id is not then it is more difficult find and correct.
	delete(t.IdMap, id)
	delete(t.NameToIdMap, nameCaseSen)

	return id, object, nil
}

func (t *MapWithId[T]) UpdateUnsync(nameCurrentCaseSen, nameNewCaseSen string) error {
	//COMMENT: do caller need to send id for checking?
	// It can be checked before calling this through GetId first.

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

// HasNameUnsync finds wether the new entry name has alredy exist or not
func (t *MapWithId[T]) HasNameUnsync(nameCaseSen string) bool {
	_, exists := t.NameToIdMap[nameCaseSen]
	return exists //TODO: make it more user friendly
}

func (t *MapWithId[T]) GetLastIdUnsync() int {
	return t.LastId
}
