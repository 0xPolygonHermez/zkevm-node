package syncinterfaces

import "context"

type Halter interface {
	Halt(ctx context.Context, err error)
}
