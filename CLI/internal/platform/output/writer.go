package output

import (
	"encoding/json"
	"fmt"
)

func WriteRaw(streams Streams, body []byte) error {
	_, err := streams.Stdout.Write(body)
	return err
}

func WriteJSON(streams Streams, v any) error {
	enc := json.NewEncoder(streams.Stdout)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func WriteError(streams Streams, err error) error {
	_, writeErr := fmt.Fprintln(streams.Stderr, err.Error())
	return writeErr
}
