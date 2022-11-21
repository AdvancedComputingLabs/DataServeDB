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

package dbstorage

import (
	"os"
	"path/filepath"

	"DataServeDB/paths"
)

func DeleteDirFromDisk(path string) error {
	// TODO: check if it may not remove all the files, then it will create table in inconsistent state.
	//  Handle it correctly.
	dir := filepath.Dir(path)
	return os.RemoveAll(dir)
}

func SaveToDisk(data []byte, path string) error {
	//println(string(data))

	//no need to check if path was created here.
	paths.CreatePathIfNotExist(path)

	fo, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	if err != nil {
		return err
	}
	defer fo.Close() // TODO: handle error

	_, err = fo.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromDisk(path string) ([]byte, error) {
	return os.ReadFile(path)
}
