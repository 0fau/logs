// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: copyfrom.go

package sql

import (
	"context"
)

// iteratorForInsertPlayer implements pgx.CopyFromSource.
type iteratorForInsertPlayer struct {
	rows                 []InsertPlayerParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertPlayer) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertPlayer) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Encounter,
		r.rows[0].Boss,
		r.rows[0].Difficulty,
		r.rows[0].Class,
		r.rows[0].Name,
		r.rows[0].Dead,
		r.rows[0].GearScore,
		r.rows[0].Dps,
		r.rows[0].Place,
	}, nil
}

func (r iteratorForInsertPlayer) Err() error {
	return nil
}

func (q *Queries) InsertPlayer(ctx context.Context, arg []InsertPlayerParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"players"}, []string{"encounter", "boss", "difficulty", "class", "name", "dead", "gear_score", "dps", "place"}, &iteratorForInsertPlayer{rows: arg})
}
