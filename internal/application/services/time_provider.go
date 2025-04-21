package services

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func NewRealTimeProvider() *RealTimeProvider {
	return &RealTimeProvider{}
}

func (p *RealTimeProvider) Now() time.Time {
	return time.Now()
}
