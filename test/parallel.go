package test

import (
	"fmt"
	"time"
)

func main() {
	list := make([]string, 0, 4)
	list = append(list, "a")
	list = append(list, "b")
	list = append(list, "c")
	list = append(list, "d")
	
	taskChan := make(chan string, 10)
	for _, item := range list {
		taskChan <- item
	}

	outter:for {
		select {
		case task, ok := <- taskChan:
			if !ok {
				fmt.Println("Channel closed, exiting loop.")
				break outter
			}
			fmt.Println("Processing task:", task)
			if task == "b" {
				taskChan <- "b1"
			}
			if task == "b1" {
				taskChan <- "b2"
			}
		default:
			fmt.Println("No tasks to process, exiting loop.")
			if len(taskChan) == 0 {
				close(taskChan)
				fmt.Println("Channel closed.")
				break outter
			}
			time.Sleep(time.Second)
		}
	}
	fmt.Println("All tasks completed.")
}