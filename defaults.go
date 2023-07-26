package cog

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type getValue func(reflect.StructField) string

var tagHandlers = []getValue{
	environmentVariable("env"),
	defaultValue("default"),
}

func SetDefaults[T any](data *T) error {
	return setNested(reflect.ValueOf(data).Elem())
}

func environmentVariable(tag string) getValue {
	return func(sf reflect.StructField) string {
		if env := sf.Tag.Get(tag); env != "" {
			return os.Getenv(env)
		}

		return ""
	}
}

func defaultValue(tag string) getValue {
	return func(sf reflect.StructField) string {
		if val := sf.Tag.Get(tag); val != "" {
			return val
		}

		return ""
	}
}

func setNested(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			setNested(v.Field(i))
		} else {
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				if err := setField(t.Field(i), v.Field(i)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func setField(sf reflect.StructField, f reflect.Value) error {
	for _, getValue := range tagHandlers {
		if err := setValue(f, getValue(sf)); err != nil {
			return err
		}
	}

	return nil
}

func setValue(field reflect.Value, val string) error {
	if val == "" {
		return nil
	}

	if !field.CanSet() {
		return fmt.Errorf("can't set value: %s", val)
	}

	if !isEmpty(field) {
		// field already set.
		return nil
	}

	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.Atoi(val); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(val).Convert(field.Type()))
	case reflect.Bool:
		if val, err := strconv.ParseBool(val); err == nil {
			field.Set(reflect.ValueOf(bool(val)).Convert(field.Type()))
		}
	}

	return nil
}

func isEmpty(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}