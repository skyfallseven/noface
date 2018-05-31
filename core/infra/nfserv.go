/*
NoFace Pre-Dev Testing
Assembling a server, or sort of
*/
package main

import (
	"net"
	"os"
	"fmt"
	"time"
)
func handleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()
	//var buf [512]byte
	//for {
		/*
		//read up to 512 bytes
		_, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		*/
		//write the n bytes read
	fmt.Println("New Connection")
	_, err2 := conn.Write([]byte("Holding for 60 seconds."))
	if err2 != nil {
		return
	}
	//}
	time.Sleep(60 * time.Second)

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nFatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	port := "65000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:"+ port)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// run as a goroutine
		go handleClient(conn)
	}
}
