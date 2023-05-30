package main

import (
	"errors"
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	val := reflect.ValueOf(out)

	if val.Kind() != reflect.Pointer {
		return errors.New("out is not a pointer")
	}
	if val.IsNil() {
		return errors.New("out is nil")
	}

	val = val.Elem()

	switch val.Kind() {
	case reflect.Struct:
		d, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type assertion failed: interface is %T, not map[string]interface{}", data)
		}

		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			dataValue, ok := d[val.Type().Field(i).Name]
			if !ok {
				continue
			}

			err := i2s(dataValue, valueField.Addr().Interface())
			if err != nil {
				return err
			}
		}

	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("type assertion failed: interface is %T, not []interface{}", data)
		}

		v := reflect.MakeSlice(val.Type(), len(d), len(d))
		for idx, elem := range d {
			err := i2s(elem, v.Index(idx).Addr().Interface())
			if err != nil {
				return err
			}
		}
		val.Set(v)

	case reflect.Bool:
		d, ok := data.(bool)
		if !ok {
			return fmt.Errorf("type assertion failed: interface is %T, not bool", data)
		}

		val.SetBool(d)

	case reflect.Int:
		d, ok := data.(float64)
		if !ok {
			return fmt.Errorf("type assertion failed: interface is %T, not float64", data)
		}

		val.SetInt(int64(d))

	case reflect.String:
		d, ok := data.(string)
		if !ok {
			return fmt.Errorf("type assertion failed: interface is %T, not string", data)
		}

		val.SetString(d)

	default:
		return fmt.Errorf("reflect type %s is not supported", val.Kind().String())
	}

	return nil
}
