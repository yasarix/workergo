package workergo

import (
	"github.com/pborman/uuid"
)

// Job is the sturture that will be passed into workers
type Job struct {
	ID string

	Type JobType

	// Payload contains the pointer to the structures with methods that can
	// be called by workers
	Payload *interface{}

	// TargetFunc function to call to send payload when Type is MESSAGE.
	// method of the payload to call when Type is TASK
	TargetFunc string
}

// JobType has 2 possible values:
// 1- MESSAGE - Dummy message to be delivered to a given function
// 2- TASK - A struct with a method to be called
type JobType int

const (
	// MESSAGE Dummy message - interface{}
	MESSAGE JobType = iota
	// TASK A struct with a method to be called
	TASK
)

// NewJob Creates a new job instance with given payload
func NewJob(jobType JobType, payload interface{}, targetFunc string) Job {
	ID := uuid.NewUUID()

	return Job{
		ID:         ID.String(),
		Type:       jobType,
		TargetFunc: targetFunc,
		Payload:    &payload,
	}
}
