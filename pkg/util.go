package pkg

import (
	"iter"
	"net"
	"time"
)

func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

type DefaultTime struct{}

func (d *DefaultTime) Now() time.Time {
	return time.Now()
}

func RemovePort(hostname string) string {
	host, _, err := net.SplitHostPort(hostname)
	if err != nil {
		return hostname
	}

	return host
}
