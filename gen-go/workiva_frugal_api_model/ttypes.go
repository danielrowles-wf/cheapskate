// Autogenerated by Thrift Compiler (0.9.3-wk-2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package workiva_frugal_api_model

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

//The condition of a service, as returned in a ServiceHealthStatus
type HealthCondition int64

const (
	HealthCondition_PASS    HealthCondition = 1
	HealthCondition_WARN    HealthCondition = 2
	HealthCondition_FAIL    HealthCondition = 3
	HealthCondition_UNKNOWN HealthCondition = 4
)

func (p HealthCondition) String() string {
	switch p {
	case HealthCondition_PASS:
		return "PASS"
	case HealthCondition_WARN:
		return "WARN"
	case HealthCondition_FAIL:
		return "FAIL"
	case HealthCondition_UNKNOWN:
		return "UNKNOWN"
	}
	return "<UNSET>"
}

func HealthConditionFromString(s string) (HealthCondition, error) {
	switch s {
	case "PASS":
		return HealthCondition_PASS, nil
	case "WARN":
		return HealthCondition_WARN, nil
	case "FAIL":
		return HealthCondition_FAIL, nil
	case "UNKNOWN":
		return HealthCondition_UNKNOWN, nil
	}
	return HealthCondition(0), fmt.Errorf("not a valid HealthCondition string")
}

func HealthConditionPtr(v HealthCondition) *HealthCondition { return &v }

func (p HealthCondition) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *HealthCondition) UnmarshalText(text []byte) error {
	q, err := HealthConditionFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *HealthCondition) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = HealthCondition(v)
	return nil
}

func (p *HealthCondition) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

//Status codes for service health
//DEPRECATED: replaced by HealthCondition for new checkServiceHealth() endpoint
type Status int64

const (
	Status_STARTING        Status = 1
	Status_STOPPING        Status = 2
	Status_HEALTHY         Status = 3
	Status_DEPENDENCY_DOWN Status = 4
	Status_ERROR           Status = 5
)

func (p Status) String() string {
	switch p {
	case Status_STARTING:
		return "STARTING"
	case Status_STOPPING:
		return "STOPPING"
	case Status_HEALTHY:
		return "HEALTHY"
	case Status_DEPENDENCY_DOWN:
		return "DEPENDENCY_DOWN"
	case Status_ERROR:
		return "ERROR"
	}
	return "<UNSET>"
}

func StatusFromString(s string) (Status, error) {
	switch s {
	case "STARTING":
		return Status_STARTING, nil
	case "STOPPING":
		return Status_STOPPING, nil
	case "HEALTHY":
		return Status_HEALTHY, nil
	case "DEPENDENCY_DOWN":
		return Status_DEPENDENCY_DOWN, nil
	case "ERROR":
		return Status_ERROR, nil
	}
	return Status(0), fmt.Errorf("not a valid Status string")
}

func StatusPtr(v Status) *Status { return &v }

func (p Status) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Status) UnmarshalText(text []byte) error {
	q, err := StatusFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *Status) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = Status(v)
	return nil
}

func (p *Status) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Information about the service
//
// Attributes:
//  - Name: Name of the service
//  - Version: Deployed version of the service
//  - Repo: Repository of service source code
//  - ActiveRequests: Number of in-flight requests the service is currently processing
//  - Metadata: Additional metadata about the service. See individual service definitions.
type Info struct {
	Name           string            `thrift:"name,1" db:"name" json:"name"`
	Version        string            `thrift:"version,2" db:"version" json:"version"`
	Repo           string            `thrift:"repo,3" db:"repo" json:"repo"`
	ActiveRequests *int64            `thrift:"activeRequests,4" db:"activeRequests" json:"activeRequests,omitempty"`
	Metadata       map[string]string `thrift:"metadata,5" db:"metadata" json:"metadata,omitempty"`
}

func NewInfo() *Info {
	return &Info{}
}

func (p *Info) GetName() string {
	return p.Name
}

func (p *Info) GetVersion() string {
	return p.Version
}

func (p *Info) GetRepo() string {
	return p.Repo
}

var Info_ActiveRequests_DEFAULT int64

