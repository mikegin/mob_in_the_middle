package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

const (
	HOST         = "0.0.0.0"
	PORT         = "8080"
	TYPE         = "tcp"
	TONY_ADDRESS = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

func main() {
	fmt.Println("Starting TCP Chat Proxy...")
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
		return
	}
	defer listen.Close()
	fmt.Println("Proxy listening on", HOST+":"+PORT)

	count := 0
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go handleRequest(conn, count)
		count++
	}
}

func handleRequest(conn net.Conn, id int) {
	defer conn.Close()

	fconn, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		log.Printf("Connection to upstream server failed: %v\n", err)
		return
	}
	defer fconn.Close()

	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)

	serverReader := bufio.NewReader(fconn)
	serverWriter := bufio.NewWriter(fconn)

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine to forward messages from client to server
	go func() {
		defer wg.Done()
		for {
			message, err := clientReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from client: %s", err)
				}
				break
			}

			fmt.Println(id, "Forwarding to server:", message)

			_, err = serverWriter.WriteString(message)
			if err != nil {
				log.Printf("Error writing to server: %s", err)
				break
			}
			serverWriter.Flush()
		}
	}()

	// Goroutine to read, modify, and forward messages from server to client
	go func() {
		defer wg.Done()
		re := regexp.MustCompile(`(^|\s)(7[a-zA-Z0-9]{25,34})(\s|$)`)

		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from server: %s", err)
				}
				break
			}

			fmt.Println(id, "Original from server: ", message)

			result := re.ReplaceAllStringFunc(message, func(match string) string {
				prefix := ""
				suffix := ""
				if strings.HasPrefix(match, " ") {
					prefix = " "
				}
				if strings.HasSuffix(match, " ") {
					suffix = " "
				}
				return prefix + TONY_ADDRESS + suffix
			})

			fmt.Println(id, "Modified to client: ", result)

			_, err = clientWriter.WriteString(result)
			if err != nil {
				log.Printf("Error writing to client: %s", err)
				break
			}
			clientWriter.Flush()
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}
