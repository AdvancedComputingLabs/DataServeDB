package dberrors

import (
	"errors"
	"net/http"
)

type RestError string

// NOTE: arranged in status code order
const (
	TestError                             RestError = "TestError"                             // DbErrorCode: TestError; HttpStatusCode: 50 (no http status code for this)
	InvalidInput                                    = "InvalidInput"                          // DbErrorCode: InvalidInput; HttpStatusCode: 400 StatusBadRequest
	InvalidInputColumnNameDoesNotExist              = "InvalidInputColumnNameDoesNotExist"    // DbErrorCode: InvalidInputColumnNameDoesNotExist; HttpStatusCode: 400 StatusBadRequest
	InvalidInputDuplicateKey                        = "InvalidInputDuplicateKey"              // DbErrorCode: InvalidInputDuplicateKey; HttpStatusCode: 400 StatusBadRequest
	InvalidInputHttpMethodNotSupported              = "InvalidInputHttpMethodNotSupported"    // DbErrorCode: InvalidInputHttpMethodNotSupported; HttpStatusCode: 400 StatusBadRequest
	InvalidInputKeyNotProvided                      = "InvalidInputKeyNotProvided"            // DbErrorCode: InvalidInputKeyNotProvided; HttpStatusCode: 400 StatusBadRequest
	InvalidInputPrimaryKeyCannotBeUpdated           = "InvalidInputPrimaryKeyCannotBeUpdated" // DbErrorCode: InvalidInputPrimaryKeyCannotBeUpdated; HttpStatusCode: 400 StatusBadRequest
	StorerConstraintViolation                       = "StorerConstraintViolation"             // DbErrorCode: StorerConstraintViolation; HttpStatusCode: 400 StatusBadRequest
	DatabaseNotFound                                = "DatabaseNotFound"                      // DbErrorCode: DatabaseNotFound; HttpStatusCode: 404 StatusNotFound
	KeyNotFound                                     = "KeyNotFound"                           // DbErrorCode: KeyNotFound; HttpStatusCode: 404 StatusNotFound
	TableNotFound                                   = "TableNotFound"                         // DbErrorCode: TableNotFound; HttpStatusCode: 404 StatusNotFound
	TableIdAlreadyExists                            = "TableIdAlreadyExists"                  // DbErrorCode: TableIdAlreadyExists; HttpStatusCode: 409 StatusConflict
	TableNameAlreadyExists                          = "TableNameAlreadyExists"                // DbErrorCode: TableNameAlreadyExists; HttpStatusCode: 409 StatusConflict
	InternalServerError                             = "InternalServerError"                   // DbErrorCode: InternalServerError; HttpStatusCode: 500 StatusInternalServerError
	InternalServerErrorDiskError                    = "InternalServerErrorDiskError"          // DbErrorCode: InternalServerErrorDiskError; HttpStatusCode: 500 StatusInternalServerError
	NotImplemented                                  = "NotImplemented"                        // DbErrorCode: NotImplemented; HttpStatusCode: 501 StatusNotImplemented
)

type DbError struct {
	ErrCode RestError
	Err     error
}

func GetStatusCode(errCode RestError) int {

	// TODO: check all error codes have been added here

	switch errCode {
	case TestError:
		return 50
	case InvalidInput, InvalidInputColumnNameDoesNotExist, InvalidInputDuplicateKey, InvalidInputHttpMethodNotSupported,
		InvalidInputKeyNotProvided, InvalidInputPrimaryKeyCannotBeUpdated, StorerConstraintViolation:
		return http.StatusBadRequest // 400
	case DatabaseNotFound, KeyNotFound, TableNotFound:
		return http.StatusNotFound // 404
	case TableIdAlreadyExists, TableNameAlreadyExists:
		return http.StatusConflict // 409
	case InternalServerError, InternalServerErrorDiskError:
		return http.StatusInternalServerError // 500
	case NotImplemented:
		return http.StatusNotImplemented // 501
	}

	return 0
}

func NewDbError(errCode RestError, err error) *DbError {
	return &DbError{ErrCode: errCode, Err: err}
}

func (r RestError) Error() string {
	return string(r)
}

// Error kept here for testing.
//func (d *DbError) Error() string {
//	return d.ErrCode.Error() + ": " + d.Err.Error()
//}

// ToError Needed to use ToError() instead of interface method Error(), because returning nil DbError was making error == nil check faulty.
// See test TestDbErrorToErrorConversion. -- HY 19-Nov-2022
func (d *DbError) ToError() error {
	return ToError(d)
}

func ToError(d *DbError) error {
	if d == nil {
		return nil
	}
	return errors.New(d.ErrCode.Error() + ": " + d.Err.Error())
}
