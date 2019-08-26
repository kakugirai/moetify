package app

type RedisStorage interface {
	Shorten(url string, exp int64) (string, error)
	ShortLinkInfo(eid string) (interface{}, error)
	Unshorten(eid string) (string, error)
}