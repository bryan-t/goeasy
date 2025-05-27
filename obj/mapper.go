package obj

import (
	"fmt"
	"reflect"
)

// ErrMismatchType returned when field of source can't be mapped to destination due to mismatched types.
var ErrMismatchType error = fmt.Errorf("type mismatch")
var ErrInsufficientCapacity error = fmt.Errorf("insufficient capacity")
var ErrNotAddresable error = fmt.Errorf("not addressable")

// AI generated code start
var ErrFieldNotFound error = fmt.Errorf("field not found")

// AI generated code end

type Mapper struct {
	cfg mapperConfig
}

// NewMapper creates a new instance of Mapper
func NewMapper() *Mapper {
	return &Mapper{
		cfg: mapperConfig{
			//fieldMaps: make(map[structMapKey]map[string]*FieldMapConfig),
			converters: make(map[structMapKey]func(src any) (dst any, err error)),
		},
	}
}

// Map copies src field values to dst fields. Fields must have the same name.
// Sample usage:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/bryan-t/goeasy/obj"
//	)
//
//	type UserDTO struct {
//		ID   int
//		Name string
//	}
//
//	type User struct {
//		ID   int
//		Name string
//	}
//
//	func main() {
//		dto := UserDTO{
//			ID:   1,
//			Name: "John",
//		}
//
//		user := User{}
//		mapper := obj.NewMapper()
//		err := mapper.Map(dto, &user)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Printf("user: %+v\n", user)
//	}
func (m *Mapper) Map(src interface{}, dst interface{}) error {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)
	if dstValue.Type().Kind() == reflect.Pointer {
		dstValue = dstValue.Elem()
	}
	if !dstValue.CanAddr() {
		return ErrNotAddresable
	}
	return m.mapValue(srcValue, dstValue)

}

func (m *Mapper) mapValue(src reflect.Value, dst reflect.Value) error {
	if !src.IsValid() || !dst.IsValid() {
		return nil
	}
	if src.Type().Kind() == reflect.Pointer || src.Type().Kind() == reflect.Interface {
		return m.mapValue(src.Elem(), dst)
	}

	if m.hasTypeConverter(src, dst) {
		converter := m.cfg.converters[structMapKey{
			source:      src.Type(),
			destination: dst.Type(),
		}]
		newVal, err := converter(src.Interface())
		if err != nil {
			return err
		}

		if newVal == nil {
			dst.Set(reflect.Zero(dst.Type()))
			return nil
		}
		if reflect.TypeOf(newVal) != dst.Type() {
			return ErrMismatchType
		}
		dst.Set(reflect.ValueOf(newVal))
		return nil
	}

	if src.Type().AssignableTo(dst.Type()) {
		if dst.CanSet() {
			dst.Set(src)
		}
		return nil
	}

	switch dst.Type().Kind() {
	case reflect.Uintptr:
		return nil // ignore
	case reflect.Array:
		if src.Type().Kind() != reflect.Array && src.Type().Kind() != reflect.Slice {
			return ErrMismatchType
		}

		if dst.Cap() < src.Len() {
			return ErrInsufficientCapacity
		}

		for i := 0; i < src.Len(); i++ {
			dstItem := dst.Index(i)
			err := m.mapValue(src.Index(i), dstItem)
			if err != nil {
				return err
			}
		}
	case reflect.Chan:
		return nil // ignore
	case reflect.Func:
		return nil // ignore
	case reflect.Interface:
		if dst.Elem().IsValid() {
			if !dst.Elem().CanAddr() { // for structs with interface fields, value in it is always not addressable
				return ErrNotAddresable
			}
			return m.mapValue(src, dst.Elem())
		}

		newVal := reflect.New(src.Type())
		err := m.mapValue(src, newVal)
		if err != nil {
			return err
		}
		dst.Set(newVal.Elem())
	case reflect.Map:
		if src.Type().Kind() != reflect.Map {
			return ErrMismatchType
		}

		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}

		iter := src.MapRange()
		for iter.Next() {
			// map key
			srcKey := iter.Key()
			dstKey := reflect.New(dst.Type().Key())
			err := m.mapValue(srcKey, dstKey)
			if err != nil {
				return err
			}

			// map value
			srcVal := iter.Value()
			dstVal := reflect.New(dst.Type().Elem())
			err = m.mapValue(srcVal, dstVal)
			if err != nil {
				return err
			}

			dst.SetMapIndex(dstKey.Elem(), dstVal.Elem())
		}
	case reflect.Pointer:
		if dst.IsNil() {
			new := reflect.New(dst.Type().Elem())
			dst.Set(new)
		}
		return m.mapValue(src, dst.Elem())
	case reflect.Slice:
		if src.Type().Kind() != reflect.Array && src.Type().Kind() != reflect.Slice {
			return ErrMismatchType
		}
		for i := 0; i < src.Len(); i++ {
			dstElem := reflect.New(dst.Type().Elem())

			err := m.mapValue(src.Index(i), dstElem.Elem())
			if err != nil {
				return err
			}
			dst.Set(reflect.Append(dst, dstElem.Elem()))
		}
	case reflect.Struct:
		if src.Type().Kind() != reflect.Struct {
			return ErrMismatchType
		}

		err := m.mapStructFields(src, dst)
		if err != nil {
			return err
		}
		err = m.mapStructSetters(src, dst)
		if err != nil {
			return err
		}
	case reflect.UnsafePointer:
		return nil // ignore
	default:
		return ErrMismatchType
	}
	return nil
}
func (m *Mapper) hasTypeConverter(src reflect.Value, dst reflect.Value) bool {
	structMapKey := structMapKey{
		source:      src.Type(),
		destination: dst.Type(),
	}
	_, ok := m.cfg.converters[structMapKey]
	return ok
}

