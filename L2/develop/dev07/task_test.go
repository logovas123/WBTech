package main

import (
	"testing"
	"time"
)

type TestCase struct {
	SliceOfChanels []<-chan interface{}
	ExpectedTime   time.Duration
}

func TestGetCurrentTime(t *testing.T) {
	cases := []TestCase{
		{
			SliceOfChanels: []<-chan interface{}{
				sig(2 * time.Hour),
				sig(5 * time.Minute),
				sig(1 * time.Second),
				sig(1 * time.Hour),
				sig(1 * time.Minute),
			},
			ExpectedTime: 1 * time.Second,
		},
	}

	for caseNum, item := range cases {
		start := time.Now()

		<-or(item.SliceOfChanels...)

		if time.Since(start)-item.ExpectedTime > 500*time.Millisecond {
			t.Errorf("[%d] incorrect time", caseNum)
		}

	}
}
