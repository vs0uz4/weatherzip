package service

import "time"

type UptimeService interface {
	GetUptime() string
}

type uptimeService struct {
	startTime time.Time
}

func NewUptimeService() UptimeService {
	return &uptimeService{
		startTime: time.Now(),
	}
}

func (u *uptimeService) GetUptime() string {
	return time.Since(u.startTime).String()
}
