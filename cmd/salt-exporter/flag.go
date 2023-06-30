package main

import (
	"flag"
	"reflect"

	"github.com/spf13/viper"
)

type viperFlag struct {
	original flag.Flag
	alias    string
}

func (f viperFlag) HasChanged() bool { return true } // TODO: fix?

func (f viperFlag) Name() string {
	if f.alias != "" {
		return f.alias
	}
	return f.original.Name
}

func (f viperFlag) ValueString() string { return f.original.Value.String() }

func (f viperFlag) ValueType() string {
	t := reflect.TypeOf(f.original.Value)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Kind().String()
	}
	return t.Kind().String()
}

type viperFlagSet struct {
	flags []viperFlag
}

func (f viperFlagSet) VisitAll(fn func(viper.FlagValue)) {
	for _, flag := range f.flags {
		fn(flag)
	}
}
