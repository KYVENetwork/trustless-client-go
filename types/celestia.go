package types

import "encoding/json"

type CelestiaValue struct {
	SharesByNamespace []NamespacedShares `json:"sharesByNamespace"`
}

type NamespacedShares struct {
	Data        json.RawMessage `json:"data"`
	NamespaceId string          `json:"namespace_id"`
}
