package vp

import (
	"context"
	"strconv"
	"time"

	"github.com/bbars/whispar/pkg/vpencoding"
)

type Time time.Time

func (t Time) MarshalVp(ctx context.Context) ([]byte, error) {
	if time.Time(t).IsZero() {
		return vpencoding.Marshal(ctx, int64(0))
	} else {
		return vpencoding.Marshal(ctx, time.Time(t).UnixMilli())
	}
}

type TimeString time.Time

func (t TimeString) MarshalVp(ctx context.Context) ([]byte, error) {
	if time.Time(t).IsZero() {
		return vpencoding.Marshal(ctx, "0")
	} else {
		return vpencoding.Marshal(ctx, strconv.FormatInt(time.Time(t).UnixMilli(), 10))
	}
}
