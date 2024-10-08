package benchmarking

import (
	"hash/fnv"
)

// Assign a worker to run a CPU usage query for a given hostname
// This function uses a hash function to determinsitically assign the same
// worker to the same hostname
func assignHostToWorker(hostname string, numberOfWorkers int) int {
	// Use a relatively inexpensive hash function to convert a hostname
	// into a somewhat-randomized integer value
	hashedHostname := fnv.New32a()
	hashedHostname.Write([]byte(hostname))
	hashValue := hashedHostname.Sum32()

	// Given the hash value, determine which worker that hash value should go to
	return int(hashValue) % numberOfWorkers
}

type WorkerAssigner struct {
	numberOfWorkers int
	assignedWorkers map[string]int
}

func NewWorkerAssigner(numberOfWorkers int) WorkerAssigner {
	return WorkerAssigner{
		numberOfWorkers: numberOfWorkers,
		assignedWorkers: make(map[string]int),
	}
}

// Assign a the given host to a worker. This function will
// reuse the same worker for the same host if called multiple times
// for that host.
func (workerAssigner *WorkerAssigner) AssignHostToWorker(host string) int {
	// Use the already-assigned worker if there is one
	alreadyAssignedWorker, ok := workerAssigner.assignedWorkers[host]
	if (ok) {
		return alreadyAssignedWorker
	}

	// Assign a worker to the host if we haven't already
	assignedWorker := assignHostToWorker(host, workerAssigner.numberOfWorkers)
	workerAssigner.assignedWorkers[host] = assignedWorker
	return assignedWorker
}