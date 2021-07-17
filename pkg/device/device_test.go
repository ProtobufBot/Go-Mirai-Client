package device

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewDevice(t *testing.T) {
	now := time.Now().UnixNano()
	d1 := RandDevice(rand.New(rand.NewSource(now)))
	d2 := RandDevice(rand.New(rand.NewSource(now)))
	t.Log(string(d1.ToJson()))
	t.Log(string(d2.ToJson()))
	t.Log(string(d1.ToJson()) == string(d2.ToJson()))
}
