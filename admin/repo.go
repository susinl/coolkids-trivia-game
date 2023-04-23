package admin

import (
	"context"
	"database/sql"
)

type UpdateQuestionStatusFn func(ctx context.Context) (int64, error)

func NewUpdateQuestionStatusFn(db *sql.DB) UpdateQuestionStatusFn {
	return func(ctx context.Context) (int64, error) {
		result, err := db.ExecContext(ctx, `
			UPDATE	db.question
			SET status = 'ready'
			WHERE	DATE_ADD(recent_time, INTERVAL 20 SECOND) < NOW()
			AND		status = 'pending'	
		;`)
		if err != nil {
			return 0, err
		}
		row, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		return row, nil
	}
}
