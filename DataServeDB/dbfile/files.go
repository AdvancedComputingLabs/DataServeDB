package dbfile

import (
	rules "DataServeDB/dbsystem/rules"
	"DataServeDB/paths"
	"DataServeDB/utils/rest/dberrors"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ListFiles(path string) (result []byte, dberr *dberrors.DbError) {
	var fileNames []string
	path = GetPath(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
		return
	}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	result, err = json.Marshal(fileNames)
	if err != nil {
		dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}

	return
}
func GetFile(path, filename string) (result []byte, dberr *dberrors.DbError) {

	filePath := GetPath(path)
	//VerifyFileStorage(filePath, filename)
	result, err := os.ReadFile(filePath)
	if err != nil {
		dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}

	return
}

func PostFile(path string, multipartForm *multipart.Form) (dberr *dberrors.DbError) {
	path = GetPath(path)
	for _, fileHeaders := range multipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			//path := fmt.Sprintf("files/%s", fileHeader.Filename)
			if fnerr := matchFileName(fileHeader.Filename); fnerr != nil {
				return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, fnerr)
			}
			fileName := paths.Combine(path, fileHeader.Filename)
			data, _ := ioutil.ReadAll(file)

			fo, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_EXCL, 0755)
			if err != nil {
				return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
			}
			defer fo.Close()
			_, err = fo.Write(data)
			if err != nil {
				return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
			}

		}
	}
	return nil
}

func DeleteFile(path string) (dberr *dberrors.DbError) {
	filePath := GetPath(path)
	err := os.Remove(filePath)
	if err != nil {
		return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}
	return
}

func EditOrUpdateFile(path string, multipartForm *multipart.Form) (dberr *dberrors.DbError) {

	filePath := GetPath(path)
	for _, fileHeaders := range multipartForm.File {
		//TO DO-:  detect single file
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			fileName := paths.Combine(filePath, fileHeader.Filename)
			buf, _ := ioutil.ReadAll(file)
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
			}
			if _, err = f.Write(buf); err != nil {
				return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
			}

			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

// func VerifyFilePathLevels(path string) error {

// 	return nil
// }

func VerifyFileStorage(path, filename string) *dberrors.DbError {
	var size int64

	dirPath := strings.Trim(path, filename)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}
	size, err = dirSize(dirPath)
	if err != nil {
		return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}
	if len(files) >= 10000 || size >= 1024*1024*1024*100 {
		return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, errors.New("out of storage"))
	}

	return nil
}
func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func extractFileName(path string) string {
	re := regexp.MustCompile(rules.FileNameValidator)
	return re.FindString(path)
}
func GetPath(path string) string {
	path = strings.Replace(path, "files", "files_data", 1)
	return paths.Combine(paths.GetFilesPath(), path)
}

func matchFileName(fileName string) error {
	re := regexp.MustCompile(rules.FileNameValidator)
	res := re.FindString(fileName)

	if res == fileName {
		return nil
	}
	return errors.New("file name not accepted")
}
