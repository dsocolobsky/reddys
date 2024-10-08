package database

import (
	"fmt"
	"github.com/dsocolobsky/reddys/pkg/resp"
	"os"
	"sync"
	"time"
)

type Persister interface {
	Write(command string)
	Read() [][]string
}

type AOF struct {
	file *os.File
	mu   *sync.Mutex
}

func NewAOF(filepath string) *AOF {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	aof := &AOF{
		file: f,
		mu:   &sync.Mutex{},
	}
	// Sync every 1 second
	go func() {
		for {
			aof.mu.Lock()
			err := aof.file.Sync()
			if err != nil {
				panic(err)
			}
			aof.mu.Unlock()

			time.Sleep(1 * time.Second)
		}
	}()
	return aof
}

func (a *AOF) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()
	err := a.file.Close()
	if err != nil {
		panic(err)
	}
}

func (a *AOF) Write(command string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err := a.file.WriteString(command)
	if err != nil {
		panic(err)
	}
	fmt.Println("Wrote to AOF: ", command)
}

func (a *AOF) Read() [][]string {
	fmt.Println("Reading AOF from " + a.file.Name())
	a.mu.Lock()
	defer a.mu.Unlock()

	content, err := os.ReadFile(a.file.Name())
	if err != nil {
		fmt.Println("Error reading AOF file")
		return nil
	}
	if len(content) == 0 {
		fmt.Println("Empty AOF file")
		return nil
	}

	return resp.UnmarshalManyArrays(string(content))
}
