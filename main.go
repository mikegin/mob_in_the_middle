package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"unicode"
)

func main() {
	const (
		HOST = "0.0.0.0"
		PORT = "8080"
		TYPE = "tcp"
	)

	fmt.Println("Starting TCP Chat Proxy...")
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
		return
	}
	defer listen.Close()
	fmt.Println("Proxy listening on", HOST+":"+PORT)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	fconn, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}

	defer fconn.Close()

	freader := bufio.NewReader(fconn)
	fwriter := bufio.NewWriter(fconn)

	// Initial welcome message from server
	message, err := freader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading from connection: %s", err)
		}
		return
	}

	_, err = writer.WriteString(message)
	if err != nil {
		log.Printf("Error writing to client: %s", err)
		return
	}
	writer.Flush()

	// goroutine to forward messages
	go func() {
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from client: %s", err)
				}
				break
			}

			_, err = fwriter.WriteString(message)
			if err != nil {
				log.Printf("Error writing to server: %s", err)
				break
			}
			fwriter.Flush()
		}
	}()

	// read, modify and respond
	for {
		fmessage, err := freader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from server: %s", err)
			}
			break
		}

		a := strings.Split(fmessage[:len(fmessage)-1], " ")

		for i, s := range a {
			l := len(s)
			isAddress := 26 <= l && l <= 35 && s[0] == '7'

			if !isAddress {
				continue
			}

			for _, c := range s {
				if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
					isAddress = false
					break
				}
			}

			if isAddress {
				a[i] = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
			}
		}

		result := strings.Join(a, " ")
		_, err = writer.WriteString(result + "\n")
		if err != nil {
			log.Printf("Error writing to client: %s", err)
			break
		}
		writer.Flush()
	}
}
