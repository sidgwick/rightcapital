package platform

import (
	"context"
)

type Platform interface {
	Name() string
	Deliver(ctx context.Context, target string, content interface{}) error
}

type PlatformManager struct {
	platforms map[string]Platform
}

func NewPlatformManager() *PlatformManager {
	return &PlatformManager{
		platforms: make(map[string]Platform),
	}
}

func (pm *PlatformManager) RegisterPlatform(p Platform) {
	pm.platforms[p.Name()] = p
}

func (pm *PlatformManager) GetPlatform(name string) (Platform, bool) {
	p, ok := pm.platforms[name]
	return p, ok
}
