package resp

import (
	"bytes"
	"io"
	"sync"
)

type waitByteReader struct {
	*bytes.Reader
	dataCh chan []byte
	mu     sync.Mutex
}

func newWaitByteReader(data []byte) *waitByteReader {
	return &waitByteReader{
		Reader: bytes.NewReader(data),
		dataCh: make(chan []byte),
	}
}

func (w *waitByteReader) ReadByte() (byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// If we're at the end, block until new data arrives
	for {
		b, err := w.Reader.ReadByte()
		if err == nil {
			return b, nil
		}

		// If EOF is reached, we will wait for new data to be written to dataCh
		if err == io.EOF {
			newBytes, ok := <-w.dataCh
			if !ok {
				return 0, io.EOF
			}
			// When new data is received, continue reading
			w.Reader.Reset(newBytes)
		}
	}
}

func (w *waitByteReader) writeNewData(data []byte) {
	w.dataCh <- data
}

func (w *waitByteReader) close() {
	close(w.dataCh)
}
