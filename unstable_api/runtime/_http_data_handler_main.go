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

package runtime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/beevik/guid"
	"google.golang.org/api/calendar/v3"
)

type UpdateInfo struct {
	UpdateField string
	UpdateItem  string
}
type EditInfo struct {
	Property string
	Item     string
}
type Details struct {
	Info interface{}
}

type Account struct {
	Key    string
	Pwd    string
	Name   string
	Claims []guid.Guid
}
type Role struct {
	Role  string
	Claim guid.Guid
}
type FileData struct {
	Design string
	Script string
}
type File struct {
	Id          string
	Name        string
	Description string
	Data        FileData
	Created     string
	Modified    string
}

var data []FormData
var srv *calendar.Service
var Accounts = map[string]Account{}
var Roles = map[string]Role{}
var imageFiles = map[string]string{}
var propertiesFile = "unstable_api/properties.json"
var templatesFile = "unstable_api/templates.json"

func init() {
	if err := loadAccounts("unstable_api/accounts.json"); err != nil {
		println("Counld not load ACCOUNTS")
	}
	if err := loadRoles("unstable_api/roles.json"); err != nil {
		println("could not load Roles")
	}
	srv = getServiece()
	//image files init
	imageFiles["ASSETS1"] = "./files/assets1.jpg"
}

func getData(table, key string, claimID []guid.Guid) ([]byte, error) {

	//NOTE: add case here for another table call.
	switch strings.ToUpper(table) {
	case "HOMEPAGE":
		return homePage(), nil

	// case "ACCOUNTS":
	// 	if key != "" {
	// 		if data, ok := Accounts[strings.ToUpper(key)]; ok {
	// 			if b, err := json.Marshal(data); err == nil {
	// 				return b, nil
	// 			}
	// 		}
	// 	}
	// 	return nil, tableDoesNotSupportListOperation(table)

	case "PROPERTIES":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			if b, err := json.Marshal(data); err == nil {
				return b, nil
			}
		} else {
			println(err.Error())
		}
		return nil, tableDoesNotSupportListOperation(table)
	case "TENANTS":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			tenant := extractTenants(data)
			if b, err := json.Marshal(tenant); err == nil {
				return b, nil
			}
		}
		return nil, tableDoesNotSupportListOperation(table)
	case "EVENTS":
		events, err := getevents(srv)
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
			return nil, err
		}
		if b, err := json.Marshal(events.Items); err == nil {
			return b, nil
		}
		return nil, err
	case "TEMPLATES":
		if templates, err := loadTemplates(templatesFile); err == nil {
			if b, err := json.Marshal(templates); err == nil {
				return b, nil
			}
		} else {
			println(err.Error())
		}
		return nil, tableDoesNotSupportListOperation(table)
	case "IMAGE":

		if key != "" {

			if imageFileName, ok := imageFiles[strings.ToUpper(key)]; ok {
				return []byte(imageFileName), nil
			}

			return nil, errors.New("FileNotFound")
		}

		return nil, tableDoesNotSupportListOperation(table)
	}

	return nil, errors.New("TableNotFound")
}

func homePage() []byte {
	return []byte("Hi there!")
}

func tableDoesNotSupportListOperation(table string) error {
	return errors.New(fmt.Sprintf("Table '%s' does not support list operation.", table))
}

