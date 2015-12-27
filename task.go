package common

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

type Action func(actionId string)

// Task task manager
type Task struct {
	name        string
	taskCounter int32
	action        Action
	parallels   int
	index       int32

	exitChan   chan bool
	notifyChan chan bool

	mutex     *sync.Mutex
	available bool
}

// NewTask create Task instance
func NewTask(name string, action Action) *Task {
	return &Task{
		name:        name,
		taskCounter: 0,
		action:        action,
		parallels:   0,
		index:       0,

		exitChan:   make(chan bool),
		notifyChan: make(chan bool),

		mutex:     &sync.Mutex{},
		available: false,
	}
}

// Start set factory available
func (p *Task) Start() {
	defer p.mutex.Unlock()
	p.mutex.Lock()

	p.available = true
}

// Close stop all action and set available false
func (p *Task) Close() {
	p.SetParallels(0)
	p.available = false
}

// SetParallels set parallel and increase or decrease task count by parallel value
func (p *Task) SetParallels(parallels int) bool {
	defer p.mutex.Unlock()
	p.mutex.Lock()

	if !p.available {
		return false
	}

	diff := parallels - p.parallels
	p.parallels = parallels
	for i := 0; i < int(math.Abs(float64(diff))); i++ {
		if diff < 0 {
			p.releaseTask()
		} else {
			p.addTask()
		}
	}

	return true
}

func (p *Task) addTask() {
	go p.doWork()

	atomic.AddInt32(&p.taskCounter, 1)
}

func (p *Task) releaseTask() {
	p.exitChan <- true
	<-p.notifyChan

	atomic.AddInt32(&p.taskCounter, -1)
}

func (p *Task) doWork() {
	//	taskID := uuid.Rand().Hex()
	index := atomic.AddInt32(&p.index, 1)
	for {
		select {
		case <-p.exitChan:
			p.notifyChan <- true
			return
		default:
			p.action(fmt.Sprintf("%d", index))
		}
	}
}
