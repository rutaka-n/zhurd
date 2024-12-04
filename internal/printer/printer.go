package printer

import (
	"errors"
	"log/slog"
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
	if p.isConnected {
		slog.Warn("printer alredy has established connection, ignored", "ID", p.ID)
	}
	resolvedAddr, err := net.ResolveTCPAddr("tcp", p.Addr)
	if err != nil {
		return err
	}
	p.conn, err = net.DialTCP("tcp", nil, resolvedAddr)
	if err != nil {
		return err
	}
	p.isConnected = true
	return nil
}

func (p *Printer) Close() error {
	if !p.isConnected {
		return nil
	}
	p.isConnected = false
	return p.conn.Close()
}

func (p *Printer) IsConnected() bool {
	return p.isConnected
}

func (p *Printer) Enqueue(label Printable) error {
	if !p.isConnected {
		return net.ErrClosed
	}
	bs, err := label.Print(p.Type)
	if err != nil {
		return err
	}

	written, err := p.conn.Write(bs)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			p.isConnected = false
		}
		return err
	}
	if written < len(bs) {
		slog.Warn("document was not fully send to printer", "size", len(bs), "bytesSent", written)
	}

	return nil
}
