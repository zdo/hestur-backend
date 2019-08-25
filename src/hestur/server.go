package hestur

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type connection struct {
	Conn *websocket.Conn
	lock *sync.Mutex
	ID   uint64
}

type requestID uint16

const (
	requestIDPing requestID = 1
)

type responseID uint16

const (
	responseIDPing      responseID = 1
	responseIDMap                  = 2
	responseIDCharacter            = 3
)

// Server is a websocket-driven backend.
type Server struct {
	upgrader websocket.Upgrader
	game     *Game

	nextConnectionID uint64
	connectionsLock  *sync.RWMutex
	connections      map[uint64]*connection
}

// NewServer creates new server with specified game.
func NewServer(game *Game) Server {
	server := Server{game: game}
	server.nextConnectionID = 1
	server.connectionsLock = new(sync.RWMutex)
	server.connections = make(map[uint64]*connection)
	server.game = game
	return server
}

// Serve will start listening by websocket server.
func (server *Server) Serve() {
	var addr = "0.0.0.0:8080"
	server.upgrader = websocket.Upgrader{}

	server.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	http.HandleFunc("/", server.processRequest)

	go server.updateLoop()
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (server *Server) deleteConnection(c *connection) {
	server.connectionsLock.Lock()
	delete(server.connections, c.ID)
	server.connectionsLock.Unlock()
	c.Conn.Close()
}

func (server Server) writeBuffer(conn *websocket.Conn, buf *serverBuffer) {
	bytes := buf.Bytes()
	err := conn.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		conn.Close()
	}
}

func (server Server) writeResponse(conn *websocket.Conn, responseID responseID, fillFn func(*serverBuffer)) {
	buf := newServerBuffer()
	buf.Write(uint16(responseID))
	fillFn(&buf)
	server.writeBuffer(conn, &buf)
}

func (server Server) handleNewConnection(c *connection) {
	c.lock.Lock()
	server.writeResponseMap(c.Conn)

	server.game.mapObjectsLock.RLock()
	characters := server.game.Characters
	server.game.mapObjectsLock.RUnlock()

	for _, character := range characters {
		server.writeResponse(c.Conn, responseIDCharacter, func(buf *serverBuffer) {
			character.serialize(buf, SerializationTypeFull)
		})
	}

	c.lock.Unlock()
}

func (server Server) handleRequest(c *websocket.Conn, requestID requestID, reader *bytes.Reader) {
	switch requestID {
	case requestIDPing:
		server.writeResponsePing(c)
	}
}

func (server *Server) processRequest(w http.ResponseWriter, r *http.Request) {
	c, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := &connection{ID: server.nextConnectionID,
		lock: new(sync.Mutex),
		Conn: c}
	server.connectionsLock.Lock()
	server.nextConnectionID++
	server.connections[conn.ID] = conn
	server.connectionsLock.Unlock()

	defer server.deleteConnection(conn)

	server.handleNewConnection(conn)

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.BinaryMessage {
			break
		}

		conn.lock.Lock()

		reader := bytes.NewReader(message)

		var requestID requestID
		binary.Read(reader, binary.BigEndian, requestID)
		server.handleRequest(c, requestID, reader)

		conn.lock.Unlock()
	}
}

// UpdateLoop
func (server Server) updateLoop() {
	for true {
		time.Sleep(time.Second / 30.0)

		server.game.mapObjectsLock.RLock()
		characters := server.game.Characters

		server.connectionsLock.RLock()
		connections := server.connections

		for _, conn := range connections {
			conn.lock.Lock()
			for _, character := range characters {
				server.writeResponse(conn.Conn, responseIDCharacter, func(buf *serverBuffer) {
					character.serialize(buf, SerializationTypeShort)
				})
			}
			conn.lock.Unlock()
		}

		server.connectionsLock.RUnlock()
		server.game.mapObjectsLock.RUnlock()
	}
}
