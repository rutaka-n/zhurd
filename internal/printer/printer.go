package printer

import (
	"errors"
	"net"
)

type Printable interface {
	Print(pType string) ([]byte, error)
}

type Printer struct {
	addr        string
	pType       string
	description string
	conn        net.Conn
	isConnected bool
}

func New(pType, addr, descr string) Printer {
	return Printer{
		addr:        addr,
		pType:       pType,
		description: descr,
	}
}

func (p *Printer) Connect() error {
	conn, err := net.Dial("tcp", p.addr)
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
	bs, err := label.Print(p.pType)
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
