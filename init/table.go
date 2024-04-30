package main
import (
  "fmt"
  "os"
)

func main() {
  file, err := os.Create("../tables.json")
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
