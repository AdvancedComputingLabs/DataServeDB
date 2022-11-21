// Copyright (c) 2022 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

/*
	WARNING:
		This is only for testing functionality. Not meant for production use.
*/

package dbtable

import (
	"errors"

	"DataServeDB/commtypes"
	idbstorer "DataServeDB/storers/dbtable_interface_storer"
)

const storerMemAppendLogName = "MemAppendLogTestsOnly"

type MALRow struct {
	ColName  string
	ColValue any
}

type MemAppendLog struct {
	Rows [][]MALRow
}

//func NewMemAppendLog(tableId int, tableFolderPath string) (idbstorer.StorerBasic, error) {
//	return &MemAppendLog{}, nil
//}

func (m MemAppendLog) Implemented(feature idbstorer.TableStorerFeaturesType) bool {

	switch feature {
	case idbstorer.TableStorerFeature_Insert:
		return true
	}

	return false
}

func (m MemAppendLog) Delete(indexName string, key string) (int, error) {
	return -1, errors.New("NotImplemented")
}

func (m MemAppendLog) Get(indexName string, key string) (int, any, error) {
	return -1, nil, errors.New("NotImplemented")
}

func (m *MemAppendLog) Insert(rowWithProps commtypes.TableRowWithFieldProperties, data any) (int, error) {
	var row []MALRow

	for _, v := range rowWithProps {
		rowCol := MALRow{}
		rowCol.ColName = v.Name()
		rowCol.ColValue = v.Value()
		row = append(row, rowCol)
	}

	m.Rows = append(m.Rows, row)

	return 1, nil
}

func (m MemAppendLog) Update(rowWithProps commtypes.TableRowWithFieldProperties, data any) (int, error) {
	return -1, errors.New("NotImplemented")
}
