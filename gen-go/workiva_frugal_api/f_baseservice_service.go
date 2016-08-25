// Autogenerated by Frugal Compiler (1.15.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package workiva_frugal_api

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal/lib/go"
	"github.com/danielrowles-wf/cheapskate/gen-go/workiva_frugal_api_model"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type FBaseService interface {
	// A simple ping to see if the service is alive
	Ping(ctx *frugal.FContext) (err error)
	// Gen2 health check
	CheckServiceHealth(ctx *frugal.FContext) (r *workiva_frugal_api_model.ServiceHealthStatus, err error)
	// Get the current info about the service
	GetInfo(ctx *frugal.FContext) (r *workiva_frugal_api_model.Info, err error)
	// Get the current health of the service.
	// DEPRECATED: replaced by checkServiceHealth()
	GetHealth(ctx *frugal.FContext) (r *workiva_frugal_api_model.Health, err error)
}

type FBaseServiceClient struct {
	transport       frugal.FTransport
	protocolFactory *frugal.FProtocolFactory
	oprot           *frugal.FProtocol
	mu              sync.Mutex
	methods         map[string]*frugal.Method
}

func NewFBaseServiceClient(t frugal.FTransport, p *frugal.FProtocolFactory, middleware ...frugal.ServiceMiddleware) *FBaseServiceClient {
	t.SetRegistry(frugal.NewFClientRegistry())
	methods := make(map[string]*frugal.Method)
	client := &FBaseServiceClient{
		transport:       t,
		protocolFactory: p,
		oprot:           p.GetProtocol(t),
		methods:         methods,
	}
	methods["ping"] = frugal.NewMethod(client, client.ping, "ping", middleware)
	methods["checkServiceHealth"] = frugal.NewMethod(client, client.checkServiceHealth, "checkServiceHealth", middleware)
	methods["getInfo"] = frugal.NewMethod(client, client.getInfo, "getInfo", middleware)
	methods["getHealth"] = frugal.NewMethod(client, client.getHealth, "getHealth", middleware)
	return client
}

// Do Not Use. To be called only by generated code.
func (f *FBaseServiceClient) GetWriteMutex() *sync.Mutex {
	return &f.mu
}

// A simple ping to see if the service is alive
func (f *FBaseServiceClient) Ping(ctx *frugal.FContext) (err error) {
	ret := f.methods["ping"].Invoke([]interface{}{ctx})
	if len(ret) != 1 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 1", len(ret)))
	}
	if ret[0] != nil {
		err = ret[0].(error)
	}
	return err
}

func (f *FBaseServiceClient) ping(ctx *frugal.FContext) (err error) {
	errorC := make(chan error, 1)
	resultC := make(chan struct{}, 1)
	if err = f.transport.Register(ctx, f.recvPingHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	defer f.transport.Unregister(ctx)
	f.GetWriteMutex().Lock()
	if err = f.oprot.WriteRequestHeader(ctx); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageBegin("ping", thrift.CALL, 0); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	args := BaseServicePingArgs{}
	if err = args.Write(f.oprot); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageEnd(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.Flush(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	f.GetWriteMutex().Unlock()

	select {
	case err = <-errorC:
	case <-resultC:
	case <-time.After(ctx.Timeout()):
		err = frugal.ErrTimeout
	case <-f.transport.Closed():
		err = frugal.ErrTransportClosed
	}
	return
}

func (f *FBaseServiceClient) recvPingHandler(ctx *frugal.FContext, resultC chan<- struct{}, errorC chan<- error) frugal.FAsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.protocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "ping" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "ping failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 thrift.TApplicationException
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			if error1.TypeId() == frugal.RESPONSE_TOO_LARGE {
				err = thrift.NewTTransportException(frugal.RESPONSE_TOO_LARGE, "response too large for transport")
				errorC <- err
				return nil
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "ping failed: invalid message type")
			errorC <- err
			return err
		}
		result := BaseServicePingResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		if result.BErr != nil {
			errorC <- result.BErr
			return nil
		}
		resultC <- struct{}{}
		return nil
	}
}

