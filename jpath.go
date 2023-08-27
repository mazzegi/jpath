package jpath

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNotFound = fmt.Errorf("not-found")
	ErrBadArgs  = fmt.Errorf("bad-args")
)

func isRVZero(rv reflect.Value) bool {
	zv := reflect.Value{}
	return rv == zv
}

func queryValue(in any, spath string) (reflect.Value, error) {
	path := strings.Split(spath, "/")
	crv := reflect.ValueOf(in)
	for _, elt := range path {
		if elt == "" {
			//skip empty
			continue
		}
		rty := crv.Type()
		kind := rty.Kind()
		if kind == reflect.Pointer {
			crv = crv.Elem()
			kind = crv.Kind()
		}

		switch kind {
		case reflect.Struct:
			crv = crv.FieldByName(elt)
			if isRVZero(crv) {
				return reflect.Value{}, errors.Join(ErrNotFound, fmt.Errorf("no such struct field %q", elt))
			}
		case reflect.Slice:
			ix, err := strconv.ParseInt(elt, 10, 64)
			if err != nil {
				return reflect.Value{}, errors.Join(ErrBadArgs, fmt.Errorf("cannot parse %q as int for slice index: %w", elt, err))
			}
			if ix < 0 {
				return reflect.Value{}, errors.Join(ErrBadArgs, fmt.Errorf("invalid slice index %d", ix))
			}
			if ix >= int64(crv.Len()) {
				return reflect.Value{}, errors.Join(ErrBadArgs, fmt.Errorf("invalid slice index %d ( >= len=%d)", ix, crv.Len()))
			}
			crv = crv.Index(int(ix))
		case reflect.Map:
			if crv.Type().Key().Kind() != reflect.String {
				return reflect.Value{}, fmt.Errorf("cannot query non-string map keys. map keys are %s", crv.Type().Kind().String())
			}
			crv = crv.MapIndex(reflect.ValueOf(elt))
			if isRVZero(crv) {
				return reflect.Value{}, errors.Join(ErrNotFound, fmt.Errorf("no such map key %q", elt))
			}
		default:
			return reflect.Value{}, errors.Join(ErrNotFound, fmt.Errorf("cannot query reflect-kind %T", kind))
		}
	}
	return crv, nil
}

func Query(in any, spath string) (any, error) {
	rv, err := queryValue(in, spath)
	if err != nil {
		return nil, fmt.Errorf("query-value: %w", err)
	}
	return rv.Interface(), nil
}

func Set(in any, spath string, value any) error {
	if reflect.TypeOf(in).Kind() != reflect.Pointer {
		return fmt.Errorf("cannot set non-pointer type")
	}
	rv, err := queryValue(in, spath)
	if err != nil {
		return fmt.Errorf("query-value: %w", err)
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if !rv.CanSet() {
		return fmt.Errorf("cannot set %v (%s)", rv.String(), rv.Type())
	}
	setVal := reflect.ValueOf(value)
	if !setVal.CanConvert(rv.Type()) {
		return fmt.Errorf("cannot convert value of type %s to %s", setVal.Type().String(), rv.Type().String())
	}
	setValConv := setVal.Convert(rv.Type())

	rv.Set(setValConv)
	return nil
}
