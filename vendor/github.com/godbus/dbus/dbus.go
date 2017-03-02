package dbus

import (
	"errors"
	"reflect"
	"strings"
)

var (
	byteType        = reflect.TypeOf(byte(0))
	boolType        = reflect.TypeOf(false)
	uint8Type       = reflect.TypeOf(uint8(0))
	int16Type       = reflect.TypeOf(int16(0))
	uint16Type      = reflect.TypeOf(uint16(0))
	intType         = reflect.TypeOf(int(0))
	uintType        = reflect.TypeOf(uint(0))
	int32Type       = reflect.TypeOf(int32(0))
	uint32Type      = reflect.TypeOf(uint32(0))
	int64Type       = reflect.TypeOf(int64(0))
	uint64Type      = reflect.TypeOf(uint64(0))
	float64Type     = reflect.TypeOf(float64(0))
	stringType      = reflect.TypeOf("")
	signatureType   = reflect.TypeOf(Signature{""})
	objectPathType  = reflect.TypeOf(ObjectPath(""))
	variantType     = reflect.TypeOf(Variant{Signature{""}, nil})
	interfacesType  = reflect.TypeOf([]interface{}{})
	unixFDType      = reflect.TypeOf(UnixFD(0))
	unixFDIndexType = reflect.TypeOf(UnixFDIndex(0))
)

// An InvalidTypeError signals that a value which cannot be represented in the
// D-Bus wire format was passed to a function.
type InvalidTypeError struct {
	Type reflect.Type
}

func (e InvalidTypeError) Error() string {
	return "dbus: invalid type " + e.Type.String()
}

// Store copies the values contained in src to dest, which must be a slice of
// pointers. It converts slices of interfaces from src to corresponding structs
// in dest. An error is returned if the lengths of src and dest or the types of
// their elements don't match.
func Store(src []interface{}, dest ...interface{}) error {
	if len(src) != len(dest) {
		return errors.New("dbus.Store: length mismatch")
	}

	for i := range src {
		if err := storeInterfaces(src[i], dest[i]); err != nil {
			return err
		}
	}
	return nil
}

func storeInterfaces(src, dest interface{}) error {
	return store(reflect.ValueOf(src), reflect.ValueOf(dest))
}

func store(src, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Ptr:
		return store(src, dest.Elem())
	case reflect.Interface:
		return storeInterface(src, dest)
	case reflect.Slice:
		return storeSlice(src, dest)
	case reflect.Map:
		return storeMap(src, dest)
	case reflect.Struct:
		return storeStruct(src, dest)
	default:
		return storeBase(src, dest)
	}
}

func storeBase(src, dest reflect.Value) error {
	return setDest(dest, src)
}

func setDest(dest, src reflect.Value) error {
	if !isVariant(src.Type()) && isVariant(dest.Type()) {
		//special conversion for dbus.Variant
		dest.Set(reflect.ValueOf(MakeVariant(src.Interface())))
		return nil
	}
	if !src.Type().ConvertibleTo(dest.Type()) {
		return errors.New(
			"dbus.Store: type mismatch")
	}
	dest.Set(src.Convert(dest.Type()))
	return nil
}

func storeStruct(sv, rv reflect.Value) error {
	if !sv.Type().AssignableTo(interfacesType) {
		return setDest(rv, sv)
	}
	vs := sv.Interface().([]interface{})
	t := rv.Type()
	ndest := make([]interface{}, 0, rv.NumField())
	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath == "" && field.Tag.Get("dbus") != "-" {
			ndest = append(ndest,
				rv.Field(i).Addr().Interface())

		}
	}
	if len(vs) != len(ndest) {
		return errors.New("dbus.Store: type mismatch")
	}
	err := Store(vs, ndest...)
	if err != nil {
		return errors.New("dbus.Store: type mismatch")
	}
	return nil
}

func storeMap(sv, rv reflect.Value) error {
	if sv.Kind() != reflect.Map {
		return errors.New("dbus.Store: type mismatch")
	}
	keys := sv.MapKeys()
	rv.Set(reflect.MakeMap(rv.Type()))
	destElemType := rv.Type().Elem()
	for _, key := range keys {
		elemv := sv.MapIndex(key)
		v := newDestValue(elemv, destElemType)
		err := store(getVariantValue(elemv), v)
		if err != nil {
			return err
		}
		if !v.Elem().Type().ConvertibleTo(destElemType) {
			return errors.New(
				"dbus.Store: type mismatch")
		}
		rv.SetMapIndex(key, v.Elem().Convert(destElemType))
	}
	return nil
}

