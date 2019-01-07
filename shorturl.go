package main

import (
    "fmt"
    "errors"
    "time"
    "net/url"
    "database/sql"
    "math/rand"
    _ "github.com/bmizerany/pq"
)

const (
  host     = "host"
  user     = "username"
  password = "password"
  dbname   = "database"
)

func main() {
    db, err := SetupConnection()
    if err != nil {
        fmt.Println(err)
    }

    code, record, err := StoreRecord(db, "https://github.com/aavrug")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Print("Successfully ",record," insertions happened, the code is http://eg.c/",code,"\n")
}

func SetupConnection() (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s user=%s "+
    "password=%s dbname=%s sslmode=disable", host, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func StoreRecord(db *sql.DB, urlString string) (string, int64, error) {
    if urlString == "" {
        return "", 0, errors.New("URL can't be empty!")
    }

    _, err := url.ParseRequestURI(urlString)
    if err != nil {
        return "", 0, errors.New("URL is not valid!")
    }

    created_time := time.Now()
    code := GetRandString(10)
    res, err := db.Exec("INSERT INTO records(code, url, created_time) VALUES"+
    "($1, $2, $3)", code, urlString, created_time.Unix())

    if err != nil {
        return "", 0, err
    }

    record, err := res.RowsAffected()

    if err != nil {
        return "", 0, err
    }

    return code, record, nil
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func GetRandString(n int) string {
    var randomString = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    str := make([]rune, n)
    for i := range str {
        str[i] = randomString[rand.Intn(len(randomString))]
    }
    return string(str)
}

func GetRecord(db *sql.DB, code string) (string, error) {
    if code == "" {
        return "", errors.New("Short code can't be empty!")
    }

    var url string
    err := db.QueryRow("SELECT url FROM records WHERE code = $1", code).Scan(&url)

    if err != nil {
        return "", err
    }

    return url, nil
}

func GetAllRecords(db *sql.DB) (map[string]string, error) {
    rows, err := db.Query("SELECT id, code FROM records")

    if err != nil {
        return nil, err
    }
    var (
        id string
        code string
    )

    records := make(map[string]string)
    for rows.Next() {
        err := rows.Scan(&id, &code)
        if err != nil {
            return nil, err
        }
        records[id] = code
    }

    return records, nil
}
