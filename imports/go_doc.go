// this file was generated by gomacro command: import _b "go/doc"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"go/doc"
)

// reflection: allow interpreted code to import "go/doc"
func init() {
	Packages["go/doc"] = Package{
	Binds: map[string]Value{
		"AllDecls":	ValueOf(doc.AllDecls),
		"AllMethods":	ValueOf(doc.AllMethods),
		"Examples":	ValueOf(doc.Examples),
		"IllegalPrefixes":	ValueOf(&doc.IllegalPrefixes).Elem(),
		"IsPredeclared":	ValueOf(doc.IsPredeclared),
		"New":	ValueOf(doc.New),
		"Synopsis":	ValueOf(doc.Synopsis),
		"ToHTML":	ValueOf(doc.ToHTML),
		"ToText":	ValueOf(doc.ToText),
	},
	Types: map[string]Type{
		"Example":	TypeOf((*doc.Example)(nil)).Elem(),
		"Filter":	TypeOf((*doc.Filter)(nil)).Elem(),
		"Func":	TypeOf((*doc.Func)(nil)).Elem(),
		"Mode":	TypeOf((*doc.Mode)(nil)).Elem(),
		"Note":	TypeOf((*doc.Note)(nil)).Elem(),
		"Package":	TypeOf((*doc.Package)(nil)).Elem(),
		"Type":	TypeOf((*doc.Type)(nil)).Elem(),
		"Value":	TypeOf((*doc.Value)(nil)).Elem(),
	},
	Proxies: map[string]Type{
	} }
}