func (m *Mapper) mapStructFields(src reflect.Value, dst reflect.Value) error {
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		srcFieldName := dst.Type().Field(i).Name
		srcField := src.FieldByName(srcFieldName)
		var err error
		if !srcField.IsValid() {
			// AI generated code block start
			getterName := "Get" + srcFieldName
			getterMethod := src.MethodByName(getterName)
			if getterMethod.IsValid() && getterMethod.Type().NumIn() == 0 && getterMethod.Type().NumOut() == 1 {
				srcField = getterMethod.Call(nil)[0]
			}
			// AI generated code block end
		}
		if !srcField.IsValid() {
			continue
		}
		err = m.mapValue(srcField, dstField)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mapper) mapStructSetters(src reflect.Value, dst reflect.Value) error {
	// AI generated code block start
	// Handle setter methods
	dstType := dst.Addr().Type()
	for i := 0; i < dstType.NumMethod(); i++ {
		method := dstType.Method(i)
		if method.Name[:3] == "Set" && method.Type.NumIn() == 2 && method.Type.NumOut() == 0 {
			fieldName := method.Name[3:]
			srcFieldName := fieldName
			srcField := src.FieldByName(srcFieldName)
			if !srcField.IsValid() {
				getterName := "Get" + srcFieldName
				getterMethod := src.MethodByName(getterName)
				if getterMethod.IsValid() && getterMethod.Type().NumIn() == 0 && getterMethod.Type().NumOut() == 1 {
					srcField = getterMethod.Call(nil)[0]
				} else {
					return ErrFieldNotFound
				}
			}
			if srcField.IsValid() {
				paramType := method.Type.In(1)
				paramValue := reflect.New(paramType).Elem()
				err := m.mapValue(srcField, paramValue)
				if err != nil {
					return err
				}
				setterMethod := dst.Addr().MethodByName(method.Name)
				setterMethod.Call([]reflect.Value{paramValue})
			}
		}
	}
	// AI generated code block end
	return nil
}
