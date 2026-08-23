package r001

import (
	"context"
	"errors"
	"fmt"
	httpadapter "github.com/example/alert-orchestration-engine/internal/adapter/http"
	"github.com/example/alert-orchestration-engine/internal/application"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/notifier"
	"testing"
	"time"
)

func TestR001Status(t *testing.T) {
	if httpadapter.StatusFor(fmt.Errorf("wrap: %w", httpadapter.ErrNotFound)) != 404 {
		t.Fatal("wrapped not found lost")
	}
}

type evRepo struct{ err error }

func (r evRepo) SaveEvent(context.Context, event.Event) error           { return r.err }
func (r evRepo) ListEvents(context.Context, int) ([]event.Event, error) { return nil, nil }

type clk struct{}

func (clk) Now() time.Time { return time.Unix(1, 0) }
func TestR001Ingest(t *testing.T) {
	want := errors.New("storage")
	s := application.NewService(evRepo{want}, nil, nil, nil, nil, nil, clk{}, nil)
	_, _, err := s.Ingest(context.Background(), event.Input{Source: "a", Name: "b"})
	if !errors.Is(err, want) {
		t.Fatal(err)
	}
}
func TestR001Parse(t *testing.T) {
	_, err := event.ParseOccurredAt("bad")
	var pe *time.ParseError
	if !errors.As(err, &pe) {
		t.Fatal(err)
	}
}
func TestR001Email(t *testing.T) {
	err := (notifier.EmailSender{From: "ops@example.com"}).Send(context.Background(), notification.Record{Target: "bad"})
	if !errors.Is(err, notifier.ErrInvalidEmailTarget) {
		t.Fatal(err)
	}
}
