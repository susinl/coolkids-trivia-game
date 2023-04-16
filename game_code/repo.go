package gameCode

import (
	"context"
	"database/sql"
	"fmt"
)

type QueryValidateGameCodeFn func(ctx context.Context, code string) (int, error)

func NewQueryParticipantByCodeFn(db *sql.DB) QueryValidateGameCodeFn {
	return func(ctx context.Context, code string) (int, error) {
		var count int
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM db.participant
			WHERE game_code = ?
		;`, code).Scan(&count)
		if err != nil {
			return 0, err
		}
		fmt.Println(count)
		return count, nil
	}
}

type QueryCheckStatusFn func(ctx context.Context, code string) (*ParticipantAnswer, error)

func NewQueryCheckStatusFn(db *sql.DB) QueryCheckStatusFn {
	return func(ctx context.Context, code string) (*ParticipantAnswer, error) {
		var participantAnswer ParticipantAnswer
		err := db.QueryRowContext(ctx, `
			SELECT	p.game_code, 
					p.question_id,
					p.name,
					p.email,
					p.answer,
					q.correct_answer
			FROM db.participant p 
			LEFT JOIN 	db.question q
						ON p.question_id = q.id 
			WHERE p.game_code = ?
		;`, code).Scan(
			&participantAnswer.GameCode,
			&participantAnswer.QuestionId,
			&participantAnswer.Name,
			&participantAnswer.Email,
			&participantAnswer.Answer,
			&participantAnswer.CorrectAnswer,
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
