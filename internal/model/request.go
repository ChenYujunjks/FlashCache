package model

type SetCacheRequest struct {
	Value      string `json:"value"`
	TTLSeconds int    `json:"ttl_seconds"`
}
