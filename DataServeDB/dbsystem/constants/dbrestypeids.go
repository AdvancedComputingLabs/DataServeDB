package constants

type DbResTypes int

type RestMethods int

const (
	DbResTypeNone DbResTypes = iota // Empty/None case.
	DbResTypeTable
	DbResTypeFile
)

const (
	RestMethodNone RestMethods = iota // Empty/None case.
	RestMethodGet
	RestMethodPut
)
