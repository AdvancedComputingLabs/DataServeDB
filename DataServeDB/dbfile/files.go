package dbfile

import (
	"DataServeDB/paths"
	"DataServeDB/utils/rest/dberrors"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

func ListFiles() (result []byte, dberr *dberrors.DbError) {
	var fileNames []string
	files, err := ioutil.ReadDir(paths.GetFilesPath())
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
func GetFile(filename string) (result []byte, dberr *dberrors.DbError) {
	fileName := paths.Combine(paths.GetFilesPath(), filename)
	result, err := os.ReadFile(fileName)
	if err != nil {
		dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}

	return
}

func PostFile(multipartForm *multipart.Form) (dberr *dberrors.DbError) {
	for _, fileHeaders := range multipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			//path := fmt.Sprintf("files/%s", fileHeader.Filename)
			fileName := paths.Combine(paths.GetFilesPath(), fileHeader.Filename)
			buf, _ := ioutil.ReadAll(file)
			err := ioutil.WriteFile(fileName, buf, os.ModePerm)
			if err != nil {
				dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
				return dberr
			}
		}
	}
	return nil
}

func DeleteFile(filename string) (dberr *dberrors.DbError) {
	fileName := paths.Combine(paths.GetFilesPath(), filename)
	err := os.Remove(fileName)
	if err != nil {
		return dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
	}
	return
}

func EditOrUpdateFile(fileName string, multipartForm *multipart.Form) (dberr *dberrors.DbError) {

	for _, fileHeaders := range multipartForm.File {

		//TO DO-:  detect single file
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			fileName := paths.Combine(paths.GetFilesPath(), fileHeader.Filename)
			buf, _ := ioutil.ReadAll(file)
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
				return
			}
			if _, err = f.Write(buf); err != nil {
				dberr = dberrors.NewDbError(dberrors.InvalidInputKeyNotProvided, err)
				return
			}

			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}
