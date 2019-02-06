//
// Really raw access to KStat data

package kstat

// #cgo LDFLAGS: -lkstat
//
// #include <sys/types.h>
// #include <stdlib.h>
// #include <strings.h>
// #include <kstat.h>
// #include <nfs/nfs_clnt.h>
//
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// Raw is the raw data of a KStat. The actual bytes are in Data;
// Ndata is kstat_t.ks_ndata, and is not normally useful.
//
// Note that with RawStat KStats, it turns out that Ndata == len(Data).
// This is contrary to its meaning for other types of kstats.
type Raw struct {
	Data     []byte
	Ndata    uint64
	Snaptime int64
	KStat    *KStat
}

// TODO: better functionality split here
func (k *KStat) prep() error {
	if k.invalid() {
		return errors.New("invalid KStat or closed token")
	}

	// Do the initial load of the data if necessary.
	if k.ksp.ks_data == nil {
		if err := k.Refresh(); err != nil {
			return err
		}
	}
	return nil
}

// Raw returns the raw byte data of a KStat. It may be called on any
// KStat. It does not refresh the KStat's data.
func (k *KStat) Raw() (*Raw, error) {
	if err := k.prep(); err != nil {
		return nil, err
	}
	r := Raw{}
	r.KStat = k
	r.Snaptime = k.Snaptime
	r.Ndata = uint64(k.ksp.ks_ndata)
	// The forced C.int() conversion is dangerous, because C.int
	// is not necessarily large enough to contain a
	// size_t. However this is the interface that Go gives us, so
	// we live with it.
	r.Data = C.GoBytes(unsafe.Pointer(k.ksp.ks_data), C.int(k.ksp.ks_data_size))
	return &r, nil
}

func (tok *Token) prepunix(name string, size uintptr) (*KStat, error) {
	k, err := tok.Lookup("unix", 0, name)
	if err != nil {
		return nil, err
	}
	// TODO: handle better?
	if k.ksp.ks_type != C.KSTAT_TYPE_RAW {
		return nil, fmt.Errorf("%s is wrong type %s", k, k.Type)
	}
	if uintptr(k.ksp.ks_data_size) != size {
		return nil, fmt.Errorf("%s is wrong size %d (should be %d)", k, k.ksp.ks_data_size, size)
	}
	return k, nil
}

// Sysinfo returns the KStat and the statistics from unix:0:sysinfo.
// It always returns a current, refreshed copy.
func (tok *Token) Sysinfo() (*KStat, *Sysinfo, error) {
	var si Sysinfo
	k, err := tok.prepunix("sysinfo", unsafe.Sizeof(si))
	if err != nil {
		return nil, nil, err
	}
	si = *((*Sysinfo)(k.ksp.ks_data))
	return k, &si, nil
}

// Vminfo returns the KStat and the statistics from unix:0:vminfo.
// It always returns a current, refreshed copy.
func (tok *Token) Vminfo() (*KStat, *Vminfo, error) {
	var vi Vminfo
	k, err := tok.prepunix("vminfo", unsafe.Sizeof(vi))
	if err != nil {
		return nil, nil, err
	}
	vi = *((*Vminfo)(k.ksp.ks_data))
	return k, &vi, nil
}

// Var returns the KStat and the statistics from unix:0:var.
// It always returns a current, refreshed copy.
func (tok *Token) Var() (*KStat, *Var, error) {
	var vi Var
	k, err := tok.prepunix("var", unsafe.Sizeof(vi))
	if err != nil {
		return nil, nil, err
	}
	vi = *((*Var)(k.ksp.ks_data))
	return k, &vi, nil
}

// GetMntinfo retrieves a Mntinfo struct from a nfs:*:mntinfo KStat.
// It does not force a refresh of the KStat.
func (k *KStat) GetMntinfo() (*Mntinfo, error) {
	var mi Mntinfo
	if err := k.prep(); err != nil {
		return nil, err
	}
	if k.Type != RawStat || k.Module != "nfs" || k.Name != "mntinfo" {
		return nil, errors.New("KStat is not a Mntinfo kstat")
	}
	if uintptr(k.ksp.ks_data_size) != unsafe.Sizeof(mi) {
		return nil, fmt.Errorf("KStat is wrong size %d (should be %d)", k.ksp.ks_data_size, unsafe.Sizeof(mi))
	}
	mi = *((*Mntinfo)(k.ksp.ks_data))
	return &mi, nil
}

//
// Support for copying semi-arbitrary structures out of raw
// KStats.
//

// safeThing returns true if a given type is either a simple defined
// size primitive integer type or an array and/or struct composed
// entirely of safe things. A safe thing is entirely self contained
// and may be initialized from random memory without breaking Go's
// memory safety (although the values it contains may be garbage).
//
func safeThing(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Array:
		// an array is safe if it's an array of something safe
		return safeThing(t.Elem())
	case reflect.Struct:
		// a struct is safe if all its components are safe
		for i := 0; i < t.NumField(); i++ {
			if !safeThing(t.Field(i).Type) {
				return false
			}
		}
		return true
	default:
		// other things are not safe.
		return false
	}
}

// TODO: add floats to the supported list? It's unlikely to be needed
// but it should just work.

// CopyTo copies a RawStat KStat into a struct that you supply a
// pointer to. The size of the struct must exactly match the size of
// the RawStat's data.
//
// CopyStat imposes conditions on the struct that you are copying to:
// it must be composed entirely of primitive integer types with defined
// sizes (intN and uintN), or arrays and structs that ultimately only
// contain them. All fields should be exported.
//
// If you give CopyStat a bad argument, it generally panics.
//
// This API is provisional and may be changed or deleted.
func (k *KStat) CopyTo(ptr interface{}) error {
	if err := k.prep(); err != nil {
		return err
	}

	if k.Type != RawStat {
		return errors.New("KStat is not a RawStat")
	}

	// Validity checks: not nil value, not nil pointer value,
	// is a pointer to struct.
	if ptr == nil {
		panic("CopyTo given nil pointer")
	}
	vp := reflect.ValueOf(ptr)
	if vp.Kind() != reflect.Ptr {
		panic("CopyTo not given a pointer")
	}
	if vp.IsNil() {
		panic("CopyTo given nil pointer")
	}
	dst := vp.Elem()
	if dst.Kind() != reflect.Struct {
		panic("CopyTo: not pointer to struct")
	}
	// Is the struct safe to copy into, which means primitive types
	// and structs/arrays of primitive types?
	if !safeThing(dst.Type()) {
		panic("CopyTo: not a safe structure, contains unsupported fields")
	}
	if !dst.CanSet() {
		panic("CopyTo: struct cannot be set for some reason")
	}

	// Verify that the size of the target struct matches the size
	// of the raw KStat.
	if uintptr(k.ksp.ks_data_size) != dst.Type().Size() {
		return errors.New("struct size does not match KStat size")
	}

	// The following is exactly the magic that we performed for
	// specific types earlier. We take k.ksp.ks_data and turn
	// it into a typed pointer to the target object's type:
	//
	//	src := ((*<type>)(k.kps.ks_data))
	src := reflect.NewAt(dst.Type(), unsafe.Pointer(k.ksp.ks_data))

	// We now dereference that into the destination to copy the
	// data:
	//
	//	dst = *src
	dst.Set(reflect.Indirect(src))

	return nil
}
