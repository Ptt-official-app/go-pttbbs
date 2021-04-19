package bbs

import "github.com/Ptt-official-app/go-pttbbs/cache"

func ReloadUHash() (err error) {
	return cache.LoadUHash()
}
