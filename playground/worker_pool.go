package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID        int
	Payload   string
	Duration  time.Duration
	Priority  int
}

type Result struct {
	TaskID   int
	Output   string
	Error    error
	Duration time.Duration
}

type WorkerPool struct {
	workers        int
	taskQueue      chan Task
	resultQueue    chan Result
	errorQueue     chan error
	quit           chan struct{}
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
	activeWorkers  int
	completedTasks int
}

func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:      workers,
		taskQueue:    make(chan Task, 100),
		resultQueue:  make(chan Result, 100),
		errorQueue:   make(chan error, 10),
		quit:         make(chan struct{}),
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	go wp.monitor()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			fmt.Printf("Worker %d shutting down...\n", id)
			return
		case task := <-wp.taskQueue:
			wp.mu.Lock()
			wp.activeWorkers++
			wp.mu.Unlock()

			fmt.Printf("Worker %d started task %d (priority: %d)\n", id, task.ID, task.Priority)
			
			result := wp.processTask(task)
			
			wp.mu.Lock()
			wp.activeWorkers--
			wp.completedTasks++
			wp.mu.Unlock()

			select {
			case wp.resultQueue <- result:
			case <-wp.ctx.Done():
				return
			}
		}
	}
}

func (wp *WorkerPool) processTask(task Task) Result {
	startTime := time.Now()
	
	select {
	case <-time.After(task.Duration):
	case <-wp.ctx.Done():
		return Result{
			TaskID:   task.ID,
			Error:    fmt.Errorf("task cancelled"),
			Duration: time.Since(startTime),
		}
	}

	output := fmt.Sprintf("Processed task %d: %s", task.ID, task.Payload)
	
	if rand.Intn(10) == 0 {
		return Result{
			TaskID:   task.ID,
			Error:    fmt.Errorf("random error processing task %d", task.ID),
			Duration: time.Since(startTime),
		}
	}

	return Result{
		TaskID:   task.ID,
		Output:   output,
		Duration: time.Since(startTime),
	}
}

func (wp *WorkerPool) monitor() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wp.mu.Lock()
			fmt.Printf("Monitor: Active workers: %d, Completed tasks: %d, Queue length: %d\n",
				wp.activeWorkers, wp.completedTasks, len(wp.taskQueue))
			wp.mu.Unlock()
		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) SubmitTask(task Task) error {
	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	}
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.resultQueue
}

func (wp *WorkerPool) Errors() <-chan error {
	return wp.errorQueue
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.quit)
	wp.wg.Wait()
	close(wp.taskQueue)
	close(wp.resultQueue)
	close(wp.errorQueue)
}

func (wp *WorkerPool) GetStats() (int, int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.activeWorkers, wp.completedTasks
}

func taskGenerator(ctx context.Context, taskChan chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(taskChan)

	taskID := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task generator shutting down...")
			return
		default:
			task := Task{
				ID:       taskID,
				Payload:  fmt.Sprintf("Data payload %d", taskID),
				Duration: time.Duration(rand.Intn(1000)+100) * time.Millisecond,
				Priority: rand.Intn(5),
			}
			
			select {
			case taskChan <- task:
				taskID++
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			case <-ctx.Done():
				return
			}
		}
	}
}

func resultCollector(ctx context.Context, resultChan <-chan Result, errorChan <-chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				return
			}
			if result.Error != nil {
				fmt.Printf("❌ Task %d failed: %v (took %v)\n", result.TaskID, result.Error, result.Duration)
			} else {
				fmt.Printf("✅ Task %d completed: %s (took %v)\n", result.TaskID, result.Output, result.Duration)
			}
		case err, ok := <-errorChan:
			if !ok {
				return
			}
			fmt.Printf("🚨 System error: %v\n", err)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := NewWorkerPool(5)
	pool.Start()
	defer pool.Stop()

	var wg sync.WaitGroup

	taskChan := make(chan Task, 50)
	wg.Add(1)
	go taskGenerator(ctx, taskChan, &wg)

	wg.Add(1)
	go resultCollector(ctx, pool.Results(), pool.Errors(), &wg)

	go func() {
		for task := range taskChan {
			if err := pool.SubmitTask(task); err != nil {
				fmt.Printf("Failed to submit task %d: %v\n", task.ID, err)
				break
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All tasks completed!")
	case <-ctx.Done():
		fmt.Println("Context timeout reached!")
	}

	active, completed := pool.GetStats()
	fmt.Printf("Final stats - Active workers: %d, Completed tasks: %d\n", active, completed)
}