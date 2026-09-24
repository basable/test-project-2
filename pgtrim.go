package main

import "context"

// trimLockKey is the Postgres advisory lock key that serializes trims.
const trimLockKey = 727400001

// Trim deletes the oldest n todos when more than max exist. The advisory
// transaction lock makes concurrent trims run one after another, so two
// simultaneous creates at the threshold cannot both delete n rows.
func (s *PGStore) Trim(ctx context.Context, max, n int) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, trimLockKey); err != nil {
		return 0, err
	}
	var count int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM todos`).Scan(&count); err != nil {
		return 0, err
	}
	if count <= max {
		return 0, tx.Commit(ctx)
	}
	tag, err := tx.Exec(ctx,
		`DELETE FROM todos WHERE id IN (
			SELECT id FROM todos ORDER BY created_at, id LIMIT $1)`, n)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), tx.Commit(ctx)
}
