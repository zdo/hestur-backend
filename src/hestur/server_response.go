package hestur

import "bytes"

type ServerResponse interface {
	ResponseID() ResponseID
	WriteData(server *Server, buf *bytes.Buffer)
}
