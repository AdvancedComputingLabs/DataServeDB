package runtime

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeJSONBody(t *testing.T) {

	//Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"name":  "Tomy",
		"email": "Tomy@example.com",
	})
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", "http://localhost:8080/", responseBody)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	//Leverage Go's HTTP Post function to make request
	// resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
	// //Handle Error
	// if err != nil {
	// 	log.Fatalf("An Error Occured %v", err)
	// }
	println("hello")
	rec := httptest.NewRecorder()
	decodeJSONBody(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
