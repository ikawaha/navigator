package service

import (
	"context"
	"testing"
)

func TestService_StartStop(t *testing.T) {
	t.Run("command empty", func(t *testing.T) {
		s := New("localhost", []string{})
		if err := s.Start(context.Background(), false); err != nil {
			if got, want := err.Error(), "failed to parse command: empty command"; got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		} else {
			t.Errorf("expected error, but nil")
		}
	})
	t.Run("command already running", func(t *testing.T) {
		s := New("localhost", []string{"sleep", "1"})
		if err := s.Start(context.Background(), false); err != nil {
			t.Fatal("unexpected error", err)
		}
		if err := s.Start(context.Background(), false); err != nil {
			if got, want := err.Error(), "already running"; got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		} else {
			t.Errorf("expected error, but nil")
		}
	})
	t.Run("service start and stop", func(t *testing.T) {
		s := New("localhost", []string{"sleep", "1"})

		// start
		if err := s.Start(context.Background(), false); err != nil {
			t.Errorf("s.Start() faild: unexpected error %v", err)
			return
		}
		if s.command == nil {
			t.Errorf("expected command is not nil")
		}
		if s.baseURL == "" {
			t.Errorf("expected base URL is not empty")
		}

		// stop
		if err := s.Stop(); err != nil {
			t.Errorf("s.Stop() failed: unexpected error %v", err)
		}
		if got := s.command; got != nil {
			t.Errorf("expected command is nil, but %p", got)
		}
		if got := s.baseURL; got != "" {
			t.Errorf("expected base URL is empty, but %q", got)
		}
	})
	t.Run("command already stopped", func(t *testing.T) {
		s := New("localhost", []string{"sleep", "1"})
		if err := s.Start(context.Background(), false); err != nil {
			t.Errorf("s.Start() faild: unexpected error %v", err)
			return
		}
		if err := s.Stop(); err != nil {
			t.Errorf("s.Stop() failed: unexpected error %v", err)
		}
		if err := s.Stop(); err != nil {
			if got, want := err.Error(), "already stopped"; got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		} else {
			t.Errorf("expected error, but nil")
		}
	})
}
