package impl

import (
	"context"
	"fmt"
)

// 短信平台下游
type SMSPlatform struct{}

func NewSMSPlatform() *SMSPlatform {
	return &SMSPlatform{}
}

func (p *SMSPlatform) Name() string {
	return "sms"
}

func (p *SMSPlatform) Deliver(ctx context.Context, target string, content interface{}) error {
	fmt.Printf("[SMS Platform] Sending SMS to %s: %v\n", target, content)
	return nil
}
