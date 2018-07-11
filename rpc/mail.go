package rpc

import (
	"bytes"
	"encoding/gob"
)

type Mail struct {
	Type uint32

	Object interface{}
	data []byte
}

type MailHandler interface {
	OnNewMail(mail Mail)
	OnRoutineConnected(remote string)
	OnRoutineDisconnect(remote string, err error)
}

func (this* Mail) Encode() {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(this.Object)
	this.data = buffer.Bytes()
}

func (this* Mail) Decode(v interface{}) error {
	buffer := bytes.NewBuffer(this.data)
	dec := gob.NewDecoder(buffer)
	return dec.Decode(v)
}

func (this* Mail) Dump () []byte{
	return this.data
}