// Gen2 health check
func (f *FBaseServiceClient) CheckServiceHealth(ctx *frugal.FContext) (r *workiva_frugal_api_model.ServiceHealthStatus, err error) {
	ret := f.methods["checkServiceHealth"].Invoke([]interface{}{ctx})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	r = ret[0].(*workiva_frugal_api_model.ServiceHealthStatus)
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return r, err
}

func (f *FBaseServiceClient) checkServiceHealth(ctx *frugal.FContext) (r *workiva_frugal_api_model.ServiceHealthStatus, err error) {
	errorC := make(chan error, 1)
	resultC := make(chan *workiva_frugal_api_model.ServiceHealthStatus, 1)
	if err = f.transport.Register(ctx, f.recvCheckServiceHealthHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	defer f.transport.Unregister(ctx)
	f.GetWriteMutex().Lock()
	if err = f.oprot.WriteRequestHeader(ctx); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageBegin("checkServiceHealth", thrift.CALL, 0); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	args := BaseServiceCheckServiceHealthArgs{}
	if err = args.Write(f.oprot); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageEnd(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.Flush(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	f.GetWriteMutex().Unlock()

	select {
	case err = <-errorC:
	case r = <-resultC:
	case <-time.After(ctx.Timeout()):
		err = frugal.ErrTimeout
	case <-f.transport.Closed():
		err = frugal.ErrTransportClosed
	}
	return
}

func (f *FBaseServiceClient) recvCheckServiceHealthHandler(ctx *frugal.FContext, resultC chan<- *workiva_frugal_api_model.ServiceHealthStatus, errorC chan<- error) frugal.FAsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.protocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "checkServiceHealth" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "checkServiceHealth failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 thrift.TApplicationException
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			if error1.TypeId() == frugal.RESPONSE_TOO_LARGE {
				err = thrift.NewTTransportException(frugal.RESPONSE_TOO_LARGE, "response too large for transport")
				errorC <- err
				return nil
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "checkServiceHealth failed: invalid message type")
			errorC <- err
			return err
		}
		result := BaseServiceCheckServiceHealthResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		if result.BErr != nil {
			errorC <- result.BErr
			return nil
		}
		resultC <- result.GetSuccess()
		return nil
	}
}

// Get the current info about the service
func (f *FBaseServiceClient) GetInfo(ctx *frugal.FContext) (r *workiva_frugal_api_model.Info, err error) {
	ret := f.methods["getInfo"].Invoke([]interface{}{ctx})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	r = ret[0].(*workiva_frugal_api_model.Info)
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return r, err
}

func (f *FBaseServiceClient) getInfo(ctx *frugal.FContext) (r *workiva_frugal_api_model.Info, err error) {
	errorC := make(chan error, 1)
	resultC := make(chan *workiva_frugal_api_model.Info, 1)
	if err = f.transport.Register(ctx, f.recvGetInfoHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	defer f.transport.Unregister(ctx)
	f.GetWriteMutex().Lock()
	if err = f.oprot.WriteRequestHeader(ctx); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageBegin("getInfo", thrift.CALL, 0); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	args := BaseServiceGetInfoArgs{}
	if err = args.Write(f.oprot); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageEnd(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.Flush(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	f.GetWriteMutex().Unlock()

	select {
	case err = <-errorC:
	case r = <-resultC:
	case <-time.After(ctx.Timeout()):
		err = frugal.ErrTimeout
	case <-f.transport.Closed():
		err = frugal.ErrTransportClosed
	}
	return
}

func (f *FBaseServiceClient) recvGetInfoHandler(ctx *frugal.FContext, resultC chan<- *workiva_frugal_api_model.Info, errorC chan<- error) frugal.FAsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.protocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "getInfo" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "getInfo failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 thrift.TApplicationException
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			if error1.TypeId() == frugal.RESPONSE_TOO_LARGE {
				err = thrift.NewTTransportException(frugal.RESPONSE_TOO_LARGE, "response too large for transport")
				errorC <- err
				return nil
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "getInfo failed: invalid message type")
			errorC <- err
			return err
		}
		result := BaseServiceGetInfoResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		if result.BErr != nil {
			errorC <- result.BErr
			return nil
		}
		resultC <- result.GetSuccess()
		return nil
	}
}

