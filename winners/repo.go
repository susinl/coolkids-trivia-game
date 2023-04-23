package winners

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type QueryWinnerListFn func(ctx context.Context) (*Winners, error)

func NewQueryWinnerListFn(db *sql.DB) QueryWinnerListFn {
	return func(ctx context.Context) (*Winners, error) {
		rows, err := db.QueryContext(ctx, `
			SELECT	p.name,
					p.phone_number,
					p.code,
					p.registered_time,
					q.recent_time
			FROM db.participant p
			LEFT JOIN db.question q ON p.question_id = q.id
			WHERE p.answer = q.correct_answer
			ORDER BY q.recent_time ASC
		;`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		winners := &Winners{Winners: []*Winner{}}

		for rows.Next() {
			winner := &Winner{}
			if err := rows.Scan(
				&winner.Name,
				&winner.PhoneNumber,
				&winner.Code,
				&winner.Timestamp,
				&winner.RecentTime,
			); err != nil {
				return nil, err
			}
			winners.Winners = append(winners.Winners, winner)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}

		return winners, nil
	}
}

type QueryGetQuotaFn func(ctx context.Context) (int, error)

func NewQueryGetQuotaFn(db *sql.DB) QueryGetQuotaFn {
	return func(ctx context.Context) (int, error) {
		var quota int
		err := db.QueryRowContext(ctx, `
			SELECT quota
			FROM db.game_quota
		;`).Scan(&quota)
		if err != nil {
			return 0, err
		}
		// fmt.Println(count)
		return quota, nil
	}
}

type UpdateQuotaFn func(ctx context.Context, newQuota int) error

func NewUpdateQuotaFn(db *sql.DB) UpdateQuotaFn {
	return func(ctx context.Context, newQuota int) error {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}

		// fmt.Println(newQuota)

		resultQ, err := tx.ExecContext(ctx, `
			UPDATE db.game_quota
			SET quota = ?
		;`, newQuota)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}

		rowQ, err := resultQ.RowsAffected()
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}

		if rowQ != 1 {
			err := errors.Errorf("Expected to affect 1 row but got affected %d row", rowQ)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}

		return nil
	}
}
