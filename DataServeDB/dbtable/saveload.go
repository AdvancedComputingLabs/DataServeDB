// Copyright (c) 2020 Advanced Computing Labs DMCC

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
	"encoding/json"
	"errors"
)

// Description: dbtable package saving and loading file.

// !WARNING: unstable api, but this needs to be in dbtable package.

type DbTableRecreation struct {
	TableInternalId   int
	CreationStructure createTableExternalStruct
}

//!INFO: this is just for demonstration, tune it if you think there is better way to do it.
func GetSaveLoadStructure(dbtbl *DbTable) (string, error) {
	slStruct := DbTableRecreation{}

	slStruct.TableInternalId = dbtbl.tblMain.TableId
	slStruct.CreationStructure = dbtbl.createTableStructure

	//TODO: handle error
	if r, e := json.Marshal(slStruct); e == nil {
		return string(r), nil
	}

	return "", errors.New("did not convert to json")
}

func LoadFromJson(dbtblJson string) (*DbTable, error) {
	var slStruct DbTableRecreation

	if e := json.Unmarshal([]byte(dbtblJson), &slStruct); e != nil {
		//TODO: for later after version 0.5, return structured error, top error json error and in sub structure include the json message.
		return nil, e
	}

	//TODO: it creates new, but data needs to be attached, second it is keeping all in memory (for now this is ok).
	//TODO: add secondary index support
	tdc := &tableDataContainer{
		Rows:          nil,
		PkToRowMapper: map[interface{}]int64{},
	}

	dbtbl, e := createTable(slStruct.TableInternalId, &slStruct.CreationStructure, tdc)
	if e != nil {
		return nil, e
	}

	return dbtbl, nil
}

