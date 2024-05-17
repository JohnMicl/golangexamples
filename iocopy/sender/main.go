package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: sender <file> <receiver_address:port>")
	}
	fileName := os.Args[1]
	addr := os.Args[2]

	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 使用io.Copy
	startTime := time.Now()
	n, err := io.Copy(conn, file)
	if err != nil {
		log.Fatalf("io.Copy failed: %v", err)
	}
	durationCopy := time.Since(startTime)
	fmt.Printf("io.Copy sent %d bytes in %v\n", n, durationCopy)

	// 重置文件指针以便再次读取
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Seek failed: %v", err)
	}

	// 手动读写
	startTime = time.Now()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal("Read failed:", err)
		}
		if n == 0 {
			break
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Fatal("Write failed:", err)
		}
	}
	durationManual := time.Since(startTime)
	fmt.Printf("Manual copy sent %d bytes in %v\n", n, durationManual)

	// 比较两种方法的效率
	fmt.Println("io.Copy was", durationManual/durationCopy, "times faster than manual copy")
}
