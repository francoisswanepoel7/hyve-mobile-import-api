package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	rawcsvMap         = map[string]*CSVDataValidated{}
	rawcsvMapMapMutex = sync.RWMutex{}

	emailValidationMap         = map[string]bool{}
	emailValidationMapMapMutex = sync.RWMutex{}

	domainIpMap            = map[string]string{}
	domainIpMapMapMapMutex = sync.RWMutex{}
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hyve Mobile Importer")
}

func processCSV(w http.ResponseWriter, r *http.Request) {
	var ok OK
	ok.Result = "ok"

	// Put CSV Data into Struct Array
	decoder := json.NewDecoder(r.Body)
	var csv_data []*CSVDataValidated
	err := decoder.Decode(&csv_data)
	if err != nil {
		panic(err)
	}

	// Sync Email -> CSV Data Struct
	for i := 0; i < len(csv_data); i++ {
		var email = csv_data[i].Email
		if _, ok := rawcsvMap[email]; ok {
			fmt.Println("Email: " + email + " exists")
		} else {
			fmt.Println("New Email: " + email + "")
			rawcsvMapMapMutex.Lock()
			rawcsvMap[email] = csv_data[i]
			rawcsvMapMapMutex.Unlock()
		}

	}

	// Send data back to client
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ok)

	go func() {
		for email, _ := range rawcsvMap {

			if _, email_ok := emailValidationMap[email]; email_ok {
			} else {
				emailValidationMapMapMutex.Lock()
				emailValidationMap[email] = isEmailValid(email)
				emailValidationMapMapMutex.Unlock()
			}
		}
	}()

	go func() {
		for email, _ := range rawcsvMap {
			domain := getDomainFromEmail(email)
			if _, domain_ok := domainIpMap[domain]; domain_ok {
			} else {
				domainIpMapMapMapMutex.Lock()
				domainIpMap[domain] = ipLookup(domain)
				domainIpMapMapMapMutex.Unlock()
			}
		}
	}()
}

func importCSV(w http.ResponseWriter, r *http.Request) {
	var ok OK
	ok.Result = "ok"
	// Send data back to client
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ok)
	go func() {
		for email, csv_data := range rawcsvMap {
			domain := getDomainFromEmail(email)
			ip := domainIpMap[domain]
			valid_email := emailValidationMap[email]
			csv_data.IP = ip
			if ip != "" && ip != "127.0.0.1" && valid_email {

				if csv_data.Processed != true {
					if insertContact(csv_data) {
						tzId := insertTimeZone(csv_data.Tz)
						insertContactTimeZone(csv_data, tzId)
						insertCompleted(csv_data)
					}
					csv_data.Processed = true

				}
			}
		}
	}()
	fmt.Println("Processed")
}

func exportData(w http.ResponseWriter, r *http.Request) {
	data := getProcessed()
	d, _ := json.Marshal(data)
	//fmt.Println(d)
	resp, err := http.Post(os.Getenv("API_ENDPOINT"), "application/json", bytes.NewBuffer(d))

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

}

func init() {
	rawcsvMap = make(map[string]*CSVDataValidated)
	emailValidationMap = make(map[string]bool)
	domainIpMap = make(map[string]string)

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/process", processCSV).Methods("POST")
	router.HandleFunc("/import", importCSV).Methods("GET")
	router.HandleFunc("/export", exportData).Methods("GET")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_URI"), router))

}
