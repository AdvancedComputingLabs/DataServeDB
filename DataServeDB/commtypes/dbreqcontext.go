package commtypes

import (
	"DataServeDB/comminterfaces"
	"DataServeDB/dbsystem/constants"
)

// fields are public, easier.

type DbReqContext struct {
	RestMethod        string
	RestMethodId      constants.RestMethods
	ResPath           string
	MatchedPath       string
	DbName            string
	Dbi               comminterfaces.DbPtrI
	TargetName        string
	TargetDbResTypeId constants.DbResTypes
}

func NewDbReqContext(restMethod, resPath, matchedPath, dbName string, dbi comminterfaces.DbPtrI, targetName string, targetDbResTypeId constants.DbResTypes) *DbReqContext {

	dbreqCtx := DbReqContext{
		RestMethod:        restMethod,
		ResPath:           resPath,
		MatchedPath:       matchedPath,
		DbName:            dbName,
		Dbi:               dbi,
		TargetName:        targetName,
		TargetDbResTypeId: targetDbResTypeId,
	}

	return &dbreqCtx
}