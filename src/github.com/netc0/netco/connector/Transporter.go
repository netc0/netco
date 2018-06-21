package connector

type Transporter interface {
	ReadData() ([]byte, int)
	WriteData([]byte) (int, error)
}
