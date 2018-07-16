package utils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type key struct {
}

type val struct {
	fromInside bool
}

func TestMyList(t *testing.T) {
	b := context.Background()
	c1 := context.WithValue(b, key{}, &val{})
	c2 := context.WithValue(b, key{}, &val{})

	go funcA(c1, 10)
	go funcA(c2, 20)

	time.Sleep(30000)
}

func funcA(ctx context.Context, x int) {
	fmt.Printf("x=%d\n", x)
	value := ctx.Value(key{}).(*val)
	if !value.fromInside {
		value.fromInside = true
		funcA(ctx, x+1)
	}
}
