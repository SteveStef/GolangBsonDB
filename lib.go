package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
)

func init() {
  file, err := os.Create("tables.json")
  if err != nil {
    fmt.Println(err)
  }
  defer file.Close()
  file.WriteString(`
{
  "DatabaseID": "Put database ID here",
  "Tables": [
    {
      "name": "users",
      "identifier": "email",
      "requires": ["email", "password"],
      "entryTemplate": {
        "email": "string",
        "password": "string"
      }
    }
  ]
}`)
}

type BsonDB struct {
  databaseId string
}

func NewBsonDB(databaseId string) *BsonDB {
  return &BsonDB{databaseId: databaseId}
}

func apiRequest(method, path string, body interface{}) ([]byte, error) {
  url := "https://bsondb.up.railway.app" + path
  bodyBytes, err := json.Marshal(body)
  if err != nil {
    return nil, err
  }

  req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
  if err != nil {
    return nil, err
  }
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode < 200 || resp.StatusCode > 299 {
    return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
  }

  return ioutil.ReadAll(resp.Body)
}

func (db *BsonDB) CreateEntry(tableName string, entry interface{}) ([]byte, error) {
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      tableName,
    "entry":      entry,
  }
  return apiRequest("POST", "/api/add-entry", body)
}

func (db *BsonDB) UpdateEntry(table string, query map[string]interface{}) ([]byte, error) {
  if query["where"] == nil || query["set"] == nil {
    fmt.Println("Invalid query, use the form {where: 'identifier', set: object}")
    return nil, fmt.Errorf("invalid query")
  }
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      table,
    "entryId":    query["where"],
    "entry":      query["set"],
  }
  return apiRequest("PUT", "/api/update-field", body)
}

func (db *BsonDB) DeleteEntry(table string, query map[string]interface{}) ([]byte, error) {
  if query["where"] == nil {
    fmt.Println("Invalid query, use the form {where: 'identifier'}")
    return nil, fmt.Errorf("invalid query")
  }
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      table,
    "entryId":    query["where"],
  }
  return apiRequest("POST", "/api/delete-entry", body)
}

func (db *BsonDB) GetTable(table string) ([]byte, error) {
  body := map[string]interface{} {
    "databaseId": db.databaseId,
    "table":      table,
  }
  return apiRequest("POST", fmt.Sprintf("/api/table"), body)
}

func (db *BsonDB) GetEntry(tableName string, query map[string]interface{}) ([]byte, error) {
  if query["where"] == nil {
    fmt.Println("Invalid query, use the form {where: 'identifier'}")
    return nil, fmt.Errorf("invalid query")
  }
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      tableName,
    "entryId":    query["where"],
  }
  return apiRequest("POST", "/api/entry", body)
}

func (db *BsonDB) GetField(tableName string, query map[string]interface{}) ([]byte, error) {
  if query["where"] == nil || query["get"] == nil {
    fmt.Println("Invalid query, use the form {where: 'identifier', get: 'field-name'}")
    return nil, fmt.Errorf("invalid query")
  }
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      tableName,
    "entryId":    query["where"],
    "field":      query["get"],
  }
  return apiRequest("POST", "/api/field", body)
}

func (db *BsonDB) GetEntries(tableName string, query map[string]interface{}) ([]byte, error) {
  if query["where"] == nil || query["is"] == nil {
    fmt.Println("Invalid query, use the form {where: 'field', is: 'value'}")
    return nil, fmt.Errorf("invalid query")
  }
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      tableName,
    "field":      query["where"],
    "value":      query["is"],
  }
  return apiRequest("POST", "/api/entries", body)
}
