//
// The kstat package provides a Go interface to the Solaris/OmniOS
// kstat(s) system for user-level access to a lot of kernel
// statistics. For more documentation on kstats, see kstat(1) and
// kstat(3kstat).
//
// In an ideal world the package documentation would go here. This is
// not an ideal world, because any number of tools like godoc choke on
// Go files that are not for their architecture (although I'll admit
// it's a hard problem). So see doc.go for the actual package level
// documentation.
//
// However, I refuse to push function level API documentation off to another
// file, at least at the moment. It would be a horrible mess.
//

package kstat

// #cgo LDFLAGS: -lkstat
//
// #include <sys/types.h>
// #include <stdlib.h>
// #include <strings.h>
// #include <kstat.h>
//
// /* We have to reach through unions, which cgo doesn't support.
//    So we have our own cheesy little routines for it. These assume
//    they are always being called on validly-typed named kstats.
//  */
//
// char *get_named_char(kstat_named_t *knp) {
//	return knp->value.str.addr.ptr;
// }
//
// uint64_t get_named_uint(kstat_named_t *knp) {
//	if (knp->data_type == KSTAT_DATA_UINT32)
//		return knp->value.ui32;
//	else
//		return knp->value.ui64;
// }
//
// int64_t get_named_int(kstat_named_t *knp) {
//	if (knp->data_type == KSTAT_DATA_INT32)
//		return knp->value.i32;
//	else
//		return knp->value.i64;
// }
//
// /* Let's not try to do C pointer arithmetic in Go and get it wrong */
// kstat_named_t *get_nth_named(kstat_t *ks, uint_t n) {
//	kstat_named_t *knp;
//	if (!ks || !ks->ks_data || ks->ks_type != KSTAT_TYPE_NAMED || n >= ks->ks_ndata)
//		return NULL;
//	knp = KSTAT_NAMED_PTR(ks);
//	return knp + n;
// }
//
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Token is an access token for obtaining kstats.
type Token struct {
	kc *C.struct_kstat_ctl

	// ksm maps kstat_t pointers to our Go-level KStats for them.
	// kstat_t's stay constant over the lifetime of a token, so
	// we want to keep unique KStats. This holds some Go-level
	// memory down, but I wave my hands.
	ksm map[*C.struct_kstat]*KStat
}

// Open returns a kstat Token that is used to obtain kstats. It corresponds
// to kstat_open(). You should call .Close() when you're done and then not
// use any KStats or Nameds obtained through this token.
//
// (Failing to call .Close() will cause memory leaks.)
func Open() (*Token, error) {
	r, err := C.kstat_open()
	if r == nil {
		return nil, err
	}
	t := Token{}
	t.kc = r
	t.ksm = make(map[*C.struct_kstat]*KStat)
	// A 'func (t *Token) Close()' is equivalent to
	// 'func Close(t *Token)'. The latter is what SetFinalizer()
	// needs.
	runtime.SetFinalizer(&t, (*Token).Close)
	return &t, nil
}

// Close a kstat access token. A closed token cannot be used for
// anything and cannot be reopened.
//
// After a Token has been closed it remains safe to look at fields
// on KStat and Named objects obtained through the Token, but it is
// not safe to call methods on them other than String(); doing so
// may cause memory corruption, although we try to avoid that.
//
// This corresponds to kstat_close().
func (t *Token) Close() error {
	if t == nil || t.kc == nil {
		return nil
	}

	// Go through our KStats and null out fields that are no longer
	// valid. We opt to do this before we actually destroy the memory
	// KStat.ksp is pointing to by calling kstat_close().
	for _, v := range t.ksm {
		v.ksp = nil
		v.tok = nil
	}

	res, err := C.kstat_close(t.kc)
	t.kc = nil

	// clear the map to drop all references to KStats.
	t.ksm = make(map[*C.struct_kstat]*KStat)

	// cancel finalizer
	runtime.SetFinalizer(&t, nil)

	if res != 0 {
		return err
	}
	return nil
}

