// this file was generated by gomacro command: import _b "reflect"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"reflect"
)

// reflection: allow interpreted code to import "reflect"
func init() {
	Packages["reflect"] = Package{
	Binds: map[string]Value{
		"Append":	ValueOf(reflect.Append),
		"AppendSlice":	ValueOf(reflect.AppendSlice),
		"Array":	ValueOf(reflect.Array),
		"ArrayOf":	ValueOf(reflect.ArrayOf),
		"Bool":	ValueOf(reflect.Bool),
		"BothDir":	ValueOf(reflect.BothDir),
		"Chan":	ValueOf(reflect.Chan),
		"ChanOf":	ValueOf(reflect.ChanOf),
		"Complex128":	ValueOf(reflect.Complex128),
		"Complex64":	ValueOf(reflect.Complex64),
		"Copy":	ValueOf(reflect.Copy),
		"DeepEqual":	ValueOf(reflect.DeepEqual),
		"Float32":	ValueOf(reflect.Float32),
		"Float64":	ValueOf(reflect.Float64),
		"Func":	ValueOf(reflect.Func),
		"FuncOf":	ValueOf(reflect.FuncOf),
		"Indirect":	ValueOf(reflect.Indirect),
		"Int":	ValueOf(reflect.Int),
		"Int16":	ValueOf(reflect.Int16),
		"Int32":	ValueOf(reflect.Int32),
		"Int64":	ValueOf(reflect.Int64),
		"Int8":	ValueOf(reflect.Int8),
		"Interface":	ValueOf(reflect.Interface),
		"Invalid":	ValueOf(reflect.Invalid),
		"MakeChan":	ValueOf(reflect.MakeChan),
		"MakeFunc":	ValueOf(reflect.MakeFunc),
		"MakeMap":	ValueOf(reflect.MakeMap),
		"MakeSlice":	ValueOf(reflect.MakeSlice),
		"Map":	ValueOf(reflect.Map),
		"MapOf":	ValueOf(reflect.MapOf),
		"New":	ValueOf(reflect.New),
		"NewAt":	ValueOf(reflect.NewAt),
		"Ptr":	ValueOf(reflect.Ptr),
		"PtrTo":	ValueOf(reflect.PtrTo),
		"RecvDir":	ValueOf(reflect.RecvDir),
		"Select":	ValueOf(reflect.Select),
		"SelectDefault":	ValueOf(reflect.SelectDefault),
		"SelectRecv":	ValueOf(reflect.SelectRecv),
		"SelectSend":	ValueOf(reflect.SelectSend),
		"SendDir":	ValueOf(reflect.SendDir),
		"Slice":	ValueOf(reflect.Slice),
		"SliceOf":	ValueOf(reflect.SliceOf),
		"String":	ValueOf(reflect.String),
		"Struct":	ValueOf(reflect.Struct),
		"StructOf":	ValueOf(reflect.StructOf),
		"Swapper":	ValueOf(reflect.Swapper),
		"TypeOf":	ValueOf(reflect.TypeOf),
		"Uint":	ValueOf(reflect.Uint),
		"Uint16":	ValueOf(reflect.Uint16),
		"Uint32":	ValueOf(reflect.Uint32),
		"Uint64":	ValueOf(reflect.Uint64),
		"Uint8":	ValueOf(reflect.Uint8),
		"Uintptr":	ValueOf(reflect.Uintptr),
		"UnsafePointer":	ValueOf(reflect.UnsafePointer),
		"ValueOf":	ValueOf(reflect.ValueOf),
		"Zero":	ValueOf(reflect.Zero),
	},
	Types: map[string]Type{
		"ChanDir":	TypeOf((*reflect.ChanDir)(nil)).Elem(),
		"Kind":	TypeOf((*reflect.Kind)(nil)).Elem(),
		"Method":	TypeOf((*reflect.Method)(nil)).Elem(),
		"SelectCase":	TypeOf((*reflect.SelectCase)(nil)).Elem(),
		"SelectDir":	TypeOf((*reflect.SelectDir)(nil)).Elem(),
		"SliceHeader":	TypeOf((*reflect.SliceHeader)(nil)).Elem(),
		"StringHeader":	TypeOf((*reflect.StringHeader)(nil)).Elem(),
		"StructField":	TypeOf((*reflect.StructField)(nil)).Elem(),
		"StructTag":	TypeOf((*reflect.StructTag)(nil)).Elem(),
		"Type":	TypeOf((*reflect.Type)(nil)).Elem(),
		"Value":	TypeOf((*reflect.Value)(nil)).Elem(),
		"ValueError":	TypeOf((*reflect.ValueError)(nil)).Elem(),
	},
	Proxies: map[string]Type{
	} }
}
