package db

import (
	"errors"
	"fmt"
	"mime/multipart"

	"DataServeDB/commtypes"
	"DataServeDB/dbfile"
	"DataServeDB/dbsystem/constants"
	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
	"DataServeDB/utils/rest/dberrors"
)

func (d *DB) FilesGet(dbReqCtx *commtypes.DbReqContext) (int, []byte, error) {

	//do body to string
	effectivePathLevel := getEffectivePathLevel(dbReqCtx.PathLevels)
	fmt.Println(dbReqCtx.MatchedPath, effectivePathLevel.PathItemTypeId)

	switch effectivePathLevel.PathItemTypeId {
	case constants.DbResTypeFileNamespace, constants.DbResTypeDirName:
		{
			resultContent, err := dbfile.ListFiles(dbReqCtx.MatchedPath)
			if err != nil {
				return rest.HttpRestDbError(err)
			}
			return rest.HttpRestOk(resultContent)

		}
	case constants.DbResTypeFile:
		{
			fileName := effectivePathLevel.PathItem
			//fmt.Println(fileName, dbReqCtx.MatchedPath)

			resultContent, dberr := dbfile.GetFile(dbReqCtx.MatchedPath, fileName)
			if dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			// TODO := Need to set content type
			// w.WriteHeader(http.StatusOK)
			// w.Header().Set("Content-Type", "application/octet-stream")
			// w.Write(fileBytes)

			return rest.HttpRestOk(resultContent)
		}

	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}

func (d *DB) FilesPost(dbReqCtx *commtypes.DbReqContext, multipartForm *multipart.Form) (int, []byte, error) {

	lastPathLevel := getLastPathLevel(dbReqCtx.PathLevels)

	switch lastPathLevel.PathItemTypeId {
	case constants.DbResTypeFileNamespace, constants.DbResTypeDirName:
		{
			//fileName := effectivePathLevel.PathItem

			dberr := dbfile.PostFile(dbReqCtx.MatchedPath, multipartForm)
			if dberr != nil {
				return rest.HttpRestDbError(dberr)
			}
			return rest.HttpRestOkNoContent()
		}

	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}

}

func (d *DB) FilesDelete(dbReqCtx *commtypes.DbReqContext) (int, []byte, error) {

	effectivePathLevel := getEffectivePathLevel(dbReqCtx.PathLevels)

	switch effectivePathLevel.PathItemTypeId {

	case constants.DbResTypeFile:
		{
			fileName := dbReqCtx.MatchedPath
			fmt.Println(fileName)

			if dberr := dbfile.DeleteFile(fileName); dberr != nil {
				// fmt.Println(dberr.ToError().Error())
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOkNoContent()
		}
	// ToDo :- Delete Directory
	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}

func (d *DB) FilesPutorPatch(dbReqCtx *commtypes.DbReqContext, multipartForm *multipart.Form) (int, []byte, error) {
	effectivePathLevel := getEffectivePathLevel(dbReqCtx.PathLevels)
	fmt.Println(constants.DbResTypeFile, constants.DbResTypeDirName)

	switch effectivePathLevel.PathItemTypeId {
	case constants.DbResTypeFile, constants.DbResTypeDirName:
		{
			//fileName := effectivePathLevel.PathItem
			if dberr := dbfile.EditOrUpdateFile(dbReqCtx.MatchedPath, multipartForm); dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOkNoContent()
		}
	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}

func getFilePathLevel(levels []dbrouter.PathLevel) dbrouter.PathLevel {
	const tablePathLevel = 2 // always second level starting from 1; 1 if starting from 0

	if len(levels) < tablePathLevel {
		return dbrouter.PathLevel{}
	}

	return levels[tablePathLevel-1]
}
