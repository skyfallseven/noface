/*
NoFace Development Build
Infrastructure Core 
Adapting TCP/IP Networking from AppliedGo
Dis gon be tuff
*/
package main

import (
	"bufio"
	"io"
	"log"
	"net"
	//"strconv"
	//"strings"
	"sync"

	"github.com/pkg/errors"
	"encoding/gob"
	//"flag"
)

type Action struct {
	Type string //LOGIN LOGOUT SEND DELETE ACCOUNT
	UID int //User ID of client (server is ID 0)
	Data string //Each action has it's own parser
	Token string //Used to authenticate 
}

// handles incoming Actions, receives open connection in a ReadWriter
type HandleFunc func(Action) Action

// Provides an endpoint to other processes that data can be sent to
type Endpoint struct {
	listener net.Listener
	handler map[string]HandleFunc
	m sync.RWMutex
}

/*
Name: 	newEndpoint
Return:	Endpoint
*/
func NewEndpoint() *Endpoint {
	return &Endpoint {
		handler : map[string]HandleFunc{},
	}
}

/*
Name: 	AddHandleFunc
Param:	name - of handlefunction to add
		f - Actual handleFunc to add to endpoint
*/
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f // add handle func to endpoint
	e.m.Unlock()
}

/*
Name:	Listen
Param:	port - port to listen on
Return:	error on fail
*/
func (e *Endpoint) Listen(port string) error {
	var err error
	// Start listening
	e.listener, err = net.Listen("tcp", port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", port)
	}
	log.Println("Listen on", e.listener.Addr().String())

	// Start accepting connections
	// Client only needs one connection request from server
	// Should not use more
	for {
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting connection:", err)
			continue
		}
		log.Println("Connection from " + conn.RemoteAddr().String())
		log.Println("Start handling Actions")
		go e.handleMessages(conn)
	}

}

/*
Name:	decodeAction
Param:	rw - ReadWriter that sends us the gob
Return:	The Action that was decoded
*/
func DecodeAction(rw *bufio.ReadWriter) Action {
	var act Action
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&act)
	if err != nil {
		log.Println("Could not decode GOB")
	}
	return act
}

/*
Name:	handleMessages
Param:	conn - connection from user/server
Return:	who the fk knows
*/
func (e *Endpoint) handleMessages(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// The client or server will send a string 
	// Either "CLIENT" or "SERVER" 
	for {
		log.Print("Receive identifier: ")
		id, err := rw.ReadString('\n')

		switch {
		case err == io.EOF:
			log.Println("Reached EOF. Terminate connection.\n ---")
			return
		case err != nil:
			log.Println("Error reading string:", err)
			return
		case id != "CLIENT" && id != "SERVER":
			return //to prevent some junk from entering the port
		}
		//Read the GOB that was sent
		act := DecodeAction(rw)
		log.Println("Recieved type:", act.Type)

		e.m.RLock()
		handleCommand, ok := e.handler[act.Type] //select handler
		e.m.RUnlock()
		if !ok {
			log.Println("Unregistered Action:", act.Type)
		}
		handleCommand(act)
	}

}
const (
	CPort = ":65000" //Server port ~Read from config
	SPort = ":65001" //Client port ~Read from config
)

/*
Name: 	connToServer
Param:	addr - address of the NF server
Return:	bufio reader and writer
*/
func ConnToServer(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Connection to "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func tryLogin() Action{
	act := Action{
		Type: "LOGIN",
		UID: 123456789,
		Data: "PGP-Base64Stuff",
		Token: "SuperCoolAuthToken"}
	return act
}

func sendAction(act Action) {
	rw, err := ConnToServer("localhost:65001")
	if err != nil {
		log.Println("Connection failed")
	}
	n, err := rw.WriteString("CLIENT\n")
	log.Println(n)
	// Try to send gob
	enc := gob.NewEncoder(rw)
	err = enc.Encode(act)
	if err != nil {
		log.Println("Encoding failed")
	}
	rw.Flush()
	rw.WriteString("DONE\n")
}

func main() {
	sendAction(tryLogin())
}
