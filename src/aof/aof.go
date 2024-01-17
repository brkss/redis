package aof

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/brkss/redis/src/resp"
)

type AOF struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

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

func (aof *AOF) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *AOF) Write(val resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	fmt.Println("file : ", aof.file.Fd())
	_, err := aof.file.Write(val.Marshal())
	if err != nil {
		return err
	}

	return nil
}