func postData(table string, claimID []guid.Guid, form url.Values) ([]byte, error) {
	println(table)
	switch strings.ToUpper(table) {
	case "PROPERTIES":
		data1, _ := getProperties(form)
		// data1.RoleIDs = append(data1.RoleIDs, Roles["ADMIN"].RoleID)
		// TO DO
		// check claims before add,
		for _, Id := range claimID {
			data1.RoleIDs = append(data1.RoleIDs, Id)
		}
		if err := writeFile(propertiesFile, data1); err != nil {
			return nil, err
		}
		return []byte("done"), nil

	case "EVENTS":
		action := form.Get("action")
		switch action {
		case "create":
			event := getEvents(form)
			evt, err := srv.Events.Insert("primary", &event).Do()
			if err != nil {
				return nil, err
			}
			bytes, _ := evt.MarshalJSON()
			return bytes, nil
		case "edit":
			event := getEvents(form)
			evt, err := srv.Events.Update("primary", event.Id, &event).Do()
			if err != nil {
				return nil, err
			}
			bytes, _ := evt.MarshalJSON()
			return bytes, nil
		case "delete":
			eventID := form.Get("eventId")
			err := srv.Events.Delete("primary", eventID).Do()
			if err != nil {
				return nil, err
			}
			return []byte("deleted event"), nil
		}
	case "TEMPLATES":
		file := parseFile(form)
		if err := writeTemplate(templatesFile, file); err != nil {
			return nil, err
		}
		return []byte("done"), nil

	}
	return nil, errors.New("TableNotFound")
}

func updateData(table string, claimID []guid.Guid, form url.Values) ([]byte, error) {
	switch strings.ToUpper(table) {
	case "PROPERTIES":
		var updateInfo UpdateInfo
		newData, Info := getProperties(form)
		err := json.Unmarshal(Info, &updateInfo)
		if err != nil {
			return nil, err
		}
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for i, property := range data {
				if property.SlNum == updateInfo.UpdateItem {
					if updateInfo.UpdateField == "tenants" {
						data[i].RentalsMasterView = newData.RentalsMasterView
						data[i].RentalItems = newData.RentalItems
					} else if updateInfo.UpdateField == "rentals" {
						for _, rentalItem := range newData.RentalItems {
							data[i].RentalItems = append(data[i].RentalItems, rentalItem)
						}
					}
					break
				}
			}
			if err := updateFile(propertiesFile, data); err != nil {
				return nil, err
			}
			return []byte("updated property details"), nil
		}
	case "TEMPLATES":
		file := parseFile(form)
		if templates, err := loadTemplates(templatesFile); err == nil {
			for i, temp := range templates {
				if temp.Id == file.Id {
					templates[i] = file
					break
				}
			}
			if err := updateTemplate(templatesFile, templates); err != nil {
				return nil, err
			}
			return []byte("updaed templates"), nil
		}
	}
	return nil, errors.New("TableNotFound")
}

func editData(table string, claimID []guid.Guid, form url.Values) ([]byte, error) {
	var editInfo EditInfo
	newData, Info := getProperties(form)
	println(string(Info))
	err := json.Unmarshal(Info, &editInfo)
	if err != nil {
		return nil, err
	}
	switch strings.ToUpper(table) {
	case "RENTALS":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for p, property := range data {
				println(property.SlNum, "==", editInfo.Property)
				if property.SlNum == editInfo.Property {
					for i, rentalItem := range property.RentalItems {
						if rentalItem.Id == editInfo.Item {
							// replace this item with editted Item
							data[p].RentalItems[i] = newData.RentalItems[0]
						}
					}
					if err := updateFile(propertiesFile, data); err != nil {
						return nil, err
					}
					return []byte("saved property details"), nil
				}
			}
		}
	case "MAININFO":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for p, property := range data {
				if property.SlNum == editInfo.Property {
					data[p].MainInfo = newData.MainInfo
					if err := updateFile(propertiesFile, data); err != nil {
						return nil, err
					}
					return []byte("saved property details"), nil
				}
			}
		}
	case "TENANTINFO":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for p, property := range data {
				if property.SlNum == editInfo.Property {
					data[p].RentalsMasterView = newData.RentalsMasterView
					if err := updateFile(propertiesFile, data); err != nil {
						return nil, err
					}
					return []byte("saved property details"), nil
				}
			}
		}
	}
	return nil, errors.New("TableNotFound")
}

