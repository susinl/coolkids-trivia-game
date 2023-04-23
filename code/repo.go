package code

import (
	"context"
	"database/sql"
)

type QueryValidateCodeFn func(ctx context.Context, code string) (int, error)

func NewQueryParticipantByCodeFn(db *sql.DB) QueryValidateCodeFn {
	return func(ctx context.Context, code string) (int, error) {
		var count int
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM db.participant
			WHERE BINARY code = ?
		;`, code).Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}

type QueryCheckStatusFn func(ctx context.Context, code string) (*ParticipantAnswer, error)

func NewQueryCheckStatusFn(db *sql.DB) QueryCheckStatusFn {
	return func(ctx context.Context, code string) (*ParticipantAnswer, error) {
		var participantAnswer ParticipantAnswer
		err := db.QueryRowContext(ctx, `
			SELECT	p.code, 
					p.question_id,
					p.name,
					p.answer,
					q.correct_answer,
					p.phone_number
			FROM db.participant p 
			LEFT JOIN 	db.question q
						ON p.question_id = q.id 
			WHERE p.code = ?
		;`, code).Scan(
			&participantAnswer.Code,
			&participantAnswer.QuestionId,
			&participantAnswer.Name,
			&participantAnswer.Answer,
			&participantAnswer.CorrectAnswer,
			&participantAnswer.PhoneNumber,
		)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return &participantAnswer, nil
		}

	}
}
