package main

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/gookit/color"
)

type Server struct {
	ln         net.Listener
	listenAddr string
	messages   chan string
	clients    map[net.Conn]bool
	client     chan net.Conn
	deadConn   chan net.Conn
	errors     chan error
}

func newServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		messages:   make(chan string),
		clients:    make(map[net.Conn]bool),
		client:     make(chan net.Conn),
		deadConn:   make(chan net.Conn),
		errors:     make(chan error),
	}
}

func (s *Server) Start(ctx context.Context, cancel context.CancelFunc) error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	go s.acceptConnections(ctx)

	select {
	case <-ctx.Done():
		return nil
	case err := <-s.errors:
		cancel()
		return err
	}
}

func (s *Server) acceptConnections(ctx context.Context) {
	go func() {
		for {
			conn, err := s.ln.Accept()
			logFatal(err)

			color.Info.Printf("%s --> has connected\n", conn.RemoteAddr().String())
			s.clients[conn] = true
			s.client <- conn
		}
	}()

	go s.handleConn()
}

func (s *Server) handleConn() {
	for {
		select {
		case conn := <-s.client:
			go s.broadcastMessages(conn)
		case conn := <-s.deadConn:
			msg := color.Red.Sprintf("%s -> disconnected\n", conn.RemoteAddr().String())

			for client := range s.clients {
				client.Write([]byte(msg))

				if client == conn {
					break
				}
			}

			fmt.Println(msg)

			delete(s.clients, conn)
		}
	}
}

func (s *Server) broadcastMessages(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		for client := range s.clients {
			if client != conn {
				msg = color.Yellow.Sprint(msg)
				client.Write([]byte(msg))
			}
		}
	}

	s.deadConn <- conn
}

// func (s *Server) broadcastMessages() {
// 	for {
// 		msg := <-s.messages

// 		for client := range s.clients {
// 			_, err := client.Write([]byte(msg + "\n"))
// 			if err != nil {
// 				delete(s.clients, client)
// 				fmt.Println("__error broadcasting msg__")
// 				return
// 			}
// 		}
// 	}
// }

// package main

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"net"
// )

// type Server struct {
// 	ln         net.Listener
// 	listenAddr string
// }

// var clients = make(map[net.Conn]bool)
// var messages = make(chan string)

// func newServer(listenAddr string) *Server {
// 	return &Server{
// 		listenAddr: listenAddr,
// 	}
// }

// func (s *Server) Start(ctx context.Context) error {
// 	ln, err := net.Listen("tcp", s.listenAddr)
// 	if err != nil {
// 		return err
// 	}
// 	defer ln.Close()
// 	s.ln = ln
// 	fmt.Println("Server started")

// 	go s.broadcastMessages()

// 	go s.acceptConn(ctx)

// 	<-ctx.Done()

// 	return nil
// }

// func (s *Server) acceptConn(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("__server closed__")
// 			return
// 		default:
// 			conn, err := s.ln.Accept()
// 			if err != nil {
// 				fmt.Println("error accepting => ", err.Error())
// 				return
// 			}
// 			clients[conn] = true

// 			conn.Write([]byte("Welcome to the server :)\n"))

// 			go s.handleConn(conn, ctx)
// 		}
// 	}
// }

// func (s *Server) handleConn(conn net.Conn, ctx context.Context) {

// 	username := conn.RemoteAddr().String()
// 	fmt.Printf("%s connected\n", username)

// 	defer func() {
// 		delete(clients, conn)
// 		fmt.Printf("%s --> disconnected\n", username)
// 	}()

// 	scanner := bufio.NewScanner(conn)
// 	for scanner.Scan() {
// 		message := scanner.Text()
// 		messages <- fmt.Sprintf("%s: %s", username, message)
// 	}
// }

// func (s *Server) broadcastMessages() {
// 	for {
// 		msg := <-messages

// 		for client := range clients {
// 			_, err := client.Write([]byte(msg + "\n"))
// 			if err != nil {
// 				fmt.Println("Error broadcasting message:", err.Error())
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }

// var msg = make(chan string, 1)

// func (s *Server) handleConn(conn net.Conn, ctx context.Context) {

// 	reader := bufio.NewReader(conn)
// 	msg, err := reader.ReadString('\n')
// 	logFatal(err)
// 	fmt.Println("msg: ", msg)
// 	// buf := make([]byte, 2048)

// 	// go func() {
// 	// 	for {
// 	// 		n, err := conn.Read(buf)
// 	// 		logFatal(err)

// 	// 		msg <- string(buf[:n])
// 	// 	}
// 	// }()

// 	// go func() {
// 	// 	for msgw := range msg {
// 	// 		fmt.Println("msgF: ", msgw)
// 	// 	}
// 	// }()

// 	usr := conn.RemoteAddr().String()
// 	fmt.Printf("%s --> connected\n", usr)

// 	defer func() {
// 		fmt.Printf("%s --> disconnected\n", usr)
// 		delete(s.clients, conn)
// 	}()

// 	scanner := bufio.NewScanner(conn)

// 	for scanner.Scan() {
// 		msg := scanner.Text()
// 		s.messages <- fmt.Sprintf("%s: %s", usr, msg)
// 	}
// }
