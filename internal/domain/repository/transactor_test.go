package repository

import (
	"context"
	"errors"
	"testing"
)

type fakeTransactor struct {
	committed bool
	rolled    bool
}

func (f *fakeTransactor) WithinTx(ctx context.Context, fn func(context.Context, Repositories) error) error {
	if err := fn(ctx, Repositories{}); err != nil {
		f.rolled = true
		return err
	}
	f.committed = true
	return nil
}

func TestTransactorCommitsOnSuccess(t *testing.T) {
	tr := &fakeTransactor{}
	if err := tr.WithinTx(context.Background(), func(context.Context, Repositories) error { return nil }); err != nil {
		t.Fatalf("WithinTx error: %v", err)
	}
	if !tr.committed || tr.rolled {
		t.Fatalf("committed=%v rolled=%v", tr.committed, tr.rolled)
	}
}

func TestTransactorRollsBackOnError(t *testing.T) {
	tr := &fakeTransactor{}
	want := errors.New("boom")
	if err := tr.WithinTx(context.Background(), func(context.Context, Repositories) error { return want }); !errors.Is(err, want) {
		t.Fatalf("WithinTx error = %v", err)
	}
	if tr.committed || !tr.rolled {
		t.Fatalf("committed=%v rolled=%v", tr.committed, tr.rolled)
	}
}
