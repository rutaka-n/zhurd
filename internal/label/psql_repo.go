package label

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PSQL struct {
	pool *pgxpool.Pool
}

func NewPSQL(pool *pgxpool.Pool) (*PSQL, error) {
	return &PSQL{pool: pool}, nil
}

func (repo *PSQL) StoreLabel(ctx context.Context, l *Label) error {
	sql := "INSERT INTO labels (name, comment) VALUES ($1, $2) RETURNING id"
	row := repo.pool.QueryRow(ctx, sql, l.Name, l.Comment)
	if err := row.Scan(&l.ID); err != nil {
		return err
	}
	return nil
}

func (repo *PSQL) ListLabels(ctx context.Context) ([]Label, error) {
	sql := "SELECT id, name, comment FROM labels"
	rows, err := repo.pool.Query(ctx, sql)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	labels := []Label{}
	for rows.Next() {
		l := Label{}
		if err := rows.Scan(&l.ID, &l.Name, &l.Comment); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (repo *PSQL) GetLabel(ctx context.Context, id int64) (Label, error) {
	sql := "SELECT id, name, comment FROM labels WHERE id = $1"
	row := repo.pool.QueryRow(ctx, sql, id)
	l := Label{}
	if err := row.Scan(&l.ID, &l.Name, &l.Comment); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Label{}, ErrNotFound
		}
		return Label{}, err
	}

	tmplts, err := repo.ListTemplates(ctx, id)
	if err != nil {
		return Label{}, err
	}
	l.templates = make(map[string]Template, len(tmplts))
	for _, tmplt := range tmplts {
		l.templates[tmplt.Type] = tmplt
	}

	return l, nil
}

func (repo *PSQL) DeleteLabel(ctx context.Context, id int64) error {
	sql := "DELETE FROM labels WHERE id = $1"
	_, err := repo.pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PSQL) StoreTemplate(ctx context.Context, t *Template) error {
	sql := "INSERT INTO templates (label_id, type, body) VALUES ($1, $2, $3) RETURNING id"
	row := repo.pool.QueryRow(ctx, sql, t.LabelID, t.Type, t.Body)
	if err := row.Scan(&t.ID); err != nil {
		return err
	}
	return nil
}

func (repo *PSQL) ListTemplates(ctx context.Context, labelID int64) ([]Template, error) {
	sql := "SELECT id, label_id, type, body FROM templates WHERE label_id = $1"
	rows, err := repo.pool.Query(ctx, sql, labelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	templates := []Template{}
	for rows.Next() {
		t := Template{}
		if err := rows.Scan(&t.ID, &t.LabelID, &t.Type, &t.Body); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (repo *PSQL) GetTemplate(ctx context.Context, labelID, templateID int64) (Template, error) {
	sql := "SELECT id, label_id, type, body FROM templates WHERE id = $1 AND label_id = $2"
	row := repo.pool.QueryRow(ctx, sql, templateID, labelID)
	t := Template{}
	if err := row.Scan(&t.ID, &t.LabelID, &t.Type, &t.Body); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Template{}, ErrNotFound
		}
		return Template{}, err
	}

	return t, nil
}

func (repo *PSQL) DeleteTemplate(ctx context.Context, labelID, templateID int64) error {
	sql := "DELETE FROM templates WHERE id = $1 AND label_id = $2"
	_, err := repo.pool.Exec(ctx, sql, templateID, labelID)
	if err != nil {
		return err
	}
	return nil
}
