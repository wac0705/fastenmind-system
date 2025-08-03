package concurrent

import (
	"context"
	"errors"
	"sync"
	"time"
)

// WorkerPool 工作池
type WorkerPool struct {
	workers    int
	taskQueue  chan Task
	workerWg   sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	maxRetries int
	retryDelay time.Duration
}

// Task 任務介面
type Task interface {
	Execute(ctx context.Context) error
	OnError(err error)
	OnSuccess()
}

// NewWorkerPool 創建新的工作池
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workers:    workers,
		taskQueue:  make(chan Task, queueSize),
		ctx:        ctx,
		cancel:     cancel,
		maxRetries: 3,
		retryDelay: time.Second,
	}
}

// Start 啟動工作池
func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.workerWg.Add(1)
		go p.worker(i)
	}
}

// Submit 提交任務
func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.taskQueue <- task:
		return nil
	case <-p.ctx.Done():
		return errors.New("worker pool is shutting down")
	default:
		return errors.New("task queue is full")
	}
}

// SubmitWithTimeout 提交任務（帶超時）
func (p *WorkerPool) SubmitWithTimeout(task Task, timeout time.Duration) error {
	select {
	case p.taskQueue <- task:
		return nil
	case <-time.After(timeout):
		return errors.New("submit timeout")
	case <-p.ctx.Done():
		return errors.New("worker pool is shutting down")
	}
}

// Shutdown 關閉工作池
func (p *WorkerPool) Shutdown() {
	p.cancel()
	close(p.taskQueue)
	p.workerWg.Wait()
}

// ShutdownWithTimeout 關閉工作池（帶超時）
func (p *WorkerPool) ShutdownWithTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	
	go func() {
		p.Shutdown()
		close(done)
	}()
	
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.New("shutdown timeout")
	}
}

// worker 工作協程
func (p *WorkerPool) worker(id int) {
	defer p.workerWg.Done()
	
	for {
		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				return // 隊列已關閉
			}
			
			p.executeTaskWithRetry(task)
			
		case <-p.ctx.Done():
			return
		}
	}
}

// executeTaskWithRetry 執行任務（帶重試）
func (p *WorkerPool) executeTaskWithRetry(task Task) {
	var err error
	
	for i := 0; i <= p.maxRetries; i++ {
		ctx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
		err = task.Execute(ctx)
		cancel()
		
		if err == nil {
			task.OnSuccess()
			return
		}
		
		if i < p.maxRetries {
			time.Sleep(p.retryDelay * time.Duration(i+1))
		}
	}
	
	task.OnError(err)
}

// ConnectionPool 連接池
type ConnectionPool[T any] struct {
	pool      chan T
	factory   func() (T, error)
	validator func(T) bool
	closer    func(T) error
	maxSize   int
	mu        sync.Mutex
}

// NewConnectionPool 創建新的連接池
func NewConnectionPool[T any](maxSize int, factory func() (T, error), validator func(T) bool, closer func(T) error) *ConnectionPool[T] {
	return &ConnectionPool[T]{
		pool:      make(chan T, maxSize),
		factory:   factory,
		validator: validator,
		closer:    closer,
		maxSize:   maxSize,
	}
}

// Get 獲取連接
func (p *ConnectionPool[T]) Get() (T, error) {
	select {
	case conn := <-p.pool:
		if p.validator != nil && !p.validator(conn) {
			// 連接無效，創建新連接
			if p.closer != nil {
				p.closer(conn)
			}
			return p.factory()
		}
		return conn, nil
	default:
		// 池中沒有可用連接，創建新連接
		return p.factory()
	}
}

// Put 歸還連接
func (p *ConnectionPool[T]) Put(conn T) {
	if p.validator != nil && !p.validator(conn) {
		// 連接無效，關閉它
		if p.closer != nil {
			p.closer(conn)
		}
		return
	}
	
	select {
	case p.pool <- conn:
		// 成功放回池中
	default:
		// 池已滿，關閉連接
		if p.closer != nil {
			p.closer(conn)
		}
	}
}

// Close 關閉連接池
func (p *ConnectionPool[T]) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	close(p.pool)
	
	// 關閉所有連接
	for conn := range p.pool {
		if p.closer != nil {
			if err := p.closer(conn); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// Size 獲取池中連接數
func (p *ConnectionPool[T]) Size() int {
	return len(p.pool)
}