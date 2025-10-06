package morphotactics

import (
	"strings"
	"sync"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// AttributeToSurfaceCache caches surface forms for phonetic attributes
type AttributeToSurfaceCache struct {
	attributeMap map[string]string
	lock         sync.RWMutex
}

// NewAttributeToSurfaceCache creates a new cache
func NewAttributeToSurfaceCache() *AttributeToSurfaceCache {
	return &AttributeToSurfaceCache{
		attributeMap: make(map[string]string),
	}
}

// AddSurface adds a surface form for the given attributes
func (c *AttributeToSurfaceCache) AddSurface(attributes map[turkish.PhoneticAttribute]bool, surface string) {
	key := c.getKey(attributes)
	c.lock.Lock()
	c.attributeMap[key] = surface
	c.lock.Unlock()
}

// GetSurface retrieves a surface form for the given attributes
func (c *AttributeToSurfaceCache) GetSurface(attributes map[turkish.PhoneticAttribute]bool) string {
	key := c.getKey(attributes)
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.attributeMap[key]
}

func (c *AttributeToSurfaceCache) getKey(attributes map[turkish.PhoneticAttribute]bool) string {
	var keyParts []string
	for attr := range attributes {
		keyParts = append(keyParts, attr.GetStringForm())
	}
	// Sort for consistent hashing
	return strings.Join(keyParts, "_")
}
