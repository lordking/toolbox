package common

import (
	"fmt"
	"time"
)

type (
	Timer struct {
		ticker *time.Ticker
	}
)

func (t *Timer) Start() {

	t.ticker = time.NewTicker(time.Second * 1)

	fmt.Print("xx")

	go func() {
		for _ = range t.ticker.C {
			fmt.Printf("ticked at %v", time.Now())
		}
	}()
}

func (t *Timer) Stop() {
	t.ticker.Stop()
}

func NewTimer() *Timer {
	return &Timer{}
}
