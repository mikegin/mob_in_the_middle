package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
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

		re := regexp.MustCompile(`(^|\s)(7[a-zA-Z0-9]{25,34})(\s|$)`)

		result := re.ReplaceAllStringFunc(fmessage, func(match string) string {
			// Check if the match contains leading or trailing whitespace
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

		_, err = writer.WriteString(result)
		if err != nil {
			log.Printf("Error writing to client: %s", err)
			break
		}
		writer.Flush()
	}
}