func storeSlice(sv, rv reflect.Value) error {
	if sv.Kind() != reflect.Slice {
		return errors.New("dbus.Store: type mismatch")
	}
	rv.Set(reflect.MakeSlice(rv.Type(), sv.Len(), sv.Len()))
	destElemType := rv.Type().Elem()
	for i := 0; i < sv.Len(); i++ {
		v := newDestValue(sv.Index(i), destElemType)
		err := store(getVariantValue(sv.Index(i)), v)
		if err != nil {
			return err
		}
		err = setDest(rv.Index(i), v.Elem())
		if err != nil {
			return err
		}
	}
	return nil
}

func storeInterface(sv, rv reflect.Value) error {
	return setDest(rv, getVariantValue(sv))
}

func getVariantValue(in reflect.Value) reflect.Value {
	if isVariant(in.Type()) {
		return reflect.ValueOf(in.Interface().(Variant).Value())
	}
	return in
}

func newDestValue(srcValue reflect.Value, destType reflect.Type) reflect.Value {
	switch srcValue.Kind() {
	case reflect.Map:
		switch {
		case !isVariant(srcValue.Type().Elem()):
			return reflect.New(destType)
		case destType.Kind() == reflect.Map:
			return reflect.New(destType)
		default:
			return reflect.New(
				reflect.MapOf(srcValue.Type().Key(), destType))
		}

	case reflect.Slice:
		switch {
		case !isVariant(srcValue.Type().Elem()):
			return reflect.New(destType)
		case destType.Kind() == reflect.Slice:
			return reflect.New(destType)
		default:
			return reflect.New(
				reflect.SliceOf(destType))
		}
	default:
		if !isVariant(srcValue.Type()) {
			return reflect.New(destType)
		}
		return newDestValue(getVariantValue(srcValue), destType)
	}
}

func isVariant(t reflect.Type) bool {
	return t == variantType
}

// An ObjectPath is an object path as defined by the D-Bus spec.
type ObjectPath string

// IsValid returns whether the object path is valid.
func (o ObjectPath) IsValid() bool {
	s := string(o)
	if len(s) == 0 {
		return false
	}
	if s[0] != '/' {
		return false
	}
	if s[len(s)-1] == '/' && len(s) != 1 {
		return false
	}
	// probably not used, but technically possible
	if s == "/" {
		return true
	}
	split := strings.Split(s[1:], "/")
	for _, v := range split {
		if len(v) == 0 {
			return false
		}
		for _, c := range v {
			if !isMemberChar(c) {
				return false
			}
		}
	}
	return true
}

// A UnixFD is a Unix file descriptor sent over the wire. See the package-level
// documentation for more information about Unix file descriptor passsing.
type UnixFD int32

// A UnixFDIndex is the representation of a Unix file descriptor in a message.
type UnixFDIndex uint32

// alignment returns the alignment of values of type t.
func alignment(t reflect.Type) int {
	switch t {
	case variantType:
		return 1
	case objectPathType:
		return 4
	case signatureType:
		return 1
	case interfacesType:
		return 4
	}
	switch t.Kind() {
	case reflect.Uint8:
		return 1
	case reflect.Uint16, reflect.Int16:
		return 2
	case reflect.Uint, reflect.Int, reflect.Uint32, reflect.Int32, reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return 4
	case reflect.Uint64, reflect.Int64, reflect.Float64, reflect.Struct:
		return 8
	case reflect.Ptr:
		return alignment(t.Elem())
	}
	return 1
}

// isKeyType returns whether t is a valid type for a D-Bus dict.
func isKeyType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float64,
		reflect.String, reflect.Uint, reflect.Int:

		return true
	}
	return false
}

// isValidInterface returns whether s is a valid name for an interface.
func isValidInterface(s string) bool {
	if len(s) == 0 || len(s) > 255 || s[0] == '.' {
		return false
	}
	elem := strings.Split(s, ".")
	if len(elem) < 2 {
		return false
	}
	for _, v := range elem {
		if len(v) == 0 {
			return false
		}
		if v[0] >= '0' && v[0] <= '9' {
			return false
		}
		for _, c := range v {
			if !isMemberChar(c) {
				return false
			}
		}
	}
	return true
}

// isValidMember returns whether s is a valid name for a member.
func isValidMember(s string) bool {
	if len(s) == 0 || len(s) > 255 {
		return false
	}
	i := strings.Index(s, ".")
	if i != -1 {
		return false
	}
	if s[0] >= '0' && s[0] <= '9' {
		return false
	}
	for _, c := range s {
		if !isMemberChar(c) {
			return false
		}
	}
	return true
}

func isMemberChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') || c == '_'
}
