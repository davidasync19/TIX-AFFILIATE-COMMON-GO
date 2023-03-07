package toggle

type ToggleExecutor[T any, V any, W any] interface {
	IsToggleOn(w W) bool
	OnToggleOn(t T) V
	OnToggleOff(t T) V
}

type ToggleHelper[T any, V any, W any] struct {
	Executor ToggleExecutor[T, V, W]
}

func (e ToggleHelper[T, V, W]) Run(t T, w W) V {
	if e.Executor.IsToggleOn(w) {
		return e.Executor.OnToggleOn(t)
	}
	return e.Executor.OnToggleOff(t)
}
