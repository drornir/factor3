package config

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/drornir/factor3/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Loader struct {
	viper    *viper.Viper
	pflagset *pflag.FlagSet

	jpath []string
	fpath []string

	loaders []func() error
	boundTo *any

	lock *sync.RWMutex
}

// Bind creates a loader from `into`, which should be a pointer to a struct.
//
// `viper` must be non-nil, because you should configure your viper instance
// to.
//
// `pflagset` can be nil, but that turns off flags support.
// To enable it, you need to pass an initialized `*pflag.FlagSet`.
// Use `cobraRootCmd.Flags()` if using cobra, or `pflag.CommandLine` unless you have a reason not to.
// There are two mechanism to register pflags:
//
//  1. Create a new flagset using spf13/pflag.NewFlagSet(), and register flags manually
//     before calling Bind(). Equivalently, you can use pflag's global and pass pflag.CommandLine here
//
//  2. Create a new flagset, but don't bind it directly. Annotated toy struct fields with
//     `flag:"flag-name"` and Bind() will discover it and register it on you pflagset
//
// In both cases, viper.BindPFlags() will be called on `pflagset` before returning from this function
func Bind(into any, viper *viper.Viper, pflagset *pflag.FlagSet) (*Loader, error) {
	l := newLoader(viper, pflagset)
	if err := l.bind(into); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Loader) Load() error {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if l.boundTo == nil {
		return fmt.Errorf("Bind() needs to be called before calling Load() for the first time")
	}

	var errs []error
	for _, loader := range l.loaders {
		err := loader()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return LoadError{Errs: errs}
	}
	return nil
}

func newLoader(viper *viper.Viper, pflagset *pflag.FlagSet) *Loader {
	return &Loader{
		viper:    viper,
		pflagset: pflagset,
		lock:     &sync.RWMutex{},
	}
}

func (p *Loader) bind(into any) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	reflected := reflect.ValueOf(into)
	if reflected.Kind() != reflect.Pointer || reflected.Elem().Kind() != reflect.Struct {
		return ParseError{Err: fmt.Errorf("input must be a pointer to struct, got %s", reflected.Type()), Value: into}
	}

	p.boundTo = &into
	if err := p.visit(reflected.Elem()); err != nil {
		return ParseError{Err: err, Value: into}
	}
	if p.pflagset != nil && p.viper != nil {
		p.viper.BindPFlags(p.pflagset)
	}

	return nil
}

func (p *Loader) visit(v reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,

		reflect.String,

		reflect.Map,
		reflect.Slice, reflect.Array:

		p.registerPflag(v)
		p.registerViper(v.Addr())

		return nil

	case reflect.Pointer: // should i have pointers?
		if v.IsNil() {
			// fmt.Fprintln(os.Stderr, "v is nil", v)
			v.Set(reflect.New(v.Type().Elem()))
		}
		if err := p.visit(v.Elem()); err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		for i := v.NumField() - 1; i >= 0; i-- {
			f := v.Type().Field(i)
			p.jpath = append(p.jpath, toJSONName(f))
			p.fpath = append(p.fpath, f.Tag.Get("flag"))
			vv := v.Field(i)
			if err := p.visit(vv); err != nil {
				return err
			}
			p.jpath = p.jpath[:len(p.jpath)-1]
			p.fpath = p.fpath[:len(p.fpath)-1]
		}
		return nil
	default:
		return p.errWithContext("value cannot be visited", v, strings.Join(p.jpath, "."))
	}
}

func (l *Loader) registerPflag(v reflect.Value) {
	if l.pflagset == nil {
		return // TODO log WARN
	}

	flagsPath := l.fpathString()
	if flagsPath == "" {
		return // TODO log DEBUG
	}

	switch v.Type().Kind() {
	case reflect.Bool:
		c := v.Interface().(bool)
		l.pflagset.BoolVar(&c, flagsPath, c, "")
	case reflect.Int:
		i := v.Interface().(int)
		l.pflagset.IntVar(&i, flagsPath, i, "")
	case reflect.Int8:
		i := v.Interface().(int8)
		l.pflagset.Int8Var(&i, flagsPath, i, "")
	case reflect.Int16:
		i := v.Interface().(int16)
		l.pflagset.Int16Var(&i, flagsPath, i, "")
	case reflect.Int32:
		i := v.Interface().(int32)
		l.pflagset.Int32Var(&i, flagsPath, i, "")
	case reflect.Int64:
		i := v.Interface().(int64)
		l.pflagset.Int64Var(&i, flagsPath, i, "")
	case reflect.Uint:
		i := v.Interface().(uint)
		l.pflagset.UintVar(&i, flagsPath, i, "")
	case reflect.Uint8:
		i := v.Interface().(uint8)
		l.pflagset.Uint8Var(&i, flagsPath, i, "")
	case reflect.Uint16:
		i := v.Interface().(uint16)
		l.pflagset.Uint16Var(&i, flagsPath, i, "")
	case reflect.Uint32:
		i := v.Interface().(uint32)
		l.pflagset.Uint32Var(&i, flagsPath, i, "")
	case reflect.Uint64:
		i := v.Interface().(uint64)
		l.pflagset.Uint64Var(&i, flagsPath, i, "")
	case reflect.Float32:
		f := v.Interface().(float32)
		l.pflagset.Float32Var(&f, flagsPath, f, "")
	case reflect.Float64:
		f := v.Interface().(float64)
		l.pflagset.Float64Var(&f, flagsPath, f, "")
	case reflect.String:
		s := v.Interface().(string)
		l.pflagset.StringVar(&s, flagsPath, s, "")
	}
}

func (l *Loader) registerViper(vAddr reflect.Value) {
	viperPath := strings.Join(l.jpath, ".")
	loader := func() error {
		// if !l.viper.IsSet(viperPath) {
		// 	fmt.Fprintln(os.Stderr, "vAddr", vAddr, "T", vAddr.Type())
		// 	vAddr.Set(reflect.Zero(vAddr.Type()))
		// 	return nil
		// }

		log.GG().D(context.Background(), "loading viper value", "path", viperPath)
		untypedVal := l.viper.Get(viperPath)
		if untypedVal == nil {
			log.GG().D(context.TODO(), "value is nil", "path", viperPath)
			return nil
		}
		log.GG().D(context.Background(), "loaded viper value", "path", viperPath)
		if err := unmarshalViper(vAddr, untypedVal); err != nil {
			return l.errWithContext(err.Error(), vAddr.Elem(), viperPath)
		}
		return nil
	}
	l.loaders = append(l.loaders, loader)
}

func (l *Loader) errWithContext(msg string, v reflect.Value, jsonPath string) error {
	// TODO make this error really nice
	return fmt.Errorf("in json path %q and value %q : %s", jsonPath, v, msg)
}

func (l *Loader) fpathString() string {
	fp := append([]string(nil), l.fpath...)
	fp = slices.DeleteFunc(fp, func(s string) bool { return s == "" })
	return strings.Join(fp, "-")
}

func unmarshalViper(into reflect.Value, data any) error {
	// TODO check if I can reuse something from viper
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("viper data is not json serializale: %w", err)
	}

	err = json.Unmarshal(jsonBytes, into.Interface())
	if err != nil {
		return fmt.Errorf("unable to parse data %s into type %s: %w", jsonBytes, into.Type(), err)
	}
	return nil
}

func toJSONName(f reflect.StructField) string {
	jsonName := f.Tag.Get("json")
	jsonName, _, _ = strings.Cut(jsonName, ",")
	if jsonName == "" {
		jsonName = f.Name
	}
	return jsonName
}