func deleteData(table string, claimID []guid.Guid, details Id) ([]byte, error) {
	switch strings.ToUpper(table) {
	case "RENTALS":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for p, property := range data {
				if property.SlNum == details.Property {
					for i, rentalItem := range property.RentalItems {
						if rentalItem.Id == details.Item {
							// delete this item with editted Item
							data[p].RentalItems = append(data[p].RentalItems[:i], data[p].RentalItems[i+1:]...)
						}
					}
					if err := updateFile(propertiesFile, data); err != nil {
						return nil, err
					}
					return []byte("saved property details"), nil
				}
			}
		}
	case "TENANT":
		var temp RentalInfo
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for p, property := range data {
				if property.SlNum == details.Property {
					data[p].RentalsMasterView = temp
				}
			}
		}
	case "PROPERTY":
		if data, err := loadProperties(propertiesFile, claimID); err == nil {
			for p, property := range data {
				if property.SlNum == details.Property {
					data = append(data[:p], data[p+1:]...)
				}
			}
		}
	}
	return nil, errors.New("TableNotFound")
}

func updateFile(dbfile string, data []FormData) error {
	db, err := os.OpenFile(dbfile, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer db.Close()
	enc := json.NewEncoder(db)
	for _, update := range data {
		err := enc.Encode(update)
		if err != nil {
			return err
		}
	}
	return nil
}
func updateTemplate(dbfile string, data []File) error {
	db, err := os.OpenFile(dbfile, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer db.Close()
	enc := json.NewEncoder(db)
	for _, update := range data {
		err := enc.Encode(update)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(dbfile string, data FormData) error {
	db, err := os.OpenFile(dbfile, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		db, err = os.OpenFile(dbfile, os.O_APPEND, 0644)
	}
	defer db.Close()
	enc := json.NewEncoder(db)
	err = enc.Encode(&data)
	if err != nil {
		return err
	}
	return nil
}
func writeTemplate(dbfile string, file File) error {
	db, err := os.OpenFile(dbfile, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		db, err = os.OpenFile(dbfile, os.O_APPEND, 0644)
	}
	defer db.Close()
	enc := json.NewEncoder(db)
	err = enc.Encode(&file)
	if err != nil {
		return err
	}
	return nil
}

func loadProperties(dbFile string, claimID []guid.Guid) ([]FormData, error) {
	var lastBlock FormData

	db, err := os.Open(dbFile)
	if err != nil {
		return nil, fmt.Errorf("no existing database found")
	}
	defer db.Close()
	data := []FormData{}
	dec := json.NewDecoder(db)
	for i := 0; ; i++ {
		lastBlock = FormData{}
		err := dec.Decode(&lastBlock)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		brk := false
		for _, id := range lastBlock.RoleIDs {
			for _, claimId := range claimID {
				if claimId == id {
					data = append(data, lastBlock)
					brk = true
					break
				}
			}
			if brk {
				break
			}
		}
	}
	return data, nil
}
func loadTemplates(dbFile string) ([]File, error) {
	var lastBlock File

	db, err := os.Open(dbFile)
	if err != nil {
		return nil, fmt.Errorf("no existing database found")
	}
	defer db.Close()
	data := []File{}
	dec := json.NewDecoder(db)
	for i := 0; ; i++ {
		lastBlock = File{}
		err := dec.Decode(&lastBlock)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data = append(data, lastBlock)
	}
	return data, nil
}
func extractTenants(data []FormData) (tenants []Tenants) {
	for _, value := range data {
		tenant := Tenants{}
		if value.RentalsMasterView.TenantInfo.TenantId == "" {
			continue
		}
		tenant.Name = value.Name
		tenant.SlNum = value.SlNum
		tenant.Tenant = value.RentalsMasterView
		tenants = append(tenants, tenant)
	}
	return
}

func getEvents(form url.Values) calendar.Event {
	event := calendar.Event{}
	event.Start = &calendar.EventDateTime{}
	event.End = &calendar.EventDateTime{}
	for key, value := range form {
		if key == "action" {
			continue
		}
		formValues, err := url.ParseQuery(value[0])
		if err != nil {
			log.Fatal(err)
		}
		switch key {
		case "event":
			for key1, value1 := range formValues {
				switch key1 {
				case "id":
					event.Id = value1[0]
				case "summary":
					event.Summary = value1[0]
				case "start[date]":
					if formValues.Get("start[time]") != "" {
						break
					}
					event.Start.Date = value1[0]
				case "start[dateTime]":
					event.Start.DateTime = formValues.Get("start[dateTime]")
				case "location":
					event.Location = value1[0]
				case "end[date]":
					if formValues.Get("end[time]") != "" {
						break
					}
					event.End.Date = value1[0]
				case "end[time]":
					event.End.DateTime = formValues.Get("end[dateTime]")
				case "description":
					event.Description = formValues.Get("description")
				}
			}
		}
	}
	return event
}
func getProperties(form url.Values) (data FormData, info []byte) {
	for key := range form {
		switch key {
		case "form":
			err := json.Unmarshal([]byte(form.Get(key)), &data)
			if err != nil {
				println(err.Error())
			}
		case "updateInfo", "editInfo":
			var details Details
			err := json.Unmarshal([]byte(form.Get(key)), &details.Info)
			if err != nil {
				println(err.Error())
			}
			info, _ = json.Marshal(details.Info)
		}
	}
	return
}
func parseFile(form url.Values) (file File) {
	var filedata FileData
	var name string
	var description string
	var created string
	var modified string
	var id string
	for key := range form {
		switch key {
		case "data":
			err := json.Unmarshal([]byte(form.Get(key)), &filedata)
			if err != nil {
				println(err.Error())
			}
		case "name":
			err := json.Unmarshal([]byte(form.Get(key)), &name)
			if err != nil {
				println(err.Error())
			}
		case "description":
			err := json.Unmarshal([]byte(form.Get(key)), &description)
			if err != nil {
				println(err.Error())
			}
		case "created":
			err := json.Unmarshal([]byte(form.Get(key)), &created)
			if err != nil {
				println(err.Error())
			}
		case "modified":
			err := json.Unmarshal([]byte(form.Get(key)), &modified)
			if err != nil {
				println(err.Error())
			}
		case "id":
			err := json.Unmarshal([]byte(form.Get(key)), &id)
			if err != nil {
				println("err", err.Error())
			}
		}
	}
	file.Name = name
	file.Description = description
	file.Data = filedata
	file.Created = created
	file.Modified = modified
	file.Id = id
	return
}
func prepareMultipartForm(params, files map[string]string) (*bytes.Buffer, error) {
	formBuffer := new(bytes.Buffer)
	writer := multipart.NewWriter(formBuffer)
	for fileName, path := range files {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}
		file.Close()

		part, err := writer.CreateFormFile(fileName, fi.Name())
		if err != nil {
			return nil, err
		}
		part.Write(fileContents)
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return formBuffer, nil
}
func loadAccounts(dbFile string) (err error) {
	db, err := os.Open(dbFile)
	if err != nil {
		return fmt.Errorf("no existing database found")
	}
	defer db.Close()
	dec := json.NewDecoder(db)
	for i := 0; ; i++ {
		account := Account{}
		err := dec.Decode(&account)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		Accounts[strings.ToUpper(account.Key)] = account
	}

	return nil
}
func loadRoles(dbFile string) (err error) {
	db, err := os.Open(dbFile)
	if err != nil {
		return fmt.Errorf("no existing database found")
	}
	defer db.Close()
	dec := json.NewDecoder(db)
	for i := 0; ; i++ {
		role := Role{}
		err := dec.Decode(&role)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		Roles[strings.ToUpper(role.Role)] = role
	}

	return nil
}
func setDateTime(date, time string) string {
	var buffer bytes.Buffer
	buffer.WriteString(date)
	buffer.WriteString("T")
	buffer.WriteString(time)
	buffer.WriteString("+05:30")

	return buffer.String()
}