func (p *Info) GetActiveRequests() int64 {
	if !p.IsSetActiveRequests() {
		return Info_ActiveRequests_DEFAULT
	}
	return *p.ActiveRequests
}

var Info_Metadata_DEFAULT map[string]string

func (p *Info) GetMetadata() map[string]string {
	return p.Metadata
}
func (p *Info) IsSetActiveRequests() bool {
	return p.ActiveRequests != nil
}

func (p *Info) IsSetMetadata() bool {
	return p.Metadata != nil
}

func (p *Info) Read(iprot thrift.TProtocol) error {
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
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.ReadField5(iprot); err != nil {
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

func (p *Info) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *Info) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Version = v
	}
	return nil
}

func (p *Info) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Repo = v
	}
	return nil
}

func (p *Info) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.ActiveRequests = &v
	}
	return nil
}

func (p *Info) ReadField5(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Metadata = tMap
	for i := 0; i < size; i++ {
		var _key0 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key0 = v
		}
		var _val1 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val1 = v
		}
		p.Metadata[_key0] = _val1
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *Info) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Info"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
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

func (p *Info) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:name: ", p), err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.name (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:name: ", p), err)
	}
	return err
}

func (p *Info) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("version", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:version: ", p), err)
	}
	if err := oprot.WriteString(string(p.Version)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.version (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:version: ", p), err)
	}
	return err
}

func (p *Info) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("repo", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:repo: ", p), err)
	}
	if err := oprot.WriteString(string(p.Repo)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.repo (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:repo: ", p), err)
	}
	return err
}

func (p *Info) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetActiveRequests() {
		if err := oprot.WriteFieldBegin("activeRequests", thrift.I64, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:activeRequests: ", p), err)
		}
		if err := oprot.WriteI64(int64(*p.ActiveRequests)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.activeRequests (4) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:activeRequests: ", p), err)
		}
	}
	return err
}

func (p *Info) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetMetadata() {
		if err := oprot.WriteFieldBegin("metadata", thrift.MAP, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:metadata: ", p), err)
		}
		if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Metadata)); err != nil {
			return thrift.PrependError("error writing map begin: ", err)
		}
		for k, v := range p.Metadata {
			if err := oprot.WriteString(string(k)); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
			}
			if err := oprot.WriteString(string(v)); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
			}
		}
		if err := oprot.WriteMapEnd(); err != nil {
			return thrift.PrependError("error writing map end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:metadata: ", p), err)
		}
	}
	return err
}

func (p *Info) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Info(%+v)", *p)
}

// A description of the health of a service
//
// Attributes:
//  - Version: The version of the service
//  - Status: The health status of the service
//  - Message: A descriptive message that elaborates upon a status
//  - Metadata: Additional service-specific metadata
type ServiceHealthStatus struct {
	Version  string            `thrift:"version,1,required" db:"version" json:"version"`
	Status   HealthCondition   `thrift:"status,2,required" db:"status" json:"status"`
	Message  string            `thrift:"message,3,required" db:"message" json:"message"`
	Metadata map[string]string `thrift:"metadata,4" db:"metadata" json:"metadata,omitempty"`
}

func NewServiceHealthStatus() *ServiceHealthStatus {
	return &ServiceHealthStatus{}
}

func (p *ServiceHealthStatus) GetVersion() string {
	return p.Version
}

func (p *ServiceHealthStatus) GetStatus() HealthCondition {
	return p.Status
}

func (p *ServiceHealthStatus) GetMessage() string {
	return p.Message
}

var ServiceHealthStatus_Metadata_DEFAULT map[string]string

func (p *ServiceHealthStatus) GetMetadata() map[string]string {
	return p.Metadata
}
func (p *ServiceHealthStatus) IsSetMetadata() bool {
	return p.Metadata != nil
}

func (p *ServiceHealthStatus) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetVersion bool = false
	var issetStatus bool = false
	var issetMessage bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
			issetVersion = true
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
			issetStatus = true
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
			issetMessage = true
		case 4:
			if err := p.ReadField4(iprot); err != nil {
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
	if !issetVersion {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Version is not set"))
	}
	if !issetStatus {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Status is not set"))
	}
	if !issetMessage {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Message is not set"))
	}
	return nil
}

