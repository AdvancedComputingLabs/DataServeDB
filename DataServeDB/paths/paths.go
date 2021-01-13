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

package paths

import (
	"log"
	"os"

	"path/filepath"
)

const dataDirNameRelative = "Databases"
const configDirNameRelative = "Config"

//returns true if path was created; panics if there is dir creation error
func CreatePathIfNotExist(path string) bool {
	path = filepath.Dir(path)

	_, e := os.Stat(path)
	if os.IsNotExist(e) {
		e2 := os.MkdirAll(path, 0755)
		if e2 != nil {
			panic(e2)
		}
		return true
	}
	return false
}

func ConstructDbPath(dbName, dbsPath string) string {
	return filepath.Join(dbsPath, dbName)
}

func Combine(paths ...string) string {
	return filepath.Join(paths...)
}

func GetConfigDirPath() string {
	path, e := filepath.Abs(configDirNameRelative)
	if e != nil {
		log.Fatal(e)
	}
	return path
}

func GetDatabasesMainDirPath() string {
	path, e := filepath.Abs(dataDirNameRelative)
	if e != nil {
		log.Fatal(e)
	}
	return path
}

func GetWorkingDirPath() string {
	wd, e := os.Getwd()
	if e != nil {
		log.Fatal(e)
	}
	return wd
}

func GetExeDirPath() string {
	exe, e := os.Executable()
	if e != nil {
		log.Fatal(e)
	}
	return filepath.Dir(exe)
}
