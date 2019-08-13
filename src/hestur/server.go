package hestur

import (
	"log"
	"net/http"
	"encoding/binary"
	"bytes"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn *websocket.Conn
	ID uint64
}

type RequestID uint16
const (
	RequestPing RequestID = 1
)

type ResponseID uint16
const (
	ResponsePing ResponseID = 1
	ResponseMap = 2
)

type Server struct {
	upgrader websocket.Upgrader
	Game Game

	NextConnectionID uint64
	Connections map[uint64]Connection
}

func (server Server) Start() {
	server.NextConnectionID = 1
	server.Connections = map[uint64]Connection{}

	server.Game = Game{Map: CreateEmptyMap(2048, 2048)}

	var addr = "0.0.0.0:8080"
	server.upgrader = websocket.Upgrader{}
	server.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	http.HandleFunc("/", server.ProcessRequest)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (server Server) DeleteConnection(c Connection) {
	delete(server.Connections, c.ID)
	c.Conn.Close()
}

func (server Server) WriteTo(conn *websocket.Conn, res ServerResponse) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint16(res.ResponseID()))
	res.WriteData(&server, buf)

	log.Printf("Write %d", int(res.ResponseID()))

	bytes := buf.Bytes()
	err := conn.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		conn.Close()
	}
}

func (server Server) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	c, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := Connection{ID: server.NextConnectionID, Conn: c}
	server.NextConnectionID++
	server.Connections[conn.ID] = conn

	defer server.DeleteConnection(conn)

	sendMap := SendMapResponse{}
	server.WriteTo(c, sendMap)

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.BinaryMessage {
			break
		}

		reader := bytes.NewReader(message)

		var requestID RequestID
		binary.Read(reader, binary.BigEndian, requestID)

		writer := new(bytes.Buffer)

		switch requestID {
		case RequestPing:
			var ok uint32 = 1
			binary.Write(writer, binary.BigEndian, ok)
		}

		err = c.WriteMessage(websocket.BinaryMessage, writer.Bytes())
		if err != nil {
			break
		}
	}
}
