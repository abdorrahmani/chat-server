package server

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	refill := int(elapsed / tb.refillRate)

	if refill > 0 {
		tb.tokens = min(tb.maxTokens, tb.tokens+refill)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}
