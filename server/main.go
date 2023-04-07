package main

import (
	"context"
	"fmt"
	"log"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Println("press enter to close the server...")
		fmt.Scanln()
		cancel()
	}()

	server := newServer(":8080")

	err := server.Start(ctx, cancel)
	logFatal(err)

	// text := make(chan string, 1)

	// value := read()

	// fs, err := os.Create("file.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer fs.Close()

	// // b := []byte("Hola kjajaja Hola kjajajaHola")
	// _, err = fs.Write(value)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// text <- string(value)
	// fmt.Println("Wrote byte: ", <-text)
}

// func read() []byte {
// 	fs, err := os.Open("readme.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer fs.Close()

// 	// create a buffer to hold the bytes read from the file
// 	var buffer bytes.Buffer

// 	b := make([]byte, 1024)

// 	for {
// 		n, err := fs.Read(b)
// 		if err != nil && err != io.EOF {
// 			log.Fatal(err)
// 		}

// 		if n == 0 {
// 			break
// 		}

// 		buffer.Write(b[:n])
// 	}

// 	fmt.Printf("Read byte: %v\n", b)

// 	return buffer.Bytes()
// }
