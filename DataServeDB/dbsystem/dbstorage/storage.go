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
	"io/ioutil"
	"os"

	"DataServeDB/paths"
)

func SaveToDisk(data []byte, path string) error {
	//println(string(data))

	//no need to check if path was created here.
	paths.CreatePathIfNotExist(path)

	fo, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = fo.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromDisk(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
