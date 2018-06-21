package connector

type UDPTransporter struct{
	hostPort string
	running bool
	sessions map[string]*Session
	_onNewConnection func(string, *Session)
	_onCloseConnection func(string)
	_onConnectionData func(*Session, uint32, uint32, []byte)
}
