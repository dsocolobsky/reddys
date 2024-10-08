package database

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type AOF struct {
	file   *os.File
	reader *bufio.Reader
	mu     *sync.Mutex
}

func NewAOF(filepath string) *AOF {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	aof := &AOF{
		file:   f,
		reader: bufio.NewReader(f),
		mu:     &sync.Mutex{},
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
