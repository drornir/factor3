package factor3

type Config any

func Load[C Config]() (C, error) {
	var conf C
	return conf, nil
}