// Update synchronizes the Token to the current state of available
// kernel kstats, returning true if the kernel's list of available
// kstats changed and false otherwise. If there have been no changes
// in the kernel's kstat list, all KStats remain valid. If there was a
// kstat update, some or all of the KStats obtained through the Token
// may now be invalid. Some of the now-invalid KStats may still exist
// and be the same thing, but if so they will have to be looked up
// again.
//
// (This happens if, for example, a device disappears and then
// reappears. At the kernel level, the device's kstat is deleted when
// it disappears and then is recreated when it reappears; the kernel
// considers the recreated version to be a different kstat, although
// it has the same module:instance:name. Note that the same
// module:instance:name still existing does not guarantee that the
// kstat is for the same thing; one disk might have removed and then
// an entirely different new disk added.)
//
// Update corresponds to kstat_chain_update().
func (t *Token) Update() (bool, error) {
	if t == nil || t.kc == nil {
		return true, errors.New("token is closed")
	}
	oid := t.kc.kc_chain_id
	// NOTE that we can't assume err == nil on success and just
	// check for err != nil. The error return is set from errno,
	// and kstat_chain_update() does not guarantee that errno is
	// 0 if it succeeds.
	nid, err := C.kstat_chain_update(t.kc)
	switch {
	case nid < 0:
		// We generously assume that if there has been an
		// error, the chain is intact. Otherwise we should
		// invalidate all KStats in t.ksm, as in .Close().
		// assumption: err != nil if n < 0.
		return false, err
	case nid == 0:
		// No change is good news.
		return false, nil
	case nid == oid:
		// Should never be the case, but...
		return false, fmt.Errorf("new KCID is old KCID: %d", nid)
	}

	// The simple approach to KStats after a chain update would be
	// to invalidate all existing KStats. However, we can do
	// better. kstat_chain_update() implicitly guarantees that it
	// will not reuse memory addresses of kstat_t structures for
	// different ones within a single call, so we can walk the
	// chain and look for addresses that we already know; the
	// KStats for those addresses are still valid.

	// Copy all valid chain entries that we have in the token ksm
	// map to a new map and delete them from the old (current) map.
	nksm := make(map[*C.struct_kstat]*KStat)
	for r := t.kc.kc_chain; r != nil; r = r.ks_next {
		if v, ok := t.ksm[r]; ok {
			nksm[r] = v
			delete(t.ksm, r)
		}
	}
	// Anything left in t.ksm is an old chain entry that was
	// removed by kstat_chain_update(). Explicitly zap their
	// KStat's references to make them invalid.
	for _, v := range t.ksm {
		v.ksp = nil
		v.tok = nil
	}
	// Make our new ksm map the current ksm map.
	t.ksm = nksm

	return true, nil
}

// All returns an array of all available KStats.
//
// (It has no error return because due to how kstats are implemented,
// it cannot fail.)
func (t *Token) All() []*KStat {
	n := []*KStat{}
	if t == nil || t.kc == nil {
		return n
	}

	for r := t.kc.kc_chain; r != nil; r = r.ks_next {
		n = append(n, newKStat(t, r))
	}
	return n
}

//
// allocate a C string for a non-blank string; otherwise return nil
func maybeCString(src string) *C.char {
	if src == "" {
		return nil
	}
	return C.CString(src)
}

// free a non-nil C string
func maybeFree(cs *C.char) {
	if cs != nil {
		C.free(unsafe.Pointer(cs))
	}
}

// strndup behaves like the C function; given a *C.char and a len, it
// returns a string that is up to len characters long at most.
// Shorn of casts, it is:
//	C.GoStringN(p, C.strnlen(p, len))
//
// strndup() is necessary to copy fields of the type 'char
// name[SIZE];' where a string of exactly SIZE length will not be
// null-terminated. GoStringN() will always copy trailing null bytes
// and other garbage; GoString()'s internal strlen() may run off the
// end of the 'name' field and either fault or copy too much.
func strndup(cs *C.char, len C.size_t) string {
	// credit: Ian Lance Taylor in
	// https://github.com/golang/go/issues/12428
	return C.GoStringN(cs, C.int(C.strnlen(cs, len)))
}

