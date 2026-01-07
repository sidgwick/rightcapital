package impl

import (
	"context"
	"fmt"
)

// 第二个 HTTP 下游
type HTTP2_Platform struct{}

func NewHTTP2_Platform() *HTTP2_Platform {
	return &HTTP2_Platform{}
}

func (p *HTTP2_Platform) Name() string {
	return "http2"
}

func (p *HTTP2_Platform) Deliver(ctx context.Context, target string, content interface{}) error {
	fmt.Printf("[HTTP2 Platform] Delivering to %s: %v\n", target, content)
	return nil
}
