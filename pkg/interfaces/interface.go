package interfaces

type Service[P any, R any] interface {
	Invoke(P) (R, error)
}
