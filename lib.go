package GolangBsonDB

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
)

type BsonDB struct {
  databaseId string
}

func NewBsonDB(databaseId string) *BsonDB {
  return &BsonDB{databaseId: databaseId}
}

func (db *BsonDB) CreateEntry(tableName string, entry interface{}) ([]byte, error) {
  body := map[string]interface{}{
    "databaseId": db.databaseId,
    "table":      tableName,
    "entry":      entry,
  }
  return apiRequest("POST", "/api/add-entry", body)
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


