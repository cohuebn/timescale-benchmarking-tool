package workers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerAssignerIsStickyPerHost(test *testing.T) {
	workerAssigner := WorkerAssigner{
		numberOfWorkers: 10,
		assignedWorkers: make(map[string]int),
	}
	
	host := "host_12345"
	host1Result1 := workerAssigner.AssignHostToWorker(host)
	host1Result2 := workerAssigner.AssignHostToWorker(host)

	assert.Equal(test, host1Result1, host1Result2)
}

func TestWorkerAssignerDistributesLoadSomewhatEvenly(test *testing.T) {
	numberOfWorkers := 10
	workerAssigner := WorkerAssigner{
		numberOfWorkers: numberOfWorkers,
		assignedWorkers: make(map[string]int),
	}
	
	// Generate a bunch of hosts and validate their distribution across workers
	hostCount := 200
	workerAssignmentCounts := make([]int, numberOfWorkers)
	for i := 0; i < hostCount; i++ {
		host := fmt.Sprintf("host_%d", i)
		workerId := workerAssigner.AssignHostToWorker(host)
		workerAssignmentCounts[workerId]++
	}

	// Given 10 workers and 200 hosts, expect about 20 hosts per worker (200 / 10)
	expectedHostCountPerWorker := 20
	for _, workerAssignmentCount := range workerAssignmentCounts {
		// Allow +- 5 hosts per worker to account for randomness
		assert.InDelta(test, expectedHostCountPerWorker, workerAssignmentCount, 5)
	}
}