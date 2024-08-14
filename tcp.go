package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	openTCP []int

	buffer []byte
	banner string
)

func scanTCP(ip string, port int, wg *sync.WaitGroup, openTCP *[]int, mu *sync.Mutex) {
	defer wg.Done()
	protocol := "tcp"
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return
	}

	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Second))
	buffer = make([]byte, 256)
	n, err := conn.Read(buffer)
	banner = string(buffer[:n])

	mu.Lock()
	*openTCP = append(*openTCP, port)
	mu.Unlock()

}

func main() {

	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("Choose host to scan: ")
	fmt.Fscan(
		os.Stdin,
		&ip,
	)

	startPort := 1
	endPort := 1024

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanTCP(ip, port, &wg, &openTCP, &mu)

	}

	wg.Wait()

	sort.Ints(openTCP)
	fmt.Printf("Open TCP ports: %v\n", openTCP)
	fmt.Print("Banner: ", banner)

}
