package hestur

import (
	"encoding/binary"
	"bytes"
)

type SendMapResponse struct {
}

func (SendMapResponse) ResponseID() ResponseID {
	return ResponseMap
}

func (r SendMapResponse) WriteData(server *Server, buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(server.Game.Map.Width))
	binary.Write(buf, binary.BigEndian, uint32(server.Game.Map.Height))

	cellsCount := server.Game.Map.Width * server.Game.Map.Height
	for i := 0; i < cellsCount; i++ {
		binary.Write(buf, binary.BigEndian, uint8(server.Game.Map.Cells[i].Type))
	}
}
