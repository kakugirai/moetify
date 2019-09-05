package model

// RedisStorage defined redis storage interface
type RedisStorage interface {
	Shorten(url string, exp int64) (string, error)
	ShortLinkInfo(eid string) (*DetailInfo, error)
	// Unshorten(eid string) (string, error)
}
