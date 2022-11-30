package db

import (
	"errors"
	"fmt"
	"mime/multipart"

	"DataServeDB/commtypes"
	"DataServeDB/dbsystem/constants"
	"DataServeDB/files"
	"DataServeDB/utils/rest"
	"DataServeDB/utils/rest/dberrors"
)

func (d *DB) FilesGet(dbReqCtx *commtypes.DbReqContext) (int, []byte, error) {

	//do body to string
	effectivePathLevel := getEffectivePathLevel(dbReqCtx.PathLevels)

	fmt.Println("Type ID -> ", effectivePathLevel.PathItemTypeId)

	switch effectivePathLevel.PathItemTypeId {
	case constants.DbResTypeFileNamespace:
		{
			resultContent, err := files.ListFiles()
			if err != nil {
				return rest.HttpRestDbError(err)
			}
			return rest.HttpRestOk(resultContent)

		}
	case constants.DbResTypeFile:
		{
			fileName := effectivePathLevel.PathItem
			//fmt.Println("filename", fileName)

			resultContent, dberr := files.GetFile(fileName)
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

			dberr := files.PostFile(multipartForm)
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

			if dberr := files.DeleteFile(fileName); dberr != nil {
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

			if dberr := files.EditOrUpdateFile(fileName, multipartForm); dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOkNoContent()
		}
	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}
