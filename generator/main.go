package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type FibonacciGenerator struct {
	a, b int
}

func NewFibonacciGenerator() *FibonacciGenerator {
	return &FibonacciGenerator{a: 0, b: 1}
}

func (fg *FibonacciGenerator) Next() int {
	fg.a, fg.b = fg.b, fg.a+fg.b
	return fg.a
}

type Payload struct {
	Value int `json:"value"`
}

func sendFibonacci(url string, value int) error {
	payload := Payload{Value: value}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Protocol", "HTTP/2.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %v", resp.Status)
	}

	return nil
}

func storeFibonacci(db *sql.DB, value int) error {
	_, err := db.Exec("INSERT INTO fibonacci (value) VALUES (?)", value)
	return err
}

func getLastFibonacci(db *sql.DB) (int, error) {
	var value int
	err := db.QueryRow("SELECT value FROM fibonacci ORDER BY id DESC LIMIT 1").Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return value, nil
}

func main() {
	db, err := sql.Open("sqlite", "./fibonacci.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS fibonacci (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value INTEGER
	)`)
	if err != nil {
		log.Fatal(err)
	}

	lastValue, err := getLastFibonacci(db)
	if err != nil {
		log.Fatal(err)
	}

	var fib *FibonacciGenerator
	if lastValue == 0 {
		fib = NewFibonacciGenerator()
	} else {
		fib = &FibonacciGenerator{}
		a, b := 0, 1
		for a < lastValue {
			a, b = b, a+b
			fib.a, fib.b = a, b
		}
	}

	url := "http://server:8080/fibonacci"

	for {
		value := fib.Next()

		fmt.Println(value)

		err := sendFibonacci(url, value)
		if err != nil {
			log.Printf("failed send to server: %v", err)
		}

		err = storeFibonacci(db, value)
		if err != nil {
			log.Printf("failed save Fibonacci number: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
