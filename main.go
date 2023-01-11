package main

import "github.com/akshaym-3255/mygrep/cmd"

func main() {
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	// if scanner.Err() != nil {

	// 	// Handle error.
	// }
	cmd.Execute()
}

// package main

// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"log"
// 	"os"
// )

// func parse_args() (file, pat string) {
// 	if len(os.Args) < 3 {
// 		log.Fatal("usage: petergrep <file_name> <pattern>")
// 	}
// 	file = os.Args[1]
// 	pat = os.Args[2]
// 	return
// }

// func grepFile(file string, pat []byte) int64 {
// 	patCount := int64(0)
// 	f, err := os.Open(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		if bytes.Contains(scanner.Bytes(), pat) {
// 			patCount++
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 	}
// 	return patCount
// }

// func main() {
// 	file, pat := parse_args()
// 	total := grepFile(file, []byte(pat))
// 	fmt.Printf("Total %d\n", total)
// }
