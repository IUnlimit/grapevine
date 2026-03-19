package model

import "time"

type Service struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Suffix      string     `json:"suffix" gorm:"uniqueIndex;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	Endpoints   []Endpoint `json:"endpoints,omitempty" gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	Config      *Config    `json:"config,omitempty" gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Endpoint struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ServiceID uint      `json:"service_id" gorm:"index;not null"`
	Host      string    `json:"host" gorm:"not null"`
	Port      int       `json:"port" gorm:"not null"`
	Priority  int       `json:"priority" gorm:"default:0"` // 值越小优先级越高
	Status    string    `json:"status" gorm:"default:online;not null"` // online / offline
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Config struct {
	ID                    uint `json:"id" gorm:"primaryKey"`
	ServiceID             uint `json:"service_id" gorm:"uniqueIndex;not null"`
	RateLimit             int  `json:"rate_limit" gorm:"default:100"`  // QPS
	Burst                 int  `json:"burst" gorm:"default:200"`
	CircuitBreakThreshold int  `json:"circuit_break_threshold" gorm:"default:50"` // 错误率百分比
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
