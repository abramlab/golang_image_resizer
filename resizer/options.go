package resizer

type Option func(resizer *Resizer)

func WithResolution(width, height uint) Option {
	return func(r *Resizer) {
		r.width, r.height = width, height
	}
}

func WithWorkersNum(workers int) Option {
	return func(r *Resizer) {
		r.workersNum = workers
	}
}
