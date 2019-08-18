package hestur

import "github.com/gorilla/websocket"

func (server Server) writeResponsePing(conn *websocket.Conn) {
	server.writeResponse(conn, responseIDPing, func(buf *serverBuffer) {})
}

func (server Server) writeResponseMap(conn *websocket.Conn) {
	server.writeResponse(conn, responseIDMap, func(buf *serverBuffer) {
		buf.Write(uint32(server.game.Width))
		buf.Write(uint32(server.game.Height))

		cellsCount := server.game.Width * server.game.Height
		for i := 0; i < cellsCount; i++ {
			buf.Write(uint8(server.game.Cells[i].Type))
		}
	})
}
