// Autogenerated by Thrift Compiler (0.9.3-wk-2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package stingy

import (
	"bytes"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/danielrowles-wf/cheapskate/gen-go/workiva_frugal_api"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var _ = workiva_frugal_api.GoUnusedProtection__

type StingyService interface {
	workiva_frugal_api.BaseService

	GetQuote() (r string, err error)
}

type StingyServiceClient struct {
	*workiva_frugal_api.BaseServiceClient
}

func NewStingyServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *StingyServiceClient {
	return &StingyServiceClient{BaseServiceClient: workiva_frugal_api.NewBaseServiceClientFactory(t, f)}
}

func NewStingyServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *StingyServiceClient {
	return &StingyServiceClient{BaseServiceClient: workiva_frugal_api.NewBaseServiceClientProtocol(t, iprot, oprot)}
}

func (p *StingyServiceClient) GetQuote() (r string, err error) {
	if err = p.sendGetQuote(); err != nil {
		return
	}
	return p.recvGetQuote()
}

func (p *StingyServiceClient) sendGetQuote() (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("GetQuote", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := StingyServiceGetQuoteArgs{}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *StingyServiceClient) recvGetQuote() (value string, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "GetQuote" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "GetQuote failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "GetQuote failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error1 error
		error1, err = error0.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error1
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "GetQuote failed: invalid message type")
		return
	}
	result := StingyServiceGetQuoteResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type StingyServiceProcessor struct {
	*workiva_frugal_api.BaseServiceProcessor
}

func NewStingyServiceProcessor(handler StingyService) *StingyServiceProcessor {
	self2 := &StingyServiceProcessor{workiva_frugal_api.NewBaseServiceProcessor(handler)}
	self2.AddToProcessorMap("GetQuote", &stingyServiceProcessorGetQuote{handler: handler})
	return self2
}

type stingyServiceProcessorGetQuote struct {
	handler StingyService
}

func (p *stingyServiceProcessorGetQuote) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := StingyServiceGetQuoteArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("GetQuote", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := StingyServiceGetQuoteResult{}
	var retval string
	var err2 error
	if retval, err2 = p.handler.GetQuote(); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing GetQuote: "+err2.Error())
		oprot.WriteMessageBegin("GetQuote", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = &retval
	}
	if err2 = oprot.WriteMessageBegin("GetQuote", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type StingyServiceGetQuoteArgs struct {
}

func NewStingyServiceGetQuoteArgs() *StingyServiceGetQuoteArgs {
	return &StingyServiceGetQuoteArgs{}
}

func (p *StingyServiceGetQuoteArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if err := iprot.Skip(fieldTypeId); err != nil {
			return err
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StingyServiceGetQuoteArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("GetQuote_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StingyServiceGetQuoteArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StingyServiceGetQuoteArgs(%+v)", *p)
}

// Attributes:
//  - Success
type StingyServiceGetQuoteResult struct {
	Success *string `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewStingyServiceGetQuoteResult() *StingyServiceGetQuoteResult {
	return &StingyServiceGetQuoteResult{}
}

var StingyServiceGetQuoteResult_Success_DEFAULT string

func (p *StingyServiceGetQuoteResult) GetSuccess() string {
	if !p.IsSetSuccess() {
		return StingyServiceGetQuoteResult_Success_DEFAULT
	}
	return *p.Success
}
func (p *StingyServiceGetQuoteResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *StingyServiceGetQuoteResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.ReadField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StingyServiceGetQuoteResult) ReadField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 0: ", err)
	} else {
		p.Success = &v
	}
	return nil
}

func (p *StingyServiceGetQuoteResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("GetQuote_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StingyServiceGetQuoteResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRING, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := oprot.WriteString(string(*p.Success)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.success (0) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *StingyServiceGetQuoteResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StingyServiceGetQuoteResult(%+v)", *p)
}
