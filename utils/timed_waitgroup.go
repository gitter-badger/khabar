package utils

import (
	"log"
	"sync"
	"time"
)

type TimedWaitGroup struct {
	sync.WaitGroup
	TimeOut time.Duration
}

func (self *TimedWaitGroup) GetCompletionChannel(completionChannel chan<- bool) {
	self.Wait()
	completionChannel <- true

}

func (self *TimedWaitGroup) TimedWait() bool {
	completionChannel := make(chan bool)

	select {
	case <-completionChannel:
		log.Println("Success!!!!")
		return true
	case <-time.After(self.TimeOut):
		return false
	}

}
