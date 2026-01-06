package platform

import (
	"context"
	"fmt"
)

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