// Lookup looks up a particular kstat. module and name may be "" and
// instance may be -1 to mean 'the first one that kstats can find'.
// It also refreshes (or retrieves) the kstat's data and thus sets
// Snaptime.
//
// Lookup() corresponds to kstat_lookup() *plus kstat_read()*.
func (t *Token) Lookup(module string, instance int, name string) (*KStat, error) {
	if t == nil || t.kc == nil {
		return nil, errors.New("Token not valid or closed")
	}

	ms := maybeCString(module)
	ns := maybeCString(name)
	r, err := C.kstat_lookup(t.kc, ms, C.int(instance), ns)
	maybeFree(ms)
	maybeFree(ns)

	if r == nil {
		return nil, err
	}

	k := newKStat(t, r)

	// People rarely look up kstats to not use them, so we immediately
	// attempt to kstat_read() the data. If this fails, we don't return
	// the kstat. However, we don't scrub it from the kstat_t mapping
	// that the Token maintains; we have no reason to believe that it
	// needs to be remade. Our return of nil is a convenience to avoid
	// problems in callers.
	// TODO: this may be a mistake in the API.
	//
	// NOTE: this means that calling Lookup() on an existing KStat
	// (either directly or via tok.GetNamed()) has the effect of
	// updating its statistics data to the current time. Right now
	// we consider this a feature.
	err = k.Refresh()
	if err != nil {
		return nil, err
	}
	return k, nil
}

// GetNamed obtains the Named representing a particular (named) kstat
// module:instance:name:statistic statistic. It always returns current
// data for the kstat statistic, even if it's called repeatedly for the
// same statistic.
//
// It is equivalent to .Lookup() then KStat.GetNamed().
func (t *Token) GetNamed(module string, instance int, name, stat string) (*Named, error) {
	stats, err := t.Lookup(module, instance, name)
	if err != nil {
		return nil, err
	}
	return stats.GetNamed(stat)
}

// -----

// KSType is the type of the data in a KStat.
type KSType int

// The different types of data that a KStat may contain, ie these
// are the value of a KStat.Type. We currently only support getting
// Named and IO statistics.
const (
	RawStat   KSType = C.KSTAT_TYPE_RAW
	NamedStat KSType = C.KSTAT_TYPE_NAMED
	IntrStat  KSType = C.KSTAT_TYPE_INTR
	IoStat    KSType = C.KSTAT_TYPE_IO
	TimerStat KSType = C.KSTAT_TYPE_TIMER
)

func (tp KSType) String() string {
	switch tp {
	case RawStat:
		return "raw"
	case NamedStat:
		return "named"
	case IntrStat:
		return "interrupt"
	case IoStat:
		return "io"
	case TimerStat:
		return "timer"
	default:
		return fmt.Sprintf("kstat_type:%d", tp)
	}
}

// KStat is the access handle for the collection of statistics for a
// particular module:instance:name kstat.
//
type KStat struct {
	Module   string
	Instance int
	Name     string

	// Class is eg 'net' or 'disk'. In kstat(1) it shows up as a
	// ':class' statistic.
	Class string
	// Type is the type of kstat.
	Type KSType

	// Creation time of a kstat in nanoseconds since sometime.
	// See gethrtime(3) and kstat(3kstat).
	Crtime int64
	// Snaptime is what kstat(1) reports as 'snaptime', the time
	// that this data was obtained. As with Crtime, it is in
	// nanoseconds since some arbitrary point in time.
	// Snaptime may not be valid until .Refresh() or .GetNamed()
	// has been called.
	Snaptime int64

	ksp *C.struct_kstat
	// We need access to the token to refresh the data
	tok *Token
}

// newKStat is our internal KStat constructor.
//
// This also has the responsibility of maintaining (and using) the
// kstat_t to KStat mapping cache, so that we don't recreate new
// KStats for the same kstat_t all the time.
func newKStat(tok *Token, ks *C.struct_kstat) *KStat {
	if kst, ok := tok.ksm[ks]; ok {
		return kst
	}

	kst := KStat{}
	kst.ksp = ks
	kst.tok = tok

	kst.Instance = int(ks.ks_instance)
	kst.Module = strndup((*C.char)(unsafe.Pointer(&ks.ks_module)), C.KSTAT_STRLEN)
	kst.Name = strndup((*C.char)(unsafe.Pointer(&ks.ks_name)), C.KSTAT_STRLEN)
	kst.Class = strndup((*C.char)(unsafe.Pointer(&ks.ks_class)), C.KSTAT_STRLEN)
	kst.Type = KSType(ks.ks_type)
	kst.Crtime = int64(ks.ks_crtime)

	// Inside the kernel, the ks_snaptime of a kstat is of course
	// a global thing. This 'global' snaptime is copied to user
	// level as part of the kstat header(s) on kstat_open(), which
	// means that kstats that have never been kstat_read() by us
	// are almost certain to have a non-zero ks_snaptime (because
	// someone, somewhere, will have read them since the system
	// booted, eg 'kstat -p | grep ...'  reads all kstats).
	// Because this ks_snaptime is not useful, we don't copy it
	// to Snaptime; instead we leave Snaptime unset (zero) as
	// an explicit signal that this KStat has never had its data
	// read.
	//
	//kst.Snaptime = int64(ks.ks_snaptime)

	tok.ksm[ks] = &kst
	return &kst
}

