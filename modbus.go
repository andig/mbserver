package mbserver

import (
	"errors"
	"fmt"
)

// maximum PDU size (253 bytes) minus the function code
const maxPDUPayloadLength int = 252

type pdu struct {
	unitId       uint8
	functionCode uint8
	payload      []byte
}

type Error string

// Error implements the error interface.
func (me Error) Error() (s string) {
	s = string(me)
	return
}

// ExceptionError is a modbus exception code, e.g. a vendor specific one
// without a predefined Error. Handlers may return it to send the exception
// code unchanged. ExceptionError(0) is not a valid exception code and is
// sent as a server device failure.
type ExceptionError uint8

// Error implements the error interface.
func (ee ExceptionError) Error() string {
	return fmt.Sprintf("exception code 0x%02x", uint8(ee))
}

const (
	// coils
	fcReadCoils          uint8 = 0x01
	fcWriteSingleCoil    uint8 = 0x05
	fcWriteMultipleCoils uint8 = 0x0f

	// discrete inputs
	fcReadDiscreteInputs uint8 = 0x02

	// 16-bit input/holding registers
	fcReadHoldingRegisters       uint8 = 0x03
	fcReadInputRegisters         uint8 = 0x04
	fcWriteSingleRegister        uint8 = 0x06
	fcWriteMultipleRegisters     uint8 = 0x10
	fcMaskWriteRegister          uint8 = 0x16
	fcReadWriteMultipleRegisters uint8 = 0x17
	fcReadFifoQueue              uint8 = 0x18

	// exception codes
	exIllegalFunction         uint8 = 0x01
	exIllegalDataAddress      uint8 = 0x02
	exIllegalDataValue        uint8 = 0x03
	exServerDeviceFailure     uint8 = 0x04
	exAcknowledge             uint8 = 0x05
	exServerDeviceBusy        uint8 = 0x06
	exMemoryParityError       uint8 = 0x08
	exGWPathUnavailable       uint8 = 0x0a
	exGWTargetFailedToRespond uint8 = 0x0b

	// errors
	ErrConfigurationError      Error = "configuration error"
	ErrRequestTimedOut         Error = "request timed out"
	ErrIllegalFunction         Error = "illegal function"
	ErrIllegalDataAddress      Error = "illegal data address"
	ErrIllegalDataValue        Error = "illegal data value"
	ErrServerDeviceFailure     Error = "server device failure"
	ErrAcknowledge             Error = "request acknowledged"
	ErrServerDeviceBusy        Error = "server device busy"
	ErrMemoryParityError       Error = "memory parity error"
	ErrGWPathUnavailable       Error = "gateway path unavailable"
	ErrGWTargetFailedToRespond Error = "gateway target device failed to respond"
	ErrBadCRC                  Error = "bad crc"
	ErrShortFrame              Error = "short frame"
	ErrProtocolError           Error = "protocol error"
	ErrBadUnitId               Error = "bad unit id"
	ErrBadTransactionId        Error = "bad transaction id"
	ErrUnknownProtocolId       Error = "unknown protocol identifier"
	ErrUnexpectedParameters    Error = "unexpected parameters"
)

// MapExceptionCodeToError turns a modbus exception code into a higher level Error object.
// Codes without a predefined Error are returned as ExceptionError.
func MapExceptionCodeToError(exceptionCode uint8) (err error) {
	switch exceptionCode {
	case exIllegalFunction:
		err = ErrIllegalFunction
	case exIllegalDataAddress:
		err = ErrIllegalDataAddress
	case exIllegalDataValue:
		err = ErrIllegalDataValue
	case exServerDeviceFailure:
		err = ErrServerDeviceFailure
	case exAcknowledge:
		err = ErrAcknowledge
	case exMemoryParityError:
		err = ErrMemoryParityError
	case exServerDeviceBusy:
		err = ErrServerDeviceBusy
	case exGWPathUnavailable:
		err = ErrGWPathUnavailable
	case exGWTargetFailedToRespond:
		err = ErrGWTargetFailedToRespond
	default:
		err = ExceptionError(exceptionCode)
	}

	return
}

// mapErrorToExceptionCode turns an Error object into a modbus exception code.
func mapErrorToExceptionCode(err error) (exceptionCode uint8) {
	// an ExceptionError carries the exception code itself
	var ee ExceptionError
	if errors.As(err, &ee) && ee != 0 {
		return uint8(ee)
	}

	// predefined errors are matched like ExceptionError, also when wrapped
	switch {
	case errors.Is(err, ErrIllegalFunction):
		exceptionCode = exIllegalFunction
	case errors.Is(err, ErrIllegalDataAddress):
		exceptionCode = exIllegalDataAddress
	case errors.Is(err, ErrIllegalDataValue):
		exceptionCode = exIllegalDataValue
	case errors.Is(err, ErrServerDeviceFailure):
		exceptionCode = exServerDeviceFailure
	case errors.Is(err, ErrAcknowledge):
		exceptionCode = exAcknowledge
	case errors.Is(err, ErrMemoryParityError):
		exceptionCode = exMemoryParityError
	case errors.Is(err, ErrServerDeviceBusy):
		exceptionCode = exServerDeviceBusy
	case errors.Is(err, ErrGWPathUnavailable):
		exceptionCode = exGWPathUnavailable
	case errors.Is(err, ErrGWTargetFailedToRespond):
		exceptionCode = exGWTargetFailedToRespond
	default:
		exceptionCode = exServerDeviceFailure
	}

	return
}
