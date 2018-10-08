package api

import (
	"log"

	"github.com/manifoldco/go-manifold"
)

// FeatureMap is a wrap around manifold.FeatureMap with some helpers for the workshop
type FeatureMap manifold.FeatureMap

// GetInt tries to get an int feature with the key passed. If the feature isn't an int, it panics.
func (fm FeatureMap) GetInt(key string) int {
	val, ok := fm[key]
	if !ok {
		return 0
	}

	switch t := val.(type) {
	case float64:
		return int(val.(float64))
	case int64:
		return int(val.(int64))
	default:
		log.Fatalf("invalid type %s", t)
		return 0
	}
}

// GetBool tries to get a boolean feature with the key passed. If the feature isn't a boolean, it panics.
func (fm FeatureMap) GetBool(key string) bool {
	val, ok := fm[key]
	if !ok {
		return false
	}

	b, ok := val.(bool)
	if !ok {
		log.Fatalf("invalid type %v", val)
		return false
	}

	return b
}

// GetString tries to get a string feature with the key passed. If the feature isn't a string, it panics.
func (fm FeatureMap) GetString(key string) string {
	val, ok := fm[key]
	if !ok {
		return ""
	}

	s, ok := val.(string)
	if !ok {
		log.Fatalf("invalid type %v", val)
		return ""
	}

	return s
}
