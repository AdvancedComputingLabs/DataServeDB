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

package db_system_rules

import (
	"fmt"
	"regexp"

	"DataServeDB/dbsystem"
)

// NOTE: db system rules like table naming, field naming are in central location rather than individual packages.
// Easier to manage changes from single rules package.
// These are hard coded into system so they will never be loaded from config files (currently as the system stands).

// NOTE: Naming is very restrictive to make changes forward compatible. It is easier to relax rules compared to tightening rules which will make previous naming incompatible.

//TODO: refactor it to types package
type zeroMem struct{}

//It will take a little more memory and cpu cycle, but the convience worth it.
var syscasing = dbsystem.SystemCasingHandler.Convert

// Database naming rules
// -

//NOTE: used in db router so must be public.
const DbNameValidatorRuleReStrBasic = "[A-Za-z][_0-9A-Za-z]{2,49}"

var dbNameValidatorRe = regexp.MustCompile(fmt.Sprintf("^%s$", DbNameValidatorRuleReStrBasic))

var dbNameReservedWords = map[string]zeroMem{
	//don't have any yet
}

func DbNameIsValid(name string) bool {
	if !dbNameValidatorRe.MatchString(name) {
		return false
	}
	if _, reserved := dbNameReservedWords[syscasing(name)]; reserved {
		return false
	}
	return true
}

// Table naming rules:
// - Len: 3 .. 50; Conservative length for now, might increase in future.
// - Casing: insensitive
// - Characters Allowed: alphanumeric starting with a letter
// - Regex: "^[A-Za-z][0-9A-Za-z]{2,49}$"
// - Reserved Words: table, tables

const TableNameValidatorRuleReStrBasic = "[A-Za-z][0-9A-Za-z]{2,49}"
var tblNameValidatorRe = regexp.MustCompile(fmt.Sprintf("^%s$", TableNameValidatorRuleReStrBasic))

var tableNameReservedWords = map[string]zeroMem{
	syscasing("table"): {},
	syscasing("tables"): {},
}

func TableNameIsValid(s string) bool {
	if !tblNameValidatorRe.MatchString(s) {
		return false
	}
	if _, reserved := tableNameReservedWords[syscasing(s)]; reserved {
		return false
	}
	return true
}

//Table field/column naming rules:
// - len: 1 .. 64; 64 character limit is used to keep table fields compatible with sql databases https://dev.mysql.com/doc/refman/8.0/en/identifier-length.html
// - Casing: insensitive
// - Characters Allowed: alphanumeric starting with a letter or underscore
// - Regex: "^[_A-Za-z][0-9A-Za-z]{0,63}$"
// - Reserved Words: PartitionKey, Timestamp

var tableFieldNameReservedWords = map[string]zeroMem{
	syscasing("PartitionKey"): {},
	syscasing("Timestamp"): {},
}

func TableFieldNameIsValid(s string) bool {
	if !regexp.MustCompile("^[_A-Za-z][0-9A-Za-z]{0,63}$").MatchString(s) {
		return false
	}
	if _, reserved := tableFieldNameReservedWords[syscasing(s)]; reserved {
		return false
	}
	return true
}

