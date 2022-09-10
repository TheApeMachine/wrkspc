package datura

import (
	"testing"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

func BenchmarkWrite(b *testing.B) {
	spd.InitCache()
	errnie.Tracing(false)
	store := NewS3()

	for i := 0; i < b.N; i++ {
		store.Write(spd.NewCached(
			"datapoint", "test", "test.wrkspc.org",
		))
	}
}
