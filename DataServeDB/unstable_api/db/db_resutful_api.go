package db

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"DataServeDB/commtypes"
	"DataServeDB/dbsystem/constants"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
	"DataServeDB/utils/rest/dberrors"
)

func (d *DB) TablesPatch(dbReqCtx *commtypes.DbReqContext, body io.ReadCloser) (int, []byte, error) {
	return patchOrPut(d, dbReqCtx, body, idbstorer.TableOperationPatchRow)
}

func (d *DB) TablesPut(dbReqCtx *commtypes.DbReqContext, body io.ReadCloser) (int, []byte, error) {
	return patchOrPut(d, dbReqCtx, body, idbstorer.TableOperationReplaceRow)
}

func (d *DB) TablesPost(dbReqCtx *commtypes.DbReqContext, body io.ReadCloser) (int, []byte, error) {

	//do body to string
	bodyStr, err := bodyToStr(body)
	if err != nil {
		return rest.HttpRestError(dberrors.InvalidInput, err)
	}

	lastPathLevel := getLastPathLevel(dbReqCtx.PathLevels)

	switch lastPathLevel.PathItemTypeId {
	case constants.DbResTypeTablesNamespace:
		{
			//create table

			dberr := d.CreateTableJSON(bodyStr, nil)
			if dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOkNoContent()
		}

	case constants.DbResTypeTable:
		{
			//insert data

			tableName := lastPathLevel.PathItem

			table, err := d.GetTable(tableName)
			if err != nil {
				return rest.HttpRestError(dberrors.TableNotFound, err)
			}
			dberr := table.InsertRowJSON(bodyStr)
			if dberr != nil {
				return rest.HttpRestDbError(dberr)
			}
			return rest.HttpRestOkNoContent()
		}

	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}

}

