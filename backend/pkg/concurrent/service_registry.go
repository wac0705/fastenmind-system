package concurrent

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Service 服務介面
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
	Status() ServiceStatus
}

// ServiceStatus 服務狀態
type ServiceStatus string

const (
	StatusStopped  ServiceStatus = "stopped"
	StatusStarting ServiceStatus = "starting"
	StatusRunning  ServiceStatus = "running"
	StatusStopping ServiceStatus = "stopping"
	StatusError    ServiceStatus = "error"
)

// ServiceInfo 服務資訊
type ServiceInfo struct {
	Name      string        `json:"name"`
	Status    ServiceStatus `json:"status"`
	StartTime time.Time     `json:"start_time"`
	Error     error         `json:"error,omitempty"`
}

// ServiceRegistry 服務註冊表
type ServiceRegistry struct {
	services   map[string]Service
	infos      map[string]*ServiceInfo
	mu         sync.RWMutex
	startOrder []string // 服務啟動順序
	stopOrder  []string // 服務停止順序（通常是啟動順序的反向）
}

// NewServiceRegistry 創建新的服務註冊表
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services:   make(map[string]Service),
		infos:      make(map[string]*ServiceInfo),
		startOrder: make([]string, 0),
		stopOrder:  make([]string, 0),
	}
}

// Register 註冊服務
func (r *ServiceRegistry) Register(service Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	name := service.Name()
	if _, exists := r.services[name]; exists {
		return errors.New("service already registered: " + name)
	}
	
	r.services[name] = service
	r.infos[name] = &ServiceInfo{
		Name:   name,
		Status: StatusStopped,
	}
	r.startOrder = append(r.startOrder, name)
	r.stopOrder = append([]string{name}, r.stopOrder...) // 反向順序
	
	return nil
}

// Unregister 取消註冊服務
func (r *ServiceRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	service, exists := r.services[name]
	if !exists {
		return errors.New("service not found: " + name)
	}
	
	if service.Status() == StatusRunning {
		return errors.New("cannot unregister running service: " + name)
	}
	
	delete(r.services, name)
	delete(r.infos, name)
	
	// 從啟動順序中移除
	newStartOrder := make([]string, 0, len(r.startOrder)-1)
	for _, n := range r.startOrder {
		if n != name {
			newStartOrder = append(newStartOrder, n)
		}
	}
	r.startOrder = newStartOrder
	
	// 從停止順序中移除
	newStopOrder := make([]string, 0, len(r.stopOrder)-1)
	for _, n := range r.stopOrder {
		if n != name {
			newStopOrder = append(newStopOrder, n)
		}
	}
	r.stopOrder = newStopOrder
	
	return nil
}

// StartAll 啟動所有服務
func (r *ServiceRegistry) StartAll(ctx context.Context) error {
	r.mu.RLock()
	order := make([]string, len(r.startOrder))
	copy(order, r.startOrder)
	r.mu.RUnlock()
	
	for _, name := range order {
		if err := r.Start(ctx, name); err != nil {
			// 如果啟動失敗，停止已啟動的服務
			r.StopAll(ctx)
			return err
		}
	}
	
	return nil
}

// StopAll 停止所有服務
func (r *ServiceRegistry) StopAll(ctx context.Context) error {
	r.mu.RLock()
	order := make([]string, len(r.stopOrder))
	copy(order, r.stopOrder)
	r.mu.RUnlock()
	
	var firstErr error
	for _, name := range order {
		if err := r.Stop(ctx, name); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	
	return firstErr
}

// Start 啟動指定服務
func (r *ServiceRegistry) Start(ctx context.Context, name string) error {
	r.mu.Lock()
	service, exists := r.services[name]
	info := r.infos[name]
	r.mu.Unlock()
	
	if !exists {
		return errors.New("service not found: " + name)
	}
	
	if info.Status == StatusRunning {
		return nil // 已經在運行
	}
	
	// 更新狀態
	r.updateStatus(name, StatusStarting, nil)
	
	// 啟動服務
	err := service.Start(ctx)
	if err != nil {
		r.updateStatus(name, StatusError, err)
		return err
	}
	
	r.updateStatus(name, StatusRunning, nil)
	r.mu.Lock()
	info.StartTime = time.Now()
	r.mu.Unlock()
	
	return nil
}

// Stop 停止指定服務
func (r *ServiceRegistry) Stop(ctx context.Context, name string) error {
	r.mu.Lock()
	service, exists := r.services[name]
	info := r.infos[name]
	r.mu.Unlock()
	
	if !exists {
		return errors.New("service not found: " + name)
	}
	
	if info.Status != StatusRunning {
		return nil // 不在運行
	}
	
	// 更新狀態
	r.updateStatus(name, StatusStopping, nil)
	
	// 停止服務
	err := service.Stop(ctx)
	if err != nil {
		r.updateStatus(name, StatusError, err)
		return err
	}
	
	r.updateStatus(name, StatusStopped, nil)
	
	return nil
}

// GetService 獲取服務
func (r *ServiceRegistry) GetService(name string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	service, exists := r.services[name]
	if !exists {
		return nil, errors.New("service not found: " + name)
	}
	
	return service, nil
}

// GetInfo 獲取服務資訊
func (r *ServiceRegistry) GetInfo(name string) (*ServiceInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	info, exists := r.infos[name]
	if !exists {
		return nil, errors.New("service not found: " + name)
	}
	
	// 返回副本
	infoCopy := *info
	return &infoCopy, nil
}

// ListServices 列出所有服務
func (r *ServiceRegistry) ListServices() []ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	infos := make([]ServiceInfo, 0, len(r.infos))
	for _, info := range r.infos {
		infos = append(infos, *info)
	}
	
	return infos
}

// updateStatus 更新服務狀態
func (r *ServiceRegistry) updateStatus(name string, status ServiceStatus, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if info, exists := r.infos[name]; exists {
		info.Status = status
		info.Error = err
	}
}

// WaitForService 等待服務達到指定狀態
func (r *ServiceRegistry) WaitForService(ctx context.Context, name string, status ServiceStatus, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := r.GetInfo(name)
			if err != nil {
				return err
			}
			if info.Status == status {
				return nil
			}
			if time.Now().After(deadline) {
				return errors.New("timeout waiting for service: " + name)
			}
		}
	}
}