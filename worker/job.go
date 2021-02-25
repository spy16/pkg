package worker

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Job represents the job request to the worker.
type Job struct {
	ID      string          // unique identifier for the job.
	Ack     func(err error) // acknowledge callback when message is processed.
	Time    time.Time       // time at which this job was created.
	Error   error           // error during processing if any.
	Payload interface{}     // payload of the message.
}

// EnsureValid sets defaults for unset fields where possible and validates the
// job definition.
func (j *Job) EnsureValid() error {
	if j.Ack == nil {
		j.Ack = func(_ error) { /* do nothing */ }
	}
	if j.ID == "" {
		return errors.New("job id must be present")
	}
	return nil
}

func (j Job) String() string {
	var parts []string
	if j.Error != nil {
		parts = append(parts, fmt.Sprintf("error='%s'", j.Error))
	}

	if _, ok := j.Payload.(fmt.Stringer); ok {
		parts = append(parts, fmt.Sprintf("payload='%s'", j.Payload))
	}

	return fmt.Sprintf("Job{%s}", strings.Join(parts, ", "))
}
