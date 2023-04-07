package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gookit/color"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	logFatal(err)

	defer conn.Close()
	color.Cyan.Printf("Enter your name: ")

	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	logFatal(err)

	username = strings.Trim(username, " \r\n")
	welcomeMsg := color.Info.Sprintf("welcome to the chat %s\n", color.Red.Sprintf(username))

	fmt.Println(welcomeMsg)

	go read(conn)
	write(conn, username)
}

func read(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		msg = color.White.Sprint(msg)
		if err == io.EOF {
			conn.Close()
			color.Error.Println("__server closed__")
			os.Exit(0)
		}

		fmt.Println(msg)
	}
}

func write(conn net.Conn, usr string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		msg = fmt.Sprintf("%s: %s\n", usr, strings.Trim(msg, " \r\n"))

		conn.Write([]byte(msg + "\n"))
	}
}
