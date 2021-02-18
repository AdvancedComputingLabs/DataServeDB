// Copyright (c) 2018 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package unstable_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/beevik/guid"
)

type Info struct {
	Type             string
	Description      string
	Location         string
	YearBuilt        string
	YearPurchased    string
	PurchasedPrice   string
	CurrentValuation string
	Occupied         string
	Drawings         string
}

type TenantInfo struct {
	TenantId         string
	TenantName       string
	TenantPassportNo string
}
type TenancyPeriod struct {
	Start string
	End   string
}
type ContractInfo struct {
	TenancyPeriod TenancyPeriod
	EjariNo       string
	RentalAmount  string
}
type RentalInfo struct {
	TenantInfo   TenantInfo
	ContractInfo ContractInfo
}
type RentalItem struct {
	Id              string
	Description     string
	ModeOfPayment   string
	ChequeDate      string
	TransactionDate string
	Amount          string
	Comment         string
}
type Id struct {
	Property string
	Item     string
}

//TODO: code needs cleaning and restructuring. -- 02-Sep-2019 Habib Y.

/*
** type RentalItemv2:
RentalItemCategory string Examples: {SD, SecurityDeposit}, {RE, Rental}, {EX, Expense}. Explaination: Key and text in the format of {Key, Text}.
Id string Examples: SD0001, R0001, or E0005
EntryDate datetime
Description string max len 2000 characters
Date
TransactionDate
Amount
ModeOfPayment selection (Cash = 1, Cheque = 2, Online = 3). Recommendation: make package with const values for these.
ModeOfPaymentDetail interface{}. Recommendation: make struct for each type of mode of payment.
Comment
** type Cash: nil Comment no need for cash type
** type Cheque:
Number
AccountNumber
BankName
ChequeDate
TransactionDate
Comment
** type Online
ServiceName
TransactionDate
Comment
*/
type RentalItemv2 struct {
	Id             string
	EntryDate      time.Time
	Description    string
	ModeOfPayment  int
	PaymentDetails interface{}
	Amount         string
	Comment        string
}
type Cheque struct {
	Number          string
	AccountNumber   string
	BankName        string
	ChequeDate      string
	TransactionDate string
	DepositedDate   string
	Comment         string
}
type Online struct {
	ServiceName       string
	TransactionDate   string
	TransactionNumber string
	Comment           string
}
type Cash struct {
	TransactionDate string
}

type FormData struct {
	SlNum             string
	Name              string
	RoleIDs           []guid.Guid
	MainInfo          Info
	RentalsMasterView RentalInfo
	RentalItems       []RentalItemv2
}

// type FormData struct {
// 	SlNum             string
// 	Name              string
// 	MainInfo          Info
// 	RentalsMasterView RentalInfo
// 	RentalItems       []RentalItemv2
// }
type Tenants struct {
	Tenant RentalInfo
	Name   string
	SlNum  string
}

const maxMEMORY = 1 * 1024 * 1024

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func commonHttpServReqHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(w)
	if r.Method == "OPTIONS" {
		return
	}

	//********
	path := r.URL.String()

	table, key, err := requestParser(path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if strings.ToUpper(table) == "SIGNIN" {
		result, err := Signin(w, r)
		if err != nil {
			return
		}
		w.Write(result)
		return
	} else if strings.ToUpper(table) == "AUTHTOKEN" {
		_, err := AuthenticateToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	/***************************************************************************/
	/* session cookie checking
	/*******************************************************************************/
	// We can obtain the session token from the requests cookies, which come with every request

	claimID, err := AuthenticateToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// after check get data according to authenticated user

	switch r.Method {
	case "GET":
		if result, err2 := getData(table, key, claimID); err2 == nil {
			if table != "image" {
				w.Write(result)
			} else {
				http.ServeFile(w, r, string(result))
			}
		} else {
			switch err2.Error() {
			case "TableNotFound":
				fallthrough
			case "FileNotFound":
				http.Error(w, err2.Error(), http.StatusNotFound)
			}
		}
	case "POST":

		if err := r.ParseMultipartForm(maxMEMORY); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				file, _ := fileHeader.Open()
				path := fmt.Sprintf("files/%s", fileHeader.Filename)
				buf, _ := ioutil.ReadAll(file)
				err := ioutil.WriteFile(path, buf, os.ModePerm)
				if err != nil {
					println("error", err.Error())
				}
			}
		}
		if result, err2 := postData(table, claimID, r.MultipartForm.Value); err2 == nil {
			_, err := w.Write(result)
			if err != nil {
				println(err.Error())
			}
		} else {
			http.Error(w, err2.Error(), http.StatusNotFound)
		}

		//Update
	case "PUT":
		if err := r.ParseMultipartForm(maxMEMORY); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if result, err2 := updateData(table, claimID, r.MultipartForm.Value); err2 == nil {
			_, err := w.Write(result)
			if err != nil {
				println(err.Error())
			}
		} else {
			http.Error(w, err2.Error(), http.StatusNotFound)
		}
		//Edit
	case "PATCH":
		if err := r.ParseMultipartForm(maxMEMORY); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if result, err2 := editData(table, claimID, r.MultipartForm.Value); err2 == nil {
			_, err := w.Write(result)
			if err != nil {
				println(err.Error())
			}
		} else {
			http.Error(w, err2.Error(), http.StatusNotFound)
		}
	case "DELETE":
		var u Id
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
		}
		data, err := url.QueryUnescape(key)
		if err != nil {
			println(err.Error())
		}
		json.Unmarshal([]byte(data), &u)
		deleteData(table, claimID, u)
	}

}

func requestParser(req_path string) (table string, key string, err error) {
	//NOTE: strings.Split returns empty string hence it is testing less than 3.
	req_path_tokenized := strings.Split(req_path, "/")
	println(req_path)

	//NOTE: home path is "/" which splits into array of len 2.
	if len(req_path_tokenized) == 2 && req_path_tokenized[1] == "" {
		table = "HOMEPAGE"
		return
	}

	if len(req_path_tokenized) < 3 {
		if len(req_path_tokenized) == 2 && req_path_tokenized[1] != "" {
			//only table name is given.
			table = req_path_tokenized[1]
			return
		}

		err = errors.New("BadRequest")
		return
	}

	table = req_path_tokenized[1]
	key = req_path_tokenized[2]
	return
}
