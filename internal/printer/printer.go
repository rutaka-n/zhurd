package printer

import (
	"errors"
	"net"
)

type Printable interface {
	Print(pType string) ([]byte, error)
}

type Printer struct {
	ID          int64
	Addr        string
	Type        string
	Comment     string
	conn        net.Conn
	isConnected bool
}

func New(pType, addr, comment string) Printer {
	return Printer{
		Addr:    addr,
		Type:    pType,
		Comment: comment,
	}
}

func (p *Printer) Connect() error {
	conn, err := net.Dial("tcp", p.Addr)
	if err != nil {
		return err
	}
	p.conn = conn
	p.isConnected = true
	return nil
}

func (p *Printer) Close() error {
	p.isConnected = false
	return p.conn.Close()
}

func (p *Printer) Enqueue(label Printable) error {
	if !p.isConnected {
		return net.ErrClosed
	}
	bs, err := label.Print(p.Type)
	if err != nil {
		return err
	}

	_, err = p.conn.Write(bs)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			p.isConnected = false
		}
		return err
	}

	return nil
}
