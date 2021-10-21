package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type TaskErrorCounter struct {
	sync.Mutex
	Errors int
}

func (c *TaskErrorCounter) Inc() int {
	c.Lock()
	c.Errors++
	errors := c.Errors
	c.Unlock()

	return errors
}

func (c *TaskErrorCounter) GetErrors() int {
	c.Lock()
	errors := c.Errors
	c.Unlock()

	return errors
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errCounter := TaskErrorCounter{}
	taskChan := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}

	// Запускаем worker'ов
	for i := 0; i < n; i++ {
		wg.Add(1)
		go consumer(taskChan, m, &errCounter, &wg)
	}

	// Перекинем все задачи в канал и закроем его
	for _, t := range tasks {
		taskChan <- t
	}
	close(taskChan)

	wg.Wait()

	if m >= 0 && errCounter.Errors >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func consumer(taskChan <-chan Task, maxErrors int, errCounter *TaskErrorCounter, wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range taskChan {
		err := t()
		if err != nil && maxErrors >= 0 {
			errors := errCounter.Inc()
			if errors >= maxErrors {
				return
			}
		}
	}
}
