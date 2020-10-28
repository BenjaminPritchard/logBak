// Example program for working with Google embedded NoSQL engine LevelDB
// from a GoLang program.
//
// This simple utility can convert back and forth between a plain text
// log file, and a LevelDB database file.
//
// Usage:
//	logBak b log_file NoSQL_file		-> backups log_file to NoSQL_file
//	logBak r log_file NoSQL_file		-> restores log_file from NoSQL_file
//
// Benjamin Pritchard

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func main() {

	// make sure command line is OK
	cmdLineOK := len(os.Args) == 4
	if cmdLineOK {
		cmdLineOK = os.Args[1] == "b" || os.Args[1] == "r"
	}

	if !cmdLineOK {
		fmt.Println("Usage:", "logBak b log_file NoSQL_file - backups log_file to NoSQL_file")
		fmt.Println("Usage:", "logBak r log_file NoSQL_file - restores log_file from NoSQL_file")
		os.Exit(1)
	}

	operation := os.Args[1]
	log_file := os.Args[2]
	NoSQL_file := os.Args[3]

	db, err := leveldb.OpenFile(NoSQL_file, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// backup
	if operation == "b" {
		fmt.Printf("backup %s to %s\n", log_file, NoSQL_file)

		f, err := os.Open(log_file)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		// read in each line of the log file, and append a number in front of it
		// eg: xxx becomes 0000|xxxx
		// then append the appended line into the key|value database

		scanner := bufio.NewScanner(f)
		counter := 0
		for scanner.Scan() {
			x := lpad(strconv.Itoa(counter), "0", 5) + "|" + scanner.Text()
			err = db.Put([]byte(x), []byte(x), nil)
			counter++
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("restore %s from %s\n", log_file, NoSQL_file)

		f, err := os.Create(log_file)
		if err != nil {
			fmt.Println(err)
			return
		}

		// loop through all the keys
		var values []string
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			value := iter.Value()
			values = append(values, string(value))
		}

		// the keys will be in a random order
		// so sort them based on the number we put in front of each one
		sort.Strings(values)

		// now write each value back out to disk...
		// remembering to strip off the number and the |
		for i, s := range values {
			_ = i
			l, err := f.WriteString(strings.Split(string(s), "|")[1] + "\n")
			if l == 0 {
				log.Fatal("0 bytes written??")
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		iter.Release()
	}

}
