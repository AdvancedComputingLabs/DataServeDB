package main

import (
	"DataServeDB/paths"
	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"testing"
)

type getTestModel struct {
	path string
	exp  error
}
type deleteTest struct {
	arg1 string
	exp  error
}

var addTests = []getTestModel{
	{"level1/level2", nil},
	{"level1", nil},
	{"level3", nil},
	{"level1/level2/new12.txt", nil},
	{"level1/level2/few12.txt", nil},
}

func TestListfilesRApi(t *testing.T) {
	successResult, err := restApiCall("GET", "re_db/files", "")
	if err != nil {
		// not implemented yet
		//log.Fatal(err)
		log.Println(err)
	} else {
		fmt.Println("file :- ")
		log.Println(successResult)
	}

}

func TestGetFileByNameRApi(t *testing.T) {

	for i, tc := range addTests {
		fmt.Println("Test Case ", i)
		successResult, err := restApiCall("GET", "re_db/files/"+tc.path, "")
		if err != tc.exp {
			log.Println(err)
		} else {
			log.Println(successResult)
		}
	}
}

const files_path = "../../../test_files/"

func TestPostFile(t *testing.T) {
	files, err := ioutil.ReadDir(files_path)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	for _, file := range files {
		t.Run("test POST"+file.Name(), func(t *testing.T) {
			successResult, err := restApiCallMu("POST", "re_db/files/level1/", file.Name())
			if err != nil {
				t.Fatal(err)
			} else {
				log.Println(successResult)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	files, err := ioutil.ReadDir(files_path)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	for _, file := range files {
		t.Run("test DELETE"+file.Name(), func(t *testing.T) {
			successResult, err := restApiCall("DELETE", "re_db/files/level1/"+file.Name(), "")
			if err != nil {
				log.Println(err)
			} else {
				log.Println(successResult)
			}
		})
	}
}

func restApiCallMu(method, path, fileName string) (string, error) {

	wbody := &bytes.Buffer{}
	writer := multipart.NewWriter(wbody)
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println(err)
	}
	fileName = paths.Combine(files_path, fileName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Close multipart writer.
	writer.Close()

	req, w := newHttpReqNResp(method, path, bytes.NewReader(wbody.Bytes()))

	reqPath := rest.HttpRestPathParser(req.URL.String())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	dbrouter.MatchPathAndCallHandler(w, req, reqPath, req.Method)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return "", fmt.Errorf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, string(body))
	} else {
		return fmt.Sprintf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, string(body)), nil
	}
}