// invalid is a desperate attempt to keep usage errors from causing
// memory corruption. Don't count on it.
func (k *KStat) invalid() bool {
	return k == nil || k.ksp == nil || k.tok == nil || k.tok.kc == nil
}

// setup does validity checks and setup, such as loading data via Refresh().
// It applies only to named kstats.
//
// TODO: setup() vs prep() is a code smell.
func (k *KStat) setup() error {
	if k.invalid() {
		return errors.New("invalid KStat or closed token")
	}

	if k.ksp.ks_type != C.KSTAT_TYPE_NAMED {
		return fmt.Errorf("kstat %s (type %d) is not a named kstat", k, k.ksp.ks_type)
	}

	// Do the initial load of the data if necessary.
	if k.ksp.ks_data == nil {
		if err := k.Refresh(); err != nil {
			return err
		}
	}
	return nil
}

func (k *KStat) String() string {
	return fmt.Sprintf("%s:%d:%s (%s)", k.Module, k.Instance, k.Name, k.Class)
}

// Valid returns true if a KStat is still valid after a Token.Update()
// call has returned true. If a KStat becomes invalid after an update,
// its fields remain available but you can no longer call methods on
// it. You may be able to look it up again with token.Lookup(k.Module,
// k.Instance, k.Name), although it's possible that the
// module:instance:name now refers to something else. Even if it is
// still the same thing, there is no continuity in the actual
// statistics once Valid becomes false; you must restart tracking from
// scratch.
//
// (For example, if one disk is removed from the system and another is
// added, the new disk may use the same module:instance:name as some
// of the old disk's KStats. Your .Lookup() may succeed, but what you
// get back is not in any way a continuation of the old disk's
// information.)
//
// Valid also returns false after the KStat's token has been closed.
func (k *KStat) Valid() bool {
	return !k.invalid()
}

// Refresh the statistics data for a KStat.
//
// Note that this does not update any existing Named objects for
// statistics from this KStat. You must re-do .GetNamed() to get
// new ones in order to see any updates.
//
// Under the hood this does a kstat_read(). You don't need to call it
// explicitly before obtaining statistics from a KStat.
func (k *KStat) Refresh() error {
	if k.invalid() {
		return errors.New("invalid KStat or closed token")
	}

	res, err := C.kstat_read(k.tok.kc, k.ksp, nil)
	if res == -1 {
		return err
	}
	k.Snaptime = int64(k.ksp.ks_snaptime)
	return nil
}

// GetIO retrieves the IO statistics data from an IoStat type
// KStat. It always refreshes the KStat to provide current data.
//
// It corresponds to kstat_read() followed by getting a copy of
// ks_data (which is a kstat_io_t).
func (k *KStat) GetIO() (*IO, error) {
	if err := k.Refresh(); err != nil {
		return nil, err
	}
	if k.ksp.ks_type != C.KSTAT_TYPE_IO {
		return nil, fmt.Errorf("kstat %s (type %d) is not an IO kstat", k, k.ksp.ks_type)
	}

	// We make our own copy of ks_data (as an IO) so that we don't
	// point into C-owned memory. 'go tool cgo -godef' apparently
	// guarantees that the IO struct/type it creates has exactly
	// the same in-memory layout as the C struct, so we can safely
	// do this copy and expect to get good results.
	io := IO{}
	io = *((*IO)(k.ksp.ks_data))
	return &io, nil
}

// GetNamed obtains a particular named statistic from a KStat. It does
// not refresh the KStat's statistics data, so multiple calls to
// GetNamed on a single KStat will get a coherent set of statistic
// values from it.
//
// It corresponds to kstat_data_lookup().
func (k *KStat) GetNamed(name string) (*Named, error) {
	if err := k.setup(); err != nil {
		return nil, err
	}
	ns := C.CString(name)
	r, err := C.kstat_data_lookup(k.ksp, ns)
	C.free(unsafe.Pointer(ns))
	if r == nil || err != nil {
		return nil, err
	}
	return newNamed(k, (*C.struct_kstat_named)(r)), err
}

