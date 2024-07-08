package utils

import (
	"runtime"
	"sync"
)

func RunParallelTasks(tasks ...func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(task)
	}
	wg.Wait()
}

func RunParallelTasksEx(tasks ...func() error) error {
	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, len(tasks))
		result error
	)

	for _, task := range tasks {
		wg.Add(1)
		go func(f func() error) {
			defer wg.Done()
			if err := f(); err != nil {
				errCh <- err
			}
		}(task)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil && result == nil {
			result = err // 只记录第一个错误
		}
	}

	return result
}

func RunParallel[T any](slice []T, taskFunc func(slice []T)) {
	cpuNum := runtime.NumCPU()
	numPerTask := len(slice)/cpuNum + 1
	// 计算任务
	var tasks [][]T
	{
		taskNum := len(slice)/numPerTask + 1
		for i := 0; i < taskNum; i++ {
			start := i * numPerTask
			end := (i + 1) * numPerTask

			if end > len(slice) {
				end = len(slice)
			}
			tasks = append(tasks, slice[start:end])
		}
	}

	// 跑任务
	var pendingTasks []func()
	for _, slice := range tasks {
		task := func(slice []T) func() {
			return func() {
				taskFunc(slice)
			}
		}(slice)
		pendingTasks = append(pendingTasks, task)
	}
	RunParallelTasks(pendingTasks...)
}
