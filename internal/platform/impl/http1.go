package impl

import (
	"context"
	"fmt"
)

// 第一个 HTTP 下游
type HTTPPlatform struct{}

func NewHTTPPlatform() *HTTPPlatform {
	return &HTTPPlatform{}
}

func (p *HTTPPlatform) Name() string {
	return "http"
}

func (p *HTTPPlatform) Deliver(ctx context.Context, target string, content interface{}) error {
	fmt.Printf("[HTTP Platform] Delivering to %s: %v\n", target, content)
	return nil
}

// 	Name() string
// 	Deliver(ctx context.Context, target string, content interface{}) error
