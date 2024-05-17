package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: receiver <port>")
	}
	port := os.Args[1]

	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	f, err := ioutil.TempFile("", "received_file_*.dat")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("ReadFromUDP error: %v", err)
			continue
		}
		if n == 0 {
			break
		}
		if _, err := f.Write(buffer[:n]); err != nil {
			log.Printf("Write error: %v", err)
		}
	}
	fmt.Printf("File received and saved as %s\n", f.Name())
}
