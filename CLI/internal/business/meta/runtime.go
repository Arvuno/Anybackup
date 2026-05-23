package meta

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
import "github.com/spf13/pflag"

type Runtime struct {
	Remote     inputs.RemoteOptions
	StringMap  map[string]string
	Inline     string
	BodyFile   string
	Body       []byte
	stringRefs map[string]*string
	intRefs    map[string]*int
	boolRefs   map[string]*bool
	sliceRefs  map[string]*[]string
	changed    map[string]struct{}
}

func NewRuntime() *Runtime {
	return &Runtime{
		StringMap:  map[string]string{},
		stringRefs: map[string]*string{},
		intRefs:    map[string]*int{},
		boolRefs:   map[string]*bool{},
		sliceRefs:  map[string]*[]string{},
		changed:    map[string]struct{}{},
	}
}

func (rt *Runtime) BindStringFlag(cmd FlagBinder, name string, defaultValue string, usage string) {
	value := defaultValue
	cmd.StringVar(&value, name, defaultValue, usage)
	rt.stringRefs[name] = &value
}

func (rt *Runtime) String(name string) string {
	if ref, ok := rt.stringRefs[name]; ok && ref != nil {
		return *ref
	}
	return rt.StringMap[name]
}

func (rt *Runtime) SetString(name string, value string) {
	rt.StringMap[name] = value
}

func (rt *Runtime) SetInt(name string, value int) {
	rt.intRefs[name] = &value
}

func (rt *Runtime) SetBool(name string, value bool) {
	rt.boolRefs[name] = &value
}

func (rt *Runtime) BindIntFlag(cmd IntFlagBinder, name string, defaultValue int, usage string) {
	value := defaultValue
	cmd.IntVar(&value, name, defaultValue, usage)
	rt.intRefs[name] = &value
}

func (rt *Runtime) BindBoolFlag(cmd BoolFlagBinder, name string, defaultValue bool, usage string) {
	value := defaultValue
	cmd.BoolVar(&value, name, defaultValue, usage)
	rt.boolRefs[name] = &value
}

func (rt *Runtime) BindStringSliceFlag(cmd StringSliceFlagBinder, name string, defaultValue []string, usage string) {
	value := defaultValue
	cmd.StringSliceVar(&value, name, defaultValue, usage)
	rt.sliceRefs[name] = &value
}

func (rt *Runtime) Int(name string) int {
	if ref, ok := rt.intRefs[name]; ok && ref != nil {
		return *ref
	}
	return 0
}

func (rt *Runtime) Bool(name string) bool {
	if ref, ok := rt.boolRefs[name]; ok && ref != nil {
		return *ref
	}
	return false
}

func (rt *Runtime) Strings(name string) []string {
	if ref, ok := rt.sliceRefs[name]; ok && ref != nil {
		return *ref
	}
	return nil
}

func (rt *Runtime) Changed(name string) bool {
	_, ok := rt.changed[name]
	return ok
}

func (rt *Runtime) MarkChanged(name string) {
	rt.changed[name] = struct{}{}
}

func (rt *Runtime) CaptureChangedFlags(flags *pflag.FlagSet) {
	rt.changed = map[string]struct{}{}

	if flags == nil {
		return
	}
	flags.Visit(func(f *pflag.Flag) {
		rt.MarkChanged(f.Name)
	})
}

type FlagBinder interface {
	StringVar(p *string, name string, value string, usage string)
}

type IntFlagBinder interface {
	IntVar(p *int, name string, value int, usage string)
}

type BoolFlagBinder interface {
	BoolVar(p *bool, name string, value bool, usage string)
}

type StringSliceFlagBinder interface {
	StringSliceVar(p *[]string, name string, value []string, usage string)
}
