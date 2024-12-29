package printer

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PSQL struct {
	pool *pgxpool.Pool
}

func NewPSQL(pool *pgxpool.Pool) (*PSQL, error) {
	return &PSQL{pool: pool}, nil
}

func (repo *PSQL) Store(ctx context.Context, p *Printer) error {
	sql := "INSERT INTO printers (addr, type, comment) VALUES ($1, $2, $3) RETURNING id"
	row := repo.pool.QueryRow(ctx, sql, p.Addr, p.Type, p.Comment)
	if err := row.Scan(&p.ID); err != nil {
		return err
	}
	return nil
}

func (repo *PSQL) List(ctx context.Context) ([]Printer, error) {
	sql := "SELECT id, addr, type, comment FROM printers"
	rows, err := repo.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	printers := []Printer{}
	for rows.Next() {
		p := Printer{}
		if err := rows.Scan(&p.ID, &p.Addr, &p.Type, &p.Comment); err != nil {
			return nil, err
		}
		printers = append(printers, p)
	}
	return printers, nil
}

func (repo *PSQL) Get(ctx context.Context, id int64) (Printer, error) {
	sql := "SELECT id, addr, type, comment FROM printers WHERE id = $1"
	row := repo.pool.QueryRow(ctx, sql, id)
	p := Printer{}
	if err := row.Scan(&p.ID, &p.Addr, &p.Type, &p.Comment); err != nil {
		return Printer{}, err
	}
	return p, nil
}

func (repo *PSQL) Delete(ctx context.Context, id int64) error {
	sql := "DELETE FROM printers WHERE id = $1"
	_, err := repo.pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	return nil
}
