// this file was generated by gomacro command: import _b "plugin"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"plugin"
)

// reflection: allow interpreted code to import "plugin"
func init() {
	Packages["plugin"] = Package{
	Binds: map[string]Value{
		"Open":	ValueOf(plugin.Open),
	},
	Types: map[string]Type{
		"Plugin":	TypeOf((*plugin.Plugin)(nil)).Elem(),
		"Symbol":	TypeOf((*plugin.Symbol)(nil)).Elem(),
	},
	Proxies: map[string]Type{
		"Symbol":	TypeOf((*Symbol_plugin)(nil)).Elem(),
	} }
}

// --------------- proxy for plugin.Symbol ---------------
type Symbol_plugin struct {
	Object	interface{}
}
