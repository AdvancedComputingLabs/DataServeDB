package constants

type DbResTypes int

type RestMethods int

const (
	DbResTypeNone DbResTypes = iota // Empty/None case.
	DbResTypeUndefined
	DbResTypeTablesNamespace
	DbResTypeTable
	DbResTypeFileNamespace
<<<<<<< Updated upstream
=======
	DbResTypeFile
>>>>>>> Stashed changes
)

const (
	RestMethodNone RestMethods = iota // Empty/None case.
	RestMethodGet
	RestMethodPost
	RestMethodPut
	RestMethodPatch
	RestMethodDelete
)

func (dbResType DbResTypes) String() string {
	switch dbResType {
	case DbResTypeNone:
		return "DbResTypeNone"
	case DbResTypeUndefined:
		return "DbResTypeUndefined"
	case DbResTypeTablesNamespace:
		return "DbResTypeTablesNamespace"
	case DbResTypeTable:
		return "DbResTypeTable"
	case DbResTypeFileNamespace:
		return "DbResTypeFileNamespace"
<<<<<<< Updated upstream
=======
	case DbResTypeFile:
		return "DbResTypeFile"
>>>>>>> Stashed changes
	}
	return "DbResTypeUnknown"
}
