package winners

import (
	"context"
	"database/sql"
)

type QueryWinnerListFn func(ctx context.Context) (*Winners, error)

func NewQueryWinnerListFn(db *sql.DB) QueryWinnerListFn {
	return func(ctx context.Context) (*Winners, error) {
		rows, err := db.QueryContext(ctx, `
			SELECT	p.name,
					p.email,
					p.phone_number,
					p.game_code,
					p.registered_time
			FROM db.participant p
			LEFT JOIN db.question q ON p.question_id = q.id
			WHERE p.answer = q.correct_answer
			ORDER BY p.registered_time ASC
			LIMIT 10
		;`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		winners := &Winners{Winners: []*Winner{}}

		for rows.Next() {
			winner := &Winner{}
			if err := rows.Scan(
				&winner.FullName,
				&winner.Email,
				&winner.PhoneNumber,
				&winner.Code,
				&winner.Timestamp,
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