// AllNamed returns an array of all named statistics for a particular
// named-type KStat. Entries are returned in no particular order.
func (k *KStat) AllNamed() ([]*Named, error) {
	if err := k.setup(); err != nil {
		return nil, err
	}
	lst := make([]*Named, k.ksp.ks_ndata)
	for i := C.uint_t(0); i < k.ksp.ks_ndata; i++ {
		ks := C.get_nth_named(k.ksp, i)
		if ks == nil {
			panic("get_nth_named returned surprise nil")
		}
		lst[i] = newNamed(k, ks)
	}
	return lst, nil
}

// Named represents a particular kstat named statistic, ie the full
//	module:instance:name:statistic
// and its current value.
//
// Name and Type are always valid, but only one of StringVal, IntVal,
// or UintVal is valid for any particular statistic; which one is
// valid is determined by its Type. Generally you'll already know what
// type a given named kstat statistic is; I don't believe Solaris
// changes their type once they're defined.
type Named struct {
	Name string
	Type NamedType

	// Only one of the following values is valid; the others are zero
	// values.
	//
	// StringVal holds the value for both CharData and String Type(s).
	StringVal string
	IntVal    int64
	UintVal   uint64

	// The Snaptime this Named was obtained. Note that while you
	// use the parent KStat's Crtime, you cannot use its Snaptime.
	// The KStat may have been refreshed since this Named was
	// created, which updates the Snaptime.
	Snaptime int64

	// Pointer to the parent KStat, for access to the full name
	// and the crtime associated with this Named.
	KStat *KStat
}

func (ks *Named) String() string {
	return fmt.Sprintf("%s:%d:%s:%s", ks.KStat.Module, ks.KStat.Instance, ks.KStat.Name, ks.Name)
}

// NamedType represents the various types of named kstat statistics.
type NamedType int

// The different types of data that a named kstat statistic can be
// (ie, these are the potential values of Named.Type).
const (
	CharData NamedType = C.KSTAT_DATA_CHAR
	Int32    NamedType = C.KSTAT_DATA_INT32
	Uint32   NamedType = C.KSTAT_DATA_UINT32
	Int64    NamedType = C.KSTAT_DATA_INT64
	Uint64   NamedType = C.KSTAT_DATA_UINT64
	String   NamedType = C.KSTAT_DATA_STRING

	// CharData is found in StringVal. At the moment we assume that
	// it is a real string, because this matches how it seems to be
	// used for short strings in the Solaris kernel. Someday we may
	// find something that uses it as just a data dump for 16 bytes.

	// Solaris sys/kstat.h also has _FLOAT (5) and _DOUBLE (6) types,
	// but labels them as obsolete.
)

func (tp NamedType) String() string {
	switch tp {
	case CharData:
		return "char"
	case Int32:
		return "int32"
	case Uint32:
		return "uint32"
	case Int64:
		return "int64"
	case Uint64:
		return "uint64"
	case String:
		return "string"
	default:
		return fmt.Sprintf("named_type-%d", tp)
	}
}

// Create a new Stat from the kstat_named_t
// We set the appropriate *Value field.
func newNamed(k *KStat, knp *C.struct_kstat_named) *Named {
	st := Named{}
	st.KStat = k
	st.Name = strndup((*C.char)(unsafe.Pointer(&knp.name)), C.KSTAT_STRLEN)
	st.Type = NamedType(knp.data_type)
	st.Snaptime = k.Snaptime

	switch st.Type {
	case String:
		// The comments in sys/kstat.h explicitly guarantee
		// that these strings are null-terminated, although
		// knp.value.str.len also holds the length.
		st.StringVal = C.GoString(C.get_named_char(knp))
	case CharData:
		// Solaris/etc appears to use CharData for short strings
		// so that they can be embedded directly into
		// knp.value.c[16] instead of requiring an out of line
		// allocation. In theory we may find someone who is
		// using it as 128-bit ints or the like.
		// However I scanned the Illumos kernel source and
		// everyone using it appears to really be using it for
		// strings.
		st.StringVal = strndup((*C.char)(unsafe.Pointer(&knp.value)), 16)
	case Int32, Int64:
		st.IntVal = int64(C.get_named_int(knp))
	case Uint32, Uint64:
		st.UintVal = uint64(C.get_named_uint(knp))
	default:
		// TODO: should do better.
		panic(fmt.Sprintf("unknown stat type: %d", st.Type))
	}
	return &st
}
