// Copyright (c) 2019 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package dbtable

import (
	"errors"
	"fmt"
	"testing"

	"DataServeDB/dbsystem"
	"DataServeDB/dbtypes"
)

// 1. add tests.
// 1.a. add fields with different names, should pass.
// 1.b. add field with same name, should fail.
// 2. Field Ids tests.
// 2.a. add with unique field ids, should pass.
// 2.b. add with same field id, should fail.
// 2.c. add fields with -1 id (should assigning ids internally), should pass.
// 3. Get tests.
// 3.a. Get non existing, error.
// 3.b Get existing, no error.
// 4. remove field tests.
// 4.a. RemoveField_ByName. TODO: needs to delete whole field. But should be in its own test file.
// 5. Update field name tests.
// 5.a. remove field name and test if field name has been changed.

//1.a
func TestAddFieldsName_Normal(t *testing.T) {

	dbFieldRule := tableFieldStruct{
		FieldInternalId: 0,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: 1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}
}

//1.b
func TestAddFieldsName_SameFieldNameMorethanOnce(t *testing.T) {

	dbFieldRule := tableFieldStruct{
		FieldInternalId: 0,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: 1,
		FieldName:       "Id",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err == nil {
		t.Error("This should fail.")
	}
}

//2.a
func TestAddFieldsIds_Normal(t *testing.T) {
	dbFieldRule := tableFieldStruct{
		FieldInternalId: 0,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: 1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}
}

//2.b
func TestAddFieldsIds_SameIds(t *testing.T) {
	dbFieldRule := tableFieldStruct{
		FieldInternalId: 0,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: 0,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err == nil {
		//handled with panic handler.
	}
}

//2.c
func TestAddFieldsIds_InternalIds(t *testing.T) {
	dbFieldRule := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}
}

//3.a
func TestGetFieldMeta_FieldNonExistent(t *testing.T) {
	tb01 := newTableMain(01, "Tbl01")

	if _, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("Id", dbsystem.SystemCasingHandler); err == nil {
		t.Error("This should be not nil.")
	}
}

//3.b
func TestGetFieldMeta_Normal(t *testing.T) {
	dbFieldRule := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
		//IsPk: true, //TODO: this part changed, need to redo the test with Pk.
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if fmd, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("Id", dbsystem.SystemCasingHandler); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%v\n", fmd)
	}
}

//3.a
func TestRemoveFieldMeta_Normal(t *testing.T) {
	//Deletes existing field metadata.

	dbFieldRule := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
		//IsPk: true, //TODO: this part changed, need to redo the test with Pk.
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if fmd, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("Id", dbsystem.SystemCasingHandler); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%v\n", fmd)
	}

	if err := tb01.TableFieldsMetaData.remove("Id", dbsystem.SystemCasingHandler); err != nil {
		t.Error(err)
	}

	if _, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("Id", dbsystem.SystemCasingHandler); err != nil {
		fmt.Println(err.Error())
	}
}

func TestRemoveFieldMeta_NonExistant(t *testing.T) {

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.remove("Id", dbsystem.SystemCasingHandler); err != nil {
		fmt.Println(err)
	}
}

func TestUpdateFieldMeta_Normal(t *testing.T) {
	dbFieldRule := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "Id",
		FieldType:       dbtypes.Int32,
		//IsPk: true, //TODO: this part changed, need to redo the test with Pk.
	}

	dbFieldRule2 := tableFieldStruct{
		FieldInternalId: -1,
		FieldName:       "IsTrue",
		FieldType:       dbtypes.Bool,
	}

	tb01 := newTableMain(01, "Tbl01")

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if err := tb01.TableFieldsMetaData.add(&dbFieldRule2, dbsystem.SystemCasingHandler); err != nil {
		t.Error("This should pass.")
	}

	if fmd, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("Id", dbsystem.SystemCasingHandler); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%v\n", fmd)
	}

	if err := tb01.TableFieldsMetaData.updateFieldName("Id", "NewId", dbsystem.SystemCasingHandler); err != nil {
		t.Error(err)
	}

	if _, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("Id", dbsystem.SystemCasingHandler); err != nil {
		fmt.Println(err)
	}

	if fmd, err := tb01.TableFieldsMetaData.getFieldMetadataInternal("NewId", dbsystem.SystemCasingHandler); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%v\n", fmd)
		if fmd.FieldName != "NewId" {
			t.Error(errors.New("field name did not change in field's metadata properties"))
		}
	}
}
