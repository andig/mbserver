package mbserver

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

// const (
// 	maxTCPFrameLength int = 260
// 	mbapHeaderLength  int = 7
// )

type udpTransport struct {
	logger    *logger
	conn      *net.UDPConn
	addr      *net.UDPAddr
	timeout   time.Duration
	lastTxnId uint16
}

// Returns a new TCP transport.
func newUDPTransport(conn *net.UDPConn, addr *net.UDPAddr, timeout time.Duration) *udpTransport {
	return &udpTransport{
		conn:    conn,
		addr:    addr,
		timeout: timeout,
		logger:  newLogger(fmt.Sprintf("udp-transport(%s)", addr.String())),
	}
}

// Closes the underlying tcp addr.
func (tt *udpTransport) Close() error {
	return nil
}

// Runs a request across the addr and returns a response.
func (tt *udpTransport) ExecuteRequest(req *pdu) (*pdu, error) {
	return nil, errors.New("foo")
	// // set an i/o deadline on the addr (read and write)
	// if err := tt.addr.SetDeadline(time.Now().Add(tt.timeout)); err != nil {
	// 	return nil, err
	// }

	// // increase the transaction ID counter
	// tt.lastTxnId++

	// if _, err := tt.addr.Write(tt.assembleMBAPFrame(tt.lastTxnId, req)); err != nil {
	// 	return nil, err
	// }

	// return tt.readResponse()
}

// Reads a request from the addr.
func (tt *udpTransport) ReadRequest() (*pdu, error) {
	// set an i/o deadline on the addr (read and write)
	// if err := tt.addr.SetDeadline(time.Now().Add(tt.timeout)); err != nil {
	// 	return nil, err
	// }

	// req, txnId, err := tt.readMBAPFrame()
	// if err != nil {
	// 	return nil, err
	// }

	// // store the incoming transaction id
	// tt.lastTxnId = txnId

	// return req, err
	return nil, errors.New("foo")
}

// Reads a request from the addr.
func (tt *udpTransport) HandleRequest(b []byte) (*pdu, error) {
	// // set an i/o deadline on the addr (read and write)
	// if err := tt.addr.SetDeadline(time.Now().Add(tt.timeout)); err != nil {
	// 	return nil, err
	// }

	req, txnId, err := tt.readMBAPFrame(b)
	if err != nil {
		return nil, err
	}

	// store the incoming transaction id
	tt.lastTxnId = txnId

	fmt.Println("HandleRequest", req, err, txnId)

	return req, err
}

// Writes a response to the addr.
func (tt *udpTransport) WriteResponse(res *pdu) error {
	_, err := tt.conn.WriteToUDP(tt.assembleMBAPFrame(tt.lastTxnId, res), tt.addr)
	fmt.Println("WriteResponse", res, err, tt.lastTxnId)
	return err
}

// Reads as many MBAP+modbus frames as necessary until either the response
// matching tt.lastTxnId is received or an error occurs.
func (tt *udpTransport) readResponse() (res *pdu, err error) {
	// var txnId uint16

	// for {
	// 	// grab a frame
	// 	res, txnId, err = tt.readMBAPFrame()

	// 	// ignore unknown protocol identifiers
	// 	if err == ErrUnknownProtocolId {
	// 		continue
	// 	}

	// 	// abort on any other erorr
	// 	if err != nil {
	// 		return
	// 	}

	// 	// ignore unknown transaction identifiers
	// 	if tt.lastTxnId != txnId {
	// 		tt.logger.Warningf("received unexpected transaction id (expected 0x%04x, received 0x%04x)", tt.lastTxnId, txnId)
	// 		continue
	// 	}

	// 	return
	// }
	return nil, errors.New("foo")
}

// Reads an entire frame (MBAP header + modbus PDU) from the addr.
func (tt *udpTransport) readMBAPFrame(rxbuf []byte) (p *pdu, txnId uint16, err error) {
	// read the MBAP header
	// rxbuf := make([]byte, mbapHeaderLength)
	// if _, err = io.ReadFull(tt.addr, rxbuf); err != nil {
	// 	return
	// }
	fmt.Println("rxbuf 1", string(rxbuf))

	// decode the transaction identifier
	txnId = binary.BigEndian.Uint16(rxbuf[0:2])
	// decode the protocol identifier
	protocolId := binary.BigEndian.Uint16(rxbuf[2:4])
	// store the source unit id
	unitId := rxbuf[6]

	// determine how many more bytes we need to read
	bytesNeeded := int(binary.BigEndian.Uint16(rxbuf[4:6]))

	// the byte count includes the unit ID field, which we already have
	bytesNeeded--

	// never read more than the max allowed frame length
	if bytesNeeded+mbapHeaderLength > maxTCPFrameLength {
		err = ErrProtocolError
		return
	}

	// an MBAP length of 0 is illegal
	if bytesNeeded <= 0 {
		err = ErrProtocolError
		return
	}

	// read the PDU
	// rxbuf = make([]byte, bytesNeeded)
	// if _, err = io.ReadFull(tt.addr, rxbuf); err != nil {
	// 	return
	// }
	rxbuf = rxbuf[mbapHeaderLength:]
	fmt.Println("rxbuf 2", string(rxbuf), len(rxbuf))

	// validate the protocol identifier
	if protocolId != 0x0000 {
		err = ErrUnknownProtocolId
		tt.logger.Warningf("received unexpected protocol id 0x%04x", protocolId)
		return
	}

	// store unit id, function code and payload in the PDU object
	p = &pdu{
		unitId:       unitId,
		functionCode: rxbuf[0],
		payload:      rxbuf[1:],
	}

	return
}

// Turns a PDU into an MBAP frame (MBAP header + PDU) and returns it as bytes.
func (tt *udpTransport) assembleMBAPFrame(txnId uint16, p *pdu) []byte {
	payload := make([]byte, 0, 8+len(p.payload))
	// transaction identifier
	payload = append(payload, asBytes(binary.BigEndian, txnId)...)
	// protocol identifier (always 0x0000)
	payload = append(payload, 0x00, 0x00)
	// length (covers unit identifier + function code + payload fields)
	payload = append(payload, asBytes(binary.BigEndian, uint16(2+len(p.payload)))...)
	// unit identifier
	payload = append(payload, p.unitId)
	// function code
	payload = append(payload, p.functionCode)
	// payload
	payload = append(payload, p.payload...)

	return payload
}
