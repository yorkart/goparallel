package common

import (
	"fmt"
	"testing"
	"time"
)

func TestTaskFactory(t *testing.T) {
	taskFactory := NewTask("tets", func(taskId string) {
		fmt.Println(taskId + " do work")
		time.Sleep(5 * time.Second)
	})

	taskFactory.Start()

	taskFactory.SetParallels(2)

	time.Sleep(2 * time.Second)
	taskFactory.SetParallels(4)

	time.Sleep(2 * time.Second)
	taskFactory.SetParallels(0)
	fmt.Println("set 0")

	time.Sleep(5 * time.Second)
	taskFactory.SetParallels(10)

	time.Sleep(2 * time.Second)
	taskFactory.Close()

	time.Sleep(2 * time.Second)
}
