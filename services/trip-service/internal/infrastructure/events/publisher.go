package events

import "context"

type Publisher interface {
	Publish(ctx context.Context, event string) error
}
