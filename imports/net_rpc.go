// this file was generated by gomacro command: import _b "net/rpc"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"net/rpc"
)

// reflection: allow interpreted code to import "net/rpc"
func init() {
	Packages["net/rpc"] = Package{
	Binds: map[string]Value{
		"Accept":	ValueOf(rpc.Accept),
		"DefaultDebugPath":	ValueOf(rpc.DefaultDebugPath),
		"DefaultRPCPath":	ValueOf(rpc.DefaultRPCPath),
		"DefaultServer":	ValueOf(&rpc.DefaultServer).Elem(),
		"Dial":	ValueOf(rpc.Dial),
		"DialHTTP":	ValueOf(rpc.DialHTTP),
		"DialHTTPPath":	ValueOf(rpc.DialHTTPPath),
		"ErrShutdown":	ValueOf(&rpc.ErrShutdown).Elem(),
		"HandleHTTP":	ValueOf(rpc.HandleHTTP),
		"NewClient":	ValueOf(rpc.NewClient),
		"NewClientWithCodec":	ValueOf(rpc.NewClientWithCodec),
		"NewServer":	ValueOf(rpc.NewServer),
		"Register":	ValueOf(rpc.Register),
		"RegisterName":	ValueOf(rpc.RegisterName),
		"ServeCodec":	ValueOf(rpc.ServeCodec),
		"ServeConn":	ValueOf(rpc.ServeConn),
		"ServeRequest":	ValueOf(rpc.ServeRequest),
	},
	Types: map[string]Type{
		"Call":	TypeOf((*rpc.Call)(nil)).Elem(),
		"Client":	TypeOf((*rpc.Client)(nil)).Elem(),
		"ClientCodec":	TypeOf((*rpc.ClientCodec)(nil)).Elem(),
		"Request":	TypeOf((*rpc.Request)(nil)).Elem(),
		"Response":	TypeOf((*rpc.Response)(nil)).Elem(),
		"Server":	TypeOf((*rpc.Server)(nil)).Elem(),
		"ServerCodec":	TypeOf((*rpc.ServerCodec)(nil)).Elem(),
		"ServerError":	TypeOf((*rpc.ServerError)(nil)).Elem(),
	},
	Proxies: map[string]Type{
		"ClientCodec":	TypeOf((*ClientCodec_net_rpc)(nil)).Elem(),
		"ServerCodec":	TypeOf((*ServerCodec_net_rpc)(nil)).Elem(),
	} }
}

// --------------- proxy for net/rpc.ClientCodec ---------------
type ClientCodec_net_rpc struct {
	Object	interface{}
	Close_	func() error
	ReadResponseBody_	func(interface{}) error
	ReadResponseHeader_	func(*rpc.Response) error
	WriteRequest_	func(*rpc.Request, interface{}) error
}
func (Proxy *ClientCodec_net_rpc) Close() error {
	return Proxy.Close_()
}
func (Proxy *ClientCodec_net_rpc) ReadResponseBody(unnamed0 interface{}) error {
	return Proxy.ReadResponseBody_(unnamed0)
}
func (Proxy *ClientCodec_net_rpc) ReadResponseHeader(unnamed0 *rpc.Response) error {
	return Proxy.ReadResponseHeader_(unnamed0)
}
func (Proxy *ClientCodec_net_rpc) WriteRequest(unnamed0 *rpc.Request, unnamed1 interface{}) error {
	return Proxy.WriteRequest_(unnamed0, unnamed1)
}

// --------------- proxy for net/rpc.ServerCodec ---------------
type ServerCodec_net_rpc struct {
	Object	interface{}
	Close_	func() error
	ReadRequestBody_	func(interface{}) error
	ReadRequestHeader_	func(*rpc.Request) error
	WriteResponse_	func(*rpc.Response, interface{}) error
}
func (Proxy *ServerCodec_net_rpc) Close() error {
	return Proxy.Close_()
}
func (Proxy *ServerCodec_net_rpc) ReadRequestBody(unnamed0 interface{}) error {
	return Proxy.ReadRequestBody_(unnamed0)
}
func (Proxy *ServerCodec_net_rpc) ReadRequestHeader(unnamed0 *rpc.Request) error {
	return Proxy.ReadRequestHeader_(unnamed0)
}
func (Proxy *ServerCodec_net_rpc) WriteResponse(unnamed0 *rpc.Response, unnamed1 interface{}) error {
	return Proxy.WriteResponse_(unnamed0, unnamed1)
}
