package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Response struct {
	Drivers []Driver `json:"drivers"`
}

type Driver struct {
	Id          string `json:"id"`
	DriverName  string `json:"name"`
	Nationality string `json:"nationality"`
	Team        string `json:"team"`
	Points      string `json:"points"`
}

func main() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	//router.HandleFunc("/persons", Persons).Methods("GET")
	router.HandleFunc("/drivers/{id}", Drivers).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("entering health check end point")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func Drivers(w http.ResponseWriter, r *http.Request) {
	log.Println("entering drivers end point")
	var response Response
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)

	pilots := prepareResponse(string(id))

	response.Drivers = pilots

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}

func prepareResponse(searchId string) []Driver {
	//searchId := "3"
	csvFile, _ := os.Open("resources/F1_Drivers.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var pilot []Driver
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if line[0] == searchId {
			pilot = append(pilot, Driver{
				Id:          line[0],
				DriverName:  line[1],
				Nationality: line[2],
				Team:        line[3],
				Points:      line[4],
			})
		}
	}
	return pilot
}
