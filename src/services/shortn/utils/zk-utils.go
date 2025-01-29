package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	zookeeperHosts    = "zookeeper:2181" // Zookeeper connection string
	rootPath          = "/shortn-ranges" // Root znode for ranges
	initialCounter    = 1                // Starting counter value
	rangeSize         = 10               // Range size for each node
	connectionTimeout = 5 * time.Second  // Zookeeper connection timeout
)

func CreateZkConnection() *zk.Conn {
	// Connect to Zookeeper
	conn, _, err := zk.Connect([]string{zookeeperHosts}, connectionTimeout)
	if err != nil {
		log.Fatalf("Failed to connect to Zookeeper: %v", err)
	}
	//defer conn.Close()
	ensurePathExists(conn, rootPath)
	return conn
}

// ensurePathExists ensures that the given path exists in Zookeeper
func ensurePathExists(conn *zk.Conn, path string) {
	exists, _, err := conn.Exists(path)
	if err != nil {
		log.Fatalf("Failed to check existence of path %s: %v", path, err)
	}

	if !exists {
		_, err := conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalf("Failed to create path %s: %v", path, err)
		}
	}
}

// allocateRange allocates a unique range of numbers for this node
func AllocateRange(conn *zk.Conn) (int, int, error) {
	counterPath := fmt.Sprintf("%s/counter", rootPath)

	// Check if the counter node exists
	exists, _, err := conn.Exists(counterPath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to check counter: %v", err)
	}

	if !exists {
		// Initialize the counter if it doesn't exist
		_, err := conn.Create(counterPath, intToBytes(initialCounter), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to initialize counter: %v", err)
		}
	}

	// Atomically allocate a range
	for {
		// Get the current counter value
		data, stat, err := conn.Get(counterPath)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get counter: %v", err)
		}

		currentCounter := bytesToInt(data)

		// Calculate the new range
		start := currentCounter
		end := currentCounter + rangeSize - 1

		// Update the counter atomically
		_, err = conn.Set(counterPath, intToBytes(end+1), stat.Version)
		if err == zk.ErrBadVersion {
			// Retry if another node modified the counter
			continue
		} else if err != nil {
			return 0, 0, fmt.Errorf("failed to update counter: %v", err)
		}

		// Return the allocated range
		return start, end, nil
	}
}