// Get the current health of the service.
// DEPRECATED: replaced by checkServiceHealth()
func (f *FBaseServiceClient) GetHealth(ctx *frugal.FContext) (r *workiva_frugal_api_model.Health, err error) {
	ret := f.methods["getHealth"].Invoke([]interface{}{ctx})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	r = ret[0].(*workiva_frugal_api_model.Health)
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return r, err
}

func (f *FBaseServiceClient) getHealth(ctx *frugal.FContext) (r *workiva_frugal_api_model.Health, err error) {
	errorC := make(chan error, 1)
	resultC := make(chan *workiva_frugal_api_model.Health, 1)
	if err = f.transport.Register(ctx, f.recvGetHealthHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	defer f.transport.Unregister(ctx)
	f.GetWriteMutex().Lock()
	if err = f.oprot.WriteRequestHeader(ctx); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageBegin("getHealth", thrift.CALL, 0); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	args := BaseServiceGetHealthArgs{}
	if err = args.Write(f.oprot); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageEnd(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.Flush(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	f.GetWriteMutex().Unlock()

	select {
	case err = <-errorC:
	case r = <-resultC:
	case <-time.After(ctx.Timeout()):
		err = frugal.ErrTimeout
	case <-f.transport.Closed():
		err = frugal.ErrTransportClosed
	}
	return
}

func (f *FBaseServiceClient) recvGetHealthHandler(ctx *frugal.FContext, resultC chan<- *workiva_frugal_api_model.Health, errorC chan<- error) frugal.FAsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.protocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "getHealth" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "getHealth failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 thrift.TApplicationException
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			if error1.TypeId() == frugal.RESPONSE_TOO_LARGE {
				err = thrift.NewTTransportException(frugal.RESPONSE_TOO_LARGE, "response too large for transport")
				errorC <- err
				return nil
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "getHealth failed: invalid message type")
			errorC <- err
			return err
		}
		result := BaseServiceGetHealthResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		if result.BErr != nil {
			errorC <- result.BErr
			return nil
		}
		resultC <- result.GetSuccess()
		return nil
	}
}

type FBaseServiceProcessor struct {
	*frugal.FBaseProcessor
}

func NewFBaseServiceProcessor(handler FBaseService, middleware ...frugal.ServiceMiddleware) *FBaseServiceProcessor {
	p := &FBaseServiceProcessor{frugal.NewFBaseProcessor()}
	p.AddToProcessorMap("ping", &baseserviceFPing{handler: frugal.NewMethod(handler, handler.Ping, "Ping", middleware), writeMu: p.GetWriteMutex()})
	p.AddToProcessorMap("checkServiceHealth", &baseserviceFCheckServiceHealth{handler: frugal.NewMethod(handler, handler.CheckServiceHealth, "CheckServiceHealth", middleware), writeMu: p.GetWriteMutex()})
	p.AddToProcessorMap("getInfo", &baseserviceFGetInfo{handler: frugal.NewMethod(handler, handler.GetInfo, "GetInfo", middleware), writeMu: p.GetWriteMutex()})
	p.AddToProcessorMap("getHealth", &baseserviceFGetHealth{handler: frugal.NewMethod(handler, handler.GetHealth, "GetHealth", middleware), writeMu: p.GetWriteMutex()})
	return p
}

type baseserviceFPing struct {
	handler *frugal.Method
	writeMu *sync.Mutex
}

