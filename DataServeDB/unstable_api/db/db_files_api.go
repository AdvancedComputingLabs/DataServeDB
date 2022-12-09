package db

import (
	"errors"
	"fmt"
	"mime/multipart"

	"DataServeDB/commtypes"
	"DataServeDB/dbfile"
	"DataServeDB/dbsystem/constants"
	"DataServeDB/utils/rest"
	"DataServeDB/utils/rest/dberrors"
)

func (d *DB) FilesGet(dbReqCtx *commtypes.DbReqContext) (int, []byte, error) {

	//do body to string
	effectivePathLevel := getEffectivePathLevel(dbReqCtx.PathLevels)

	switch effectivePathLevel.PathItemTypeId {
	case constants.DbResTypeFileNamespace:
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
			fmt.Println(fileName, dbReqCtx.MatchedPath)

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
	case constants.DbResTypeFileNamespace:
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

	tablePathLevel := getTablePathLevel(dbReqCtx.PathLevels)

	switch tablePathLevel.PathItemTypeId {

	case constants.DbResTypeTable:
		{
			fileName := tablePathLevel.PathItem

			if dberr := dbfile.DeleteFile(fileName); dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOkNoContent()
		}
	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}

func (d *DB) FilesPutorPatch(dbReqCtx *commtypes.DbReqContext, multipartForm *multipart.Form) (int, []byte, error) {
	tablePathLevel := getTablePathLevel(dbReqCtx.PathLevels)

	switch tablePathLevel.PathItemTypeId {

	case constants.DbResTypeTable:
		{
			fileName := tablePathLevel.PathItem

			if dberr := dbfile.EditOrUpdateFile(fileName, multipartForm); dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOkNoContent()
		}
	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}
