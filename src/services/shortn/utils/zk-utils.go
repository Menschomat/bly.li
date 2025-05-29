package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	zookeeperHosts    = "zookeeper:2181" // Zookeeper connection string
	rootPath          = "/shortn-ranges" // Root znode for ranges
	initialCounter    = 1                // Starting counter value
	rangeMin          = 200
	rangeMax          = 500             // Range size for each node
	connectionTimeout = 5 * time.Second // Zookeeper connection timeout
)

type zkLogger struct{}

func (zkLogger) Printf(format string, v ...interface{}) {
	logger.Info("zookeeper", "msg", strings.TrimSpace(fmt.Sprintf(format, v...)))
}

func createZkConnection() *zk.Conn {
	// Connect to Zookeeper
	conn, _, err := zk.Connect([]string{zookeeperHosts}, connectionTimeout, zk.WithLogger(zkLogger{}))
	if err != nil {
		logger.Error("Failed to connect to Zookeeper", "err", err)
		os.Exit(1)
	}
	//defer conn.Close()
	ensurePathExists(conn, rootPath)
	return conn
}

// ensurePathExists ensures that the given path exists in Zookeeper
func ensurePathExists(conn *zk.Conn, path string) {
	exists, _, err := conn.Exists(path)
	if err != nil {
		logger.Error("Failed to check existence of path",
			"path", path,
			"err", err,
		)
		os.Exit(1)
	}

	if !exists {
		_, err := conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			logger.Error("Failed to create path",
				"path", path,
				"err", err,
			)
			os.Exit(1)
		}
	}
}

// allocateRange allocates a unique range of numbers for this node
func AllocateRange() (int, int, error) {
	conn := createZkConnection()
	defer conn.Close()

	counterPath := fmt.Sprintf("%s/counter", rootPath)

	// Check if the counter node exists
	exists, _, err := conn.Exists(counterPath)
	if err != nil {
		logger.Error("Failed to check counter existence", "path", counterPath, "error", err)
		return 0, 0, fmt.Errorf("failed to check counter: %w", err)
	}

	if !exists {
		// Initialize the counter if it doesn't exist
		_, err := conn.Create(counterPath, intToBytes(initialCounter), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			logger.Error("Failed to initialize counter", "path", counterPath, "error", err)
			return 0, 0, fmt.Errorf("failed to initialize counter: %w", err)
		}
	}

	// Atomically allocate a range
	for {
		data, stat, err := conn.Get(counterPath)
		if err != nil {
			logger.Error("Failed to get counter value", "path", counterPath, "error", err)
			return 0, 0, fmt.Errorf("failed to get counter: %w", err)
		}

		currentCounter := bytesToInt(data)

		start := currentCounter
		end := currentCounter + GetRandomIntInRange(200, 500) - 1

		_, err = conn.Set(counterPath, intToBytes(end+1), stat.Version)
		if err == zk.ErrBadVersion {
			// Retry if another node modified the counter
			continue
		} else if err != nil {
			logger.Error("Failed to update counter", "path", counterPath, "error", err)
			return 0, 0, fmt.Errorf("failed to update counter: %w", err)
		}

		logger.Info("Allocated counter range", "start", start, "end", end)
		return start, end, nil
	}
}
