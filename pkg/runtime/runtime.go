package runtime

type Interface interface {
	Factor3Load() error
}

func Load(loadable Interface) error {
	return loadable.Factor3Load()
}