func (p *baseserviceFPing) Process(ctx *frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	args := BaseServicePingArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.writeMu.Lock()
		baseserviceWriteApplicationError(ctx, oprot, thrift.PROTOCOL_ERROR, "ping", err.Error())
		p.writeMu.Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := BaseServicePingResult{}
	var err2 error
	ret := p.handler.Invoke([]interface{}{ctx})
	if len(ret) != 1 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 1", len(ret)))
	}
	if ret[0] != nil {
		err2 = ret[0].(error)
	}
	if err2 != nil {
		switch v := err2.(type) {
		case *workiva_frugal_api_model.BaseError:
			result.BErr = v
		default:
			p.writeMu.Lock()
			baseserviceWriteApplicationError(ctx, oprot, thrift.INTERNAL_ERROR, "ping", "Internal error processing ping: "+err2.Error())
			p.writeMu.Unlock()
			return err2
		}
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "ping", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("ping", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "ping", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "ping", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "ping", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "ping", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

type baseserviceFCheckServiceHealth struct {
	handler *frugal.Method
	writeMu *sync.Mutex
}

func (p *baseserviceFCheckServiceHealth) Process(ctx *frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	args := BaseServiceCheckServiceHealthArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.writeMu.Lock()
		baseserviceWriteApplicationError(ctx, oprot, thrift.PROTOCOL_ERROR, "checkServiceHealth", err.Error())
		p.writeMu.Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := BaseServiceCheckServiceHealthResult{}
	var err2 error
	var retval *workiva_frugal_api_model.ServiceHealthStatus
	ret := p.handler.Invoke([]interface{}{ctx})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	retval = ret[0].(*workiva_frugal_api_model.ServiceHealthStatus)
	if ret[1] != nil {
		err2 = ret[1].(error)
	}
	if err2 != nil {
		switch v := err2.(type) {
		case *workiva_frugal_api_model.BaseError:
			result.BErr = v
		default:
			p.writeMu.Lock()
			baseserviceWriteApplicationError(ctx, oprot, thrift.INTERNAL_ERROR, "checkServiceHealth", "Internal error processing checkServiceHealth: "+err2.Error())
			p.writeMu.Unlock()
			return err2
		}
	} else {
		result.Success = retval
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "checkServiceHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("checkServiceHealth", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "checkServiceHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "checkServiceHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "checkServiceHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "checkServiceHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

type baseserviceFGetInfo struct {
	handler *frugal.Method
	writeMu *sync.Mutex
}

func (p *baseserviceFGetInfo) Process(ctx *frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	args := BaseServiceGetInfoArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.writeMu.Lock()
		baseserviceWriteApplicationError(ctx, oprot, thrift.PROTOCOL_ERROR, "getInfo", err.Error())
		p.writeMu.Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := BaseServiceGetInfoResult{}
	var err2 error
	var retval *workiva_frugal_api_model.Info
	ret := p.handler.Invoke([]interface{}{ctx})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	retval = ret[0].(*workiva_frugal_api_model.Info)
	if ret[1] != nil {
		err2 = ret[1].(error)
	}
	if err2 != nil {
		switch v := err2.(type) {
		case *workiva_frugal_api_model.BaseError:
			result.BErr = v
		default:
			p.writeMu.Lock()
			baseserviceWriteApplicationError(ctx, oprot, thrift.INTERNAL_ERROR, "getInfo", "Internal error processing getInfo: "+err2.Error())
			p.writeMu.Unlock()
			return err2
		}
	} else {
		result.Success = retval
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getInfo", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("getInfo", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getInfo", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getInfo", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getInfo", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getInfo", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

type baseserviceFGetHealth struct {
	handler *frugal.Method
	writeMu *sync.Mutex
}

func (p *baseserviceFGetHealth) Process(ctx *frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	args := BaseServiceGetHealthArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.writeMu.Lock()
		baseserviceWriteApplicationError(ctx, oprot, thrift.PROTOCOL_ERROR, "getHealth", err.Error())
		p.writeMu.Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := BaseServiceGetHealthResult{}
	var err2 error
	var retval *workiva_frugal_api_model.Health
	ret := p.handler.Invoke([]interface{}{ctx})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	retval = ret[0].(*workiva_frugal_api_model.Health)
	if ret[1] != nil {
		err2 = ret[1].(error)
	}
	if err2 != nil {
		switch v := err2.(type) {
		case *workiva_frugal_api_model.BaseError:
			result.BErr = v
		default:
			p.writeMu.Lock()
			baseserviceWriteApplicationError(ctx, oprot, thrift.INTERNAL_ERROR, "getHealth", "Internal error processing getHealth: "+err2.Error())
			p.writeMu.Unlock()
			return err2
		}
	} else {
		result.Success = retval
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("getHealth", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			baseserviceWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "getHealth", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

func baseserviceWriteApplicationError(ctx *frugal.FContext, oprot *frugal.FProtocol, type_ int32, method, message string) {
	x := thrift.NewTApplicationException(type_, message)
	oprot.WriteResponseHeader(ctx)
	oprot.WriteMessageBegin(method, thrift.EXCEPTION, 0)
	x.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
}
