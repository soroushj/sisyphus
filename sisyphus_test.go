package sisyphus_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/soroushj/sisyphus"
)

var (
	errDummy         = errors.New("dummy")
	errUnrecoverable = errors.New("unrecoverable")
)

func TestCtxCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := sisyphus.Do(ctx, func() error { return nil })
	if err != context.Canceled {
		t.Error(err)
	}
}

func TestCtxDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	s := sisyphus.New(1*time.Millisecond, 30*time.Millisecond)
	err := s.Do(ctx, func() error { return errDummy })
	if err != context.DeadlineExceeded {
		t.Error(err)
	}
}

func TestImmediateSuccess(t *testing.T) {
	ctx := context.Background()
	err := sisyphus.Do(ctx, func() error { return nil })
	if err != nil {
		t.Error(err)
	}
}

func TestImmediateFailure(t *testing.T) {
	ctx := context.Background()
	err := sisyphus.DoIf(ctx, func() error { return errDummy }, func(error) bool { return false })
	if err != errDummy {
		t.Error(err)
	}
}

func TestEventualSuccess(t *testing.T) {
	ctx := context.Background()
	s := sisyphus.New(1*time.Millisecond, 30*time.Millisecond)
	n := 0
	err := s.DoIf(ctx, func() error {
		n++
		if n < 2 {
			return errDummy
		}
		return nil
	}, func(error) bool {
		return true
	})
	if err != nil {
		t.Error(err)
	}
}

func TestEventualFailure(t *testing.T) {
	ctx := context.Background()
	s := sisyphus.New(1*time.Millisecond, 30*time.Millisecond)
	n := 0
	err := s.DoIf(ctx, func() error {
		n++
		if n < 2 {
			return errDummy
		}
		return errUnrecoverable
	}, func(err error) bool {
		return err != errUnrecoverable
	})
	if err != errUnrecoverable {
		t.Error(err)
	}
}

func TestBadConfig(t *testing.T) {
	ctx := context.Background()
	s := sisyphus.New(0, -1)
	n := 0
	err := s.DoIf(ctx, func() error {
		n++
		if n < 2 {
			return errDummy
		}
		return nil
	}, func(error) bool {
		return true
	})
	if err != nil {
		t.Error(err)
	}
}
