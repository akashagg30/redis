package aof

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akashagg30/redis/redis/utils"
)

// Constants
const (
	AOFFileName = "redis.aof" // Name of the append-only file
)

// AOFManager manages the append-only file.
type AOFManager struct {
	file      *os.File
	writer    *bufio.Writer
	mu        sync.Mutex // Protects file access
	lastSync  time.Time
	syncEvery time.Duration // How often to sync to disk
}

// NewAOFManager creates a new AOF manager.
func NewAOFManager(syncEvery time.Duration) (*AOFManager, error) {
	file, err := os.OpenFile(AOFFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("opening AOF file: %w", err)
	}

	return &AOFManager{
		file:      file,
		writer:    bufio.NewWriter(file),
		syncEvery: syncEvery,
	}, nil
}

// Close closes the AOF file.
func (aof *AOFManager) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	if err := aof.writer.Flush(); err != nil {
		return fmt.Errorf("flushing AOF writer: %w", err)
	}
	return aof.file.Close()
}

// LogCommand logs a command to the AOF file.
func (aof *AOFManager) LogCommand(wholeCommand ...string) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	respCmd := strconv.FormatInt(time.Now().UnixMicro(), 10) + " " + strings.Join(wholeCommand, " ") + "\n"

	_, err := aof.writer.WriteString(respCmd)
	if err != nil {
		return fmt.Errorf("writing to AOF: %w", err)
	}

	// Perform a sync if needed.
	if time.Since(aof.lastSync) >= aof.syncEvery {
		if err := aof.writer.Flush(); err != nil {
			return fmt.Errorf("flushing AOF writer: %w", err)
		}
		if err := aof.file.Sync(); err != nil {
			return fmt.Errorf("syncing AOF file: %w", err)
		}
		aof.lastSync = time.Now()
	}

	return nil
}

func (aof *AOFManager) Update(data ...any) {
	result := data[0]
	command := data[1].(string)
	switch command {
	case "SET":
		if result.(bool) { // if set command was successful and result is true
			wholeCommand, err := utils.ConvertToStringSlice(data[1:])
			if err == nil {
				aof.LogCommand(wholeCommand...)
			} else {
				log.Println(err)
			}
		}
	}
}

func (aof *AOFManager) ReplayAOF(commandChannel chan []string) { // work in progress
	file, err := os.Open(AOFFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokenizedLine := strings.Split(line, " ")
		wholeCommand := tokenizedLine[1:]
		commandChannel <- wholeCommand
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
