package main

import (
	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"testing"
)

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
	successResult, err := restApiCall("GET", "re_db/files/new.txt", "")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func TestPostFile(t *testing.T) {
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)

	file, err := os.Open("tes.txt")
	if err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile("file", "tes.txt")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}
	defer mw.Close()

	successResult, err := restApiCallMu("POST", "re_db/files", mw)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func restApiCallMu(method, path string, mw *multipart.Writer) (string, error) {
	bodybytes := new(bytes.Buffer)

	req, w := newHttpReqNResp(method, path, io.NopCloser(bodybytes))

	reqPath := rest.HttpRestPathParser(req.URL.String())
	req.Header.Add("Content-Type", mw.FormDataContentType())

	dbrouter.MatchPathAndCallHandler(w, req, reqPath, req.Method)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return "", fmt.Errorf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, string(body))
	} else {
		return fmt.Sprintf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, string(body)), nil
	}
}
