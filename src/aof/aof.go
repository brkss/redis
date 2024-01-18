package aof

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/brkss/redis/src/resp"
)

// AOF is struct that hold file / reader and mutex related to data sync
type AOF struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

// NewAOF create a new AOF object create or open file provided in params
// alse lunch a parallel thread to keep the file synced
func NewAOF(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &AOF{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// seperated thread to sync the AOF to disc every 1s
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

// Close closes the file after the problem exited
func (aof *AOF) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *AOF) Write(val resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(val.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *AOF) Read(fn func(val resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, io.SeekStart)
	reader := resp.NewReader(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fn(value)
	}

	return nil
}
