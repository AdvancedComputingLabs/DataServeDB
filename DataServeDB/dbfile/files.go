package dbfile

import (
	"DataServeDB/paths"
	"DataServeDB/utils/rest/dberrors"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
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
func GetPath(path string) string {
	path = strings.Replace(path, "files", "files_data", 1)
	return paths.Combine(paths.GetFilesPath(), path)
}

func PostFile(path string, multipartForm *multipart.Form) (dberr *dberrors.DbError) {
	path = GetPath(path)
	for _, fileHeaders := range multipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			//path := fmt.Sprintf("files/%s", fileHeader.Filename)
			fileName := paths.Combine(path, fileHeader.Filename)
			buf, _ := ioutil.ReadAll(file)
			err := ioutil.WriteFile(fileName, buf, os.ModePerm)
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
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

func VerifyFilePathLevels(path string) error {

	return nil
}

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
