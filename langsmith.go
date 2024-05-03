package mylangchaingo

import (
	"sync/atomic"
)

var runID, parentID, rootID atomic.Value

// SetRunId sets the run id
func SetRunId(runId string) {
	runID.Store(runId)
}

// GetRunId gets the run id
func GetRunId() string {
	if runID.Load() == nil {
		return ""

	}
	return runID.Load().(string)
}

// SetParentId sets the parent id
func SetParentId(parentId string) {
	parentID.Store(parentId)
}

// GetParentId gets the parent id
func GetParentId() string {
	if parentID.Load() == nil {
		return ""
	}
	return parentID.Load().(string)
}

// SetRootId sets the root id
func SetRootId(rootId string) {
	rootID.Store(rootId)
}

// GetRootId gets the root id
func GetRootId() string {
	if rootID.Load() == nil {
		return ""
	}
	return rootID.Load().(string)
}
