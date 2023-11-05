package job

import "fmt"

type JobCommand struct{}

func (j *JobCommand) Name() string {
	return "start_server_job"
}

func (j *JobCommand) Description() string {
	return "Start Server Job HTTP/1.1"
}

func (j *JobCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("JobCommand Args is required")
	}
	// Do here
	return nil
}
