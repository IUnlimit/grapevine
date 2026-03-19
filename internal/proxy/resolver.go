package proxy

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
)

type Resolver struct {
	mu    sync.RWMutex
	cache map[string][]model.Endpoint // suffix -> sorted endpoints
}

var DefaultResolver = &Resolver{
	cache: make(map[string][]model.Endpoint),
}

// Resolve 返回指定 suffix 的最高优先级在线端点
func (r *Resolver) Resolve(suffix string) (*model.Endpoint, error) {
	var svc model.Service
	if err := database.DB.Where("suffix = ? AND is_active = ?", suffix, true).First(&svc).Error; err != nil {
		return nil, fmt.Errorf("service not found for suffix: %s", suffix)
	}

	var endpoints []model.Endpoint
	if err := database.DB.Where("service_id = ? AND status = ?", svc.ID, "online").
		Order("priority ASC").Find(&endpoints).Error; err != nil {
		return nil, fmt.Errorf("no endpoints found for suffix: %s", suffix)
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no online endpoints for suffix: %s", suffix)
	}

	return &endpoints[0], nil
}

// ResolveAll 返回指定 suffix 的所有在线端点（按优先级排序）
func (r *Resolver) ResolveAll(suffix string) ([]model.Endpoint, error) {
	var svc model.Service
	if err := database.DB.Where("suffix = ? AND is_active = ?", suffix, true).First(&svc).Error; err != nil {
		return nil, fmt.Errorf("service not found for suffix: %s", suffix)
	}

	var endpoints []model.Endpoint
	if err := database.DB.Where("service_id = ? AND status = ?", svc.ID, "online").
		Order("priority ASC").Find(&endpoints).Error; err != nil {
		return nil, err
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].Priority < endpoints[j].Priority
	})

	return endpoints, nil
}

// ResolveWithFallback 尝试主端点，失败时返回备选
func (r *Resolver) ResolveWithFallback(suffix string, excludeIDs []uint) (*model.Endpoint, error) {
	endpoints, err := r.ResolveAll(suffix)
	if err != nil {
		return nil, err
	}

	excludeSet := make(map[uint]bool)
	for _, id := range excludeIDs {
		excludeSet[id] = true
	}

	for i := range endpoints {
		if !excludeSet[endpoints[i].ID] {
			return &endpoints[i], nil
		}
	}

	return nil, fmt.Errorf("all endpoints exhausted for suffix: %s", suffix)
}

// UpdateEndpointStatus 更新端点状态
func UpdateEndpointStatus(endpointID uint, status string) {
	if err := database.DB.Model(&model.Endpoint{}).Where("id = ?", endpointID).
		Update("status", status).Error; err != nil {
		log.Printf("failed to update endpoint %d status: %v", endpointID, err)
	}
}
