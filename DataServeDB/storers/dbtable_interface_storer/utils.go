package dbtable_interface_storer

import (
	"DataServeDB/utils/rest/dberrors"
	"errors"
	"fmt"
)

func CheckTableStorageConstraints(current, prev, next StorerBasic) *dberrors.DbError {

	if current.FeaturesAndConstraints(TableStorerConstraint_MustNotBeFirstStore) && prev == nil {
		return dberrors.NewDbError(dberrors.StorerConstraintViolation, errors.New(fmt.Sprintf("storage '%s' cannot be first store", current.DisplayName())))
	}

	if current.FeaturesAndConstraints(TableStorerConstraint_MustBeLastStore) && next != nil {
		return dberrors.NewDbError(dberrors.StorerConstraintViolation, errors.New(fmt.Sprintf("storage '%s' must be last store", current.DisplayName())))
	}

	return nil
}
