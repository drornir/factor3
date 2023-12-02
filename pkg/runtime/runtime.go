package runtime

type Interface interface {
	Factor3Load(argv []string) error
}

func Load(loadable Interface, argv []string) error {
	return loadable.Factor3Load(argv)
}
