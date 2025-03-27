package worker

type Result struct {
	JobID int
	State int // 0: Failed, 1: Success
}
