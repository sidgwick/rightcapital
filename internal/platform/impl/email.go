package impl

import (
	"context"
	"fmt"
)

// Email 下游
type EmailPlatform struct{}

func NewEmailPlatform() *EmailPlatform {
	return &EmailPlatform{}
}

func (p *EmailPlatform) Name() string {
	return "email"
}

func (p *EmailPlatform) Deliver(ctx context.Context, target string, content interface{}) error {
	fmt.Printf("[Email Platform] Sending email to %s: %v\n", target, content)
	return nil
}
