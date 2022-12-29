package main

import (
	"DataServeDB/unstable_api/dbrouter"
	"DataServeDB/utils/rest"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
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

	successResult, err := restApiCallMu("POST", "re_db/files/level1/")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(successResult)
	}
}

func TestDeleteFile(t *testing.T) {
	successResult, err := restApiCall("DELETE", "re_db/files/level1/level2/new12.txt", "")
	if err != nil {
		// not implemented yet
		//log.Fatal(err)
		log.Println(err)
	} else {
		log.Println(successResult)
	}
}

func restApiCallMu(method, path string) (string, error) {
	//bodybytes := new(bytes.Buffer)
	// file, er := os.ReadFile("./tes.txt")
	// if er != nil {
	// 	fmt.Println(er.Error())
	// }
	// file, err := os.Open("tes.txt")
	// if err != nil {
	// 	//t.Fatal(err)
	// }
	// data := url.Values{}
	// data.Set("name", file.)
	wbody := &bytes.Buffer{}
	writer := multipart.NewWriter(wbody)
	fw, err := writer.CreateFormFile("photo", "token.json")
	if err != nil {
	}
	file, err := os.Open("token.json")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Close multipart writer.
	writer.Close()

	req, w := newHttpReqNResp(method, path, bytes.NewReader(wbody.Bytes()))

	reqPath := rest.HttpRestPathParser(req.URL.String())
	// req.Header.Add("Content-Type", "multipart/form-data; boundary=<calculated when request is sent>")
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
func restApiPost(method, path string) (string, error) {
	file, _ := os.ReadFile("./tes.txt")
	data := url.Values{}
	data.Set("name", string(file))
	//data.Set("surname", "bar")

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode())) // URL-encoded payload

	resp, _ := client.Do(r)
	//body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return "", fmt.Errorf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, "")
	} else {
		return fmt.Sprintf("\n\tstatus-code: %v\n\tresponse: %v", resp.StatusCode, ""), nil
	}
}
