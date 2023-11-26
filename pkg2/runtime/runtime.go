package runtime

type Interface interface {
	Factor3Load() error
}

func Load(config Interface) error {
	return config.Factor3Load()
}
