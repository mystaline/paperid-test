package interfaces

import "context"

type Service[P any, R any] interface {
	Invoke(context.Context, P) (R, error)
}
