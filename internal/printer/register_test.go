package printer

import (
	"errors"
	"testing"
)

func TestRegister(t *testing.T) {
	ucs := []struct {
		desc        string
		cp          CreatePrinter
		expectedErr error
	}{
		{
			desc: "happy path",
			cp: CreatePrinter{
				Addr:    "0.0.0.0:8009",
				Type:    "ZPL",
				Comment: "test printer",
			},
			expectedErr: nil,
		},
		{
			desc: "invalid port",
			cp: CreatePrinter{
				Addr:    "0.0.256.0:128009",
				Type:    "ZPL",
				Comment: "test printer",
			},
			expectedErr: ValidationError,
		},
	}

	for _, us := range ucs {
		us := us
		t.Run(us.desc, func(t *testing.T) {
			repo, err := NewMemory()
			if err != nil {
				t.Fatalf("got error: %s\n", err)
			}
			svc := NewRegisterSvc(repo)

			p, err := svc.Call(us.cp)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil {
				if p.Addr != us.cp.Addr {
					t.Errorf("expected: %v, got: %v\n", us.cp.Addr, p.Addr)
				}
				if p.Type != us.cp.Type {
					t.Errorf("expected: %v, got: %v\n", us.cp.Type, p.Type)
				}
				if p.Comment != us.cp.Comment {
					t.Errorf("expected: %v, got: %v\n", us.cp.Comment, p.Comment)
				}
			}
		})
	}
}