func (p *ServiceHealthStatus) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Version = v
	}
	return nil
}

func (p *ServiceHealthStatus) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		temp := HealthCondition(v)
		p.Status = temp
	}
	return nil
}

func (p *ServiceHealthStatus) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Message = v
	}
	return nil
}

func (p *ServiceHealthStatus) ReadField4(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Metadata = tMap
	for i := 0; i < size; i++ {
		var _key2 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key2 = v
		}
		var _val3 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val3 = v
		}
		p.Metadata[_key2] = _val3
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *ServiceHealthStatus) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ServiceHealthStatus"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
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

func (p *ServiceHealthStatus) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("version", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:version: ", p), err)
	}
	if err := oprot.WriteString(string(p.Version)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.version (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:version: ", p), err)
	}
	return err
}

func (p *ServiceHealthStatus) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("status", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:status: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Status)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.status (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:status: ", p), err)
	}
	return err
}

func (p *ServiceHealthStatus) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("message", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:message: ", p), err)
	}
	if err := oprot.WriteString(string(p.Message)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.message (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:message: ", p), err)
	}
	return err
}

func (p *ServiceHealthStatus) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetMetadata() {
		if err := oprot.WriteFieldBegin("metadata", thrift.MAP, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:metadata: ", p), err)
		}
		if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Metadata)); err != nil {
			return thrift.PrependError("error writing map begin: ", err)
		}
		for k, v := range p.Metadata {
			if err := oprot.WriteString(string(k)); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
			}
			if err := oprot.WriteString(string(v)); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
			}
		}
		if err := oprot.WriteMapEnd(); err != nil {
			return thrift.PrependError("error writing map end: ", err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:metadata: ", p), err)
		}
	}
	return err
}

func (p *ServiceHealthStatus) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ServiceHealthStatus(%+v)", *p)
}

// An error arising while trying to serve a BaseService request
//
// Attributes:
//  - Code: The specific problem flagged by this error, if any
//  - Message: A descriptive error message
type BaseError struct {
	Code    int32  `thrift:"code,1" db:"code" json:"code"`
	Message string `thrift:"message,2" db:"message" json:"message"`
}

func NewBaseError() *BaseError {
	return &BaseError{}
}

func (p *BaseError) GetCode() int32 {
	return p.Code
}

func (p *BaseError) GetMessage() string {
	return p.Message
}
func (p *BaseError) Read(iprot thrift.TProtocol) error {
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
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
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

func (p *BaseError) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Code = v
	}
	return nil
}

func (p *BaseError) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Message = v
	}
	return nil
}

func (p *BaseError) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("BaseError"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
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

func (p *BaseError) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("code", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:code: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Code)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.code (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:code: ", p), err)
	}
	return err
}

func (p *BaseError) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("message", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:message: ", p), err)
	}
	if err := oprot.WriteString(string(p.Message)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.message (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:message: ", p), err)
	}
	return err
}

func (p *BaseError) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseError(%+v)", *p)
}

func (p *BaseError) Error() string {
	return p.String()
}

// Describes the health of the service
// DEPRECATED: replaced by ServiceHealthStatus for new checkServiceHealth() endpoint
//
// Attributes:
//  - Status: Status code
//  - Message: Descriptive message explaining the status
type Health struct {
	Status  Status `thrift:"status,1" db:"status" json:"status"`
	Message string `thrift:"message,2" db:"message" json:"message"`
}

func NewHealth() *Health {
	return &Health{}
}

func (p *Health) GetStatus() Status {
	return p.Status
}

func (p *Health) GetMessage() string {
	return p.Message
}
func (p *Health) Read(iprot thrift.TProtocol) error {
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
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
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

func (p *Health) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := Status(v)
		p.Status = temp
	}
	return nil
}

func (p *Health) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Message = v
	}
	return nil
}

func (p *Health) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Health"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
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

func (p *Health) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("status", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:status: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Status)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.status (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:status: ", p), err)
	}
	return err
}

func (p *Health) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("message", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:message: ", p), err)
	}
	if err := oprot.WriteString(string(p.Message)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.message (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:message: ", p), err)
	}
	return err
}

func (p *Health) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Health(%+v)", *p)
}
