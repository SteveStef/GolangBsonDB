package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Table struct {
	Identifier   string                 `json:"identifier"`
	Name         string                 `json:"name"`
	EntryTemplate map[string]interface{} `json:"entryTemplate"`
	Requires     []string               `json:"requires"`
}

type Database struct {
	DatabaseID string   `json:"DatabaseID"`
	Tables     []Table `json:"Tables"`
}

const server = "https://bsondb.up.railway.app"

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
}

func apirequest(method, path string, body interface{}) (map[string]interface{}, error) {
	url := server + path
	var response map[string]interface{}

	client := &http.Client{}
	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, strings.NewReader(string(jsonBody)))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	for key, value := range defaultHeaders {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func createTable(databaseID string, tables []Table) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"databaseId": databaseID,
		"tables":     tables,
	}
	return apirequest("POST", "/api/migrate-tables", body)
}

func main() {

  file, err := os.Open("tables.json")
  if err != nil {
    fmt.Println("Error opening tables.json:", err)
    return
  }

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to continue? All table data will be lost. (yes/no): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "yes" || answer == "y" {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		var db Database
		if err := json.Unmarshal(data, &db); err != nil {
			fmt.Println(err)
			return
		}

		validTypes := []string{"string", "number", "boolean", "object"}
		for _, table := range db.Tables {
			if table.Identifier == "" || table.Name == "" || table.EntryTemplate == nil || len(table.Requires) == 0 {
				fmt.Println("Invalid table object. Table object must have identifier, name, entryTemplate, and requires properties.")
				return
			}

			for key, value := range table.EntryTemplate {
				if !contains(validTypes, fmt.Sprintf("%T", value)) {
					fmt.Printf("Invalid type for %s in %s. Valid types are: string, number, boolean, object\n", key, table.Name)
					return
				}
			}

			if table.EntryTemplate[table.Identifier] != "string" {
				fmt.Println("Identifier must be of type string.")
				return
			}
		}

		response, err := createTable(db.DatabaseID, db.Tables)
		if err != nil {
			fmt.Println(err)
			return
		}

		if response["error"] != nil {
			fmt.Println("Error creating tables:", response["error"])
			return
		}

		fmt.Printf("%d tables have been created successfully.\n", len(db.Tables))
		fmt.Println("Go to https://bson-api.com/ to view your database.")
	} else {
		fmt.Println("The migration has been canceled.")
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
