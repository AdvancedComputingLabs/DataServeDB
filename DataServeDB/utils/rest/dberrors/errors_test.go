package dberrors

import (
	"fmt"
	"testing"
)

func TestDbErrorToErrorConversion(t *testing.T) {

	{
		dbbErr := NewDbError(TestError, fmt.Errorf("test error"))

		if err := dbbErr; err.ToError().Error() != "TestError: test error" {
			t.Fatal("DbError to error conversion failed")
		}
	}

	{
		//_, err := nilError()
		//
		//num, err := nilDbError()
		//_ = num
		//
		//// ok
		//if err.(*DbError) != nil {
		//	t.Fatal("DbError to error conversion failed#1")
		//}
		//
		//// this fails
		//if err != nil {
		//	t.Fatal("DbError to error conversion failed#2")
		//}
	}
}

func nilDbError() (int, *DbError) {
	return 0, nil
}

func nilError() (int, error) {
	return 0, nil
}
