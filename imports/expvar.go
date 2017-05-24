// this file was generated by gomacro command: import _b "expvar"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"expvar"
)

// reflection: allow interpreted code to import "expvar"
func init() {
	Packages["expvar"] = Package{
	Binds: map[string]Value{
		"Do":	ValueOf(expvar.Do),
		"Get":	ValueOf(expvar.Get),
		"Handler":	ValueOf(expvar.Handler),
		"NewFloat":	ValueOf(expvar.NewFloat),
		"NewInt":	ValueOf(expvar.NewInt),
		"NewMap":	ValueOf(expvar.NewMap),
		"NewString":	ValueOf(expvar.NewString),
		"Publish":	ValueOf(expvar.Publish),
	},
	Types: map[string]Type{
		"Float":	TypeOf((*expvar.Float)(nil)).Elem(),
		"Func":	TypeOf((*expvar.Func)(nil)).Elem(),
		"Int":	TypeOf((*expvar.Int)(nil)).Elem(),
		"KeyValue":	TypeOf((*expvar.KeyValue)(nil)).Elem(),
		"Map":	TypeOf((*expvar.Map)(nil)).Elem(),
		"String":	TypeOf((*expvar.String)(nil)).Elem(),
		"Var":	TypeOf((*expvar.Var)(nil)).Elem(),
	},
	Proxies: map[string]Type{
		"Var":	TypeOf((*Var_expvar)(nil)).Elem(),
	} }
}

// --------------- proxy for expvar.Var ---------------
type Var_expvar struct {
	Object	interface{}
	String_	func() string
}
func (Proxy *Var_expvar) String() string {
	return Proxy.String_()
}
