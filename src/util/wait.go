package util

import "sync"

var (
	waitGroup *sync.WaitGroup = &sync.WaitGroup{}
)

func GetWaitGroup() *sync.WaitGroup {
	return waitGroup
}
