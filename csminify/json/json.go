package json

import (
	"encoding/json"
	rep "github.com/markus-wa/cs-demo-minifier/csminify/replay"
	"io"
)

func MarshalReplay(replay rep.Replay, w io.Writer) error {
	return json.NewEncoder(w).Encode(replay)
}
