// file will contain all used public used SQL functions within the API Private SQL functions will be hosted wihtin Pirvate Repo. 

package SQLFunctions 

import (
  "fmt"
  "log"
  "database/sql"
)


func connectToDB(){
  const (
    host = "localhost"
    port = 5432
    user = "postgres"
    password = "SugaBear2025"
    dbname = "arcRaidersAPI"
  )
  
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s" + "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    log.Fatalf("Error opening the database: %q", err)
  }

  defer db.Close()

  err = db.Ping()
  if err != nil{
    log.Fatalf("Error connecting to the data base: %q", err)
  }

  fmt.Printf("Connected to the DB!!!!!!")
}