func (d *DB) TablesDelete(dbReqCtx *commtypes.DbReqContext) (int, []byte, error) {

	// NOTE: Always, there will be table level. There is no call to table namespace level as it is not deletable.
	tablePathLevel := getTablePathLevel(dbReqCtx.PathLevels)

	switch tablePathLevel.PathItemTypeId {

	case constants.DbResTypeTable:
		{
			// check if delete is on table or row

			if len(dbReqCtx.PathLevels) > 2 {
				// delete row

				//NOTE: there is no empty index 'indexName:' access at the moment.

				_, value, err := rest.ParseKeyValue(dbReqCtx.ResPath)
				if err != nil {
					return rest.HttpRestError(dberrors.InvalidInput, err)
				}

				if value == "" {
					return rest.HttpRestError(dberrors.InvalidInputKeyNotProvided, errors.New("primary key or secondary key is missing"))
				}

				tableName := tablePathLevel.PathItem

				table, err := d.GetTable(tableName)
				if err != nil {
					return rest.HttpRestError(dberrors.TableNotFound, err)
				}

				if dberr := table.DeleteRowByPrimaryKey(value); dberr != nil {
					return rest.HttpRestDbError(dberr)
				}

				return rest.HttpRestOkNoContent()
			} else {
				// delete table

				tableName := tablePathLevel.PathItem

				if dberr := d.DeleteTable(tableName); dberr != nil {
					return rest.HttpRestDbError(dberr)
				}

				return rest.HttpRestOkNoContent()
			}
		}

	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}

func (d *DB) TablesGet(dbReqCtx *commtypes.DbReqContext) (int, []byte, error) {

	effectivePathLevel := getEffectivePathLevel(dbReqCtx.PathLevels)

	switch effectivePathLevel.PathItemTypeId {
	case constants.DbResTypeTablesNamespace:
		{
			// get all tables

			// TODO:
			//tables, err := d.GetAllTables()

			return rest.HttpRestError(dberrors.NotImplemented, errors.New("not implemented"))
		}

	case constants.DbResTypeTable:
		{
			// get row or rows of a table

			tableName := effectivePathLevel.PathItem

			table, err := d.GetTable(tableName)
			if err != nil {
				return rest.HttpRestError(dberrors.TableNotFound, err)
			}

			resultContent, dberr := table.Get(dbReqCtx)
			if dberr != nil {
				return rest.HttpRestDbError(dberr)
			}

			return rest.HttpRestOk(resultContent)
		}

	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}

}

// ### Private Methods/Functions

func bodyToStr(body io.ReadCloser) (string, error) {

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func extractFullCallStrTablesParent(s, targetResName string) string {
	pos := strings.LastIndex(s, targetResName)
	if pos == -1 {
		return ""
	}

	result := s[pos:]
	if result != "" {
		if result[len(result)-1] == '/' {
			result = result[:len(result)-1]
		}
	}

	return result
}

func extractTableNameFromFullCallStrTableParent(fullCallStr string) string {
	re := regexp.MustCompile(`\('(.+)'\)`)

	matches := re.FindStringSubmatch(fullCallStr)
	if len(matches) != 2 {
		return ""
	}

	return matches[1]
}

func getEffectivePathLevel(levels []dbrouter.PathLevel) dbrouter.PathLevel {

	if len(levels) == 0 {
		return dbrouter.PathLevel{}
	}

	// check namespace level, always the first level

	switch levels[0].PathItemTypeId {
	case constants.DbResTypeTablesNamespace:
		{
			if len(levels) == 1 {
				return levels[0]
			}

			// check table level, always the second level
			if levels[1].PathItemTypeId == constants.DbResTypeTable {
				return levels[1]
			}

			// any further levels are not needed in the case of table namespace.
		}
	}

	// no effective path level found, return empty path level; does this have performance impact? then change it to return nil.
	// at the moment, nil avoids the need to check for nil in the caller. -- HY @19-Nov-2022
	return dbrouter.PathLevel{}
}

func getLastPathLevel(levels []dbrouter.PathLevel) dbrouter.PathLevel {

	if len(levels) == 0 {
		return dbrouter.PathLevel{}
	}

	return levels[len(levels)-1]
}

func getTablePathLevel(levels []dbrouter.PathLevel) dbrouter.PathLevel {
	const tablePathLevel = 2 // always second level starting from 1; 1 if starting from 0

	if len(levels) < tablePathLevel {
		return dbrouter.PathLevel{}
	}

	return levels[tablePathLevel-1]
}

func patchOrPut(d *DB, dbReqCtx *commtypes.DbReqContext, body io.ReadCloser, updateType idbstorer.TableOperationType) (int, []byte, error) {

	//do body to string
	bodyStr, err := bodyToStr(body)
	if err != nil {
		return rest.HttpRestError(dberrors.InvalidInput, err)
	}

	// NOTE: at the moment only row level patch/put has been implemented.
	// TODO: other levels patch/put support?

	tablePathLevel := getTablePathLevel(dbReqCtx.PathLevels)

	switch tablePathLevel.PathItemTypeId {
	case constants.DbResTypeTable:
		{
			// patch/put row or table?

			if len(dbReqCtx.PathLevels) > 2 {
				// patch/put row

				colName, value, err := rest.ParseKeyValue(dbReqCtx.ResPath)
				if err != nil {
					return rest.HttpRestError(dberrors.InvalidInput, err)
				}

				//NOTE: there is no empty index 'indexName:' access at the moment.
				if value == "" {
					return rest.HttpRestError(dberrors.InvalidInputKeyNotProvided, errors.New("primary key or secondary index key is missing"))
				}

				tableName := tablePathLevel.PathItem

				//update row
				table, err := d.GetTable(tableName)
				if err != nil {
					return rest.HttpRestError(dberrors.TableNotFound, err)
				}

				if colName == "" {
					// primary key
					dberr := table.UpdateRowJsonByPk(value, bodyStr, updateType)
					if dberr != nil {
						return rest.HttpRestDbError(dberr)
					}
					return rest.HttpRestOkNoContent()
				} else {
					// secondary index key
					return rest.HttpRestError(dberrors.InvalidInput, errors.New("secondary indexes are not supported"))
				}
			} else {
				// patch/put table
				return rest.HttpRestError(dberrors.NotImplemented, errors.New("table patch/put is not implemented"))
			}

		}

	default:
		return rest.HttpRestError(dberrors.InvalidInput, errors.New("invalid target"))
	}
}
