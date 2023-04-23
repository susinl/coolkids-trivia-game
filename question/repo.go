package question

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/susinl/coolkids-trivia-game/util"
)

type QueryCountTotalWinnerFn func(ctx context.Context) (int, error)

func NewQueryCountTotalWinnerFn(db *sql.DB) QueryCountTotalWinnerFn {
	return func(ctx context.Context) (int, error) {
		var count int
		err := db.QueryRowContext(ctx, `
			SELECT	COUNT(*) AS total_winner
			FROM db.participant x
			LEFT JOIN db.question y ON x.question_id = y.id
			WHERE x.answer = y.correct_answer
		;`).Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}

type QueryParticipantByCodeFn func(ctx context.Context, code string) (*Participant, error)

func NewQueryParticipantByCodeFn(db *sql.DB) QueryParticipantByCodeFn {
	return func(ctx context.Context, code string) (*Participant, error) {
		var participant Participant
		err := db.QueryRowContext(ctx, `
			SELECT	code,
					name,
					phone_number,
					question_id,
					answer,
					registered_time
			FROM db.participant
			WHERE code = ?
		;`, code).Scan(
			&participant.Code,
			&participant.Name,
			&participant.PhoneNumber,
			&participant.QuestionId,
			&participant.Answer,
			&participant.RegisteredTime,
		)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return &participant, nil
		}
	}
}

type QueryQuestionByStatusFn func(ctx context.Context, status string) (*Question, error)

func NewQueryQuestionByStatusFn(db *sql.DB) QueryQuestionByStatusFn {
	return func(ctx context.Context, status string) (*Question, error) {
		var question Question
		err := db.QueryRowContext(ctx, `
			SELECT	id,
					question_text,
					choice_a,
					choice_b,
					choice_c,
					choice_d,
					choice_e,
					choice_f,
					correct_answer,
					status,
					recent_time
			FROM db.question
			WHERE status = ?
			ORDER BY recent_time ASC
		;`, status).Scan(
			&question.Id,
			&question.QuestionText,
			&question.ChoiceA,
			&question.ChoiceB,
			&question.ChoiceC,
			&question.ChoiceD,
			&question.ChoiceE,
			&question.ChoiceF,
			&question.CorrectAnswer,
			&question.Status,
			&question.RecentTime,
		)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return &question, nil
		}
	}
}

type UpdateQuestionStatusAndParticipantInfoFn func(ctx context.Context, code string, name string, phone string, id int, byPass bool) error

func NewUpdateQuestionStatusAndParticipantInfoFn(db *sql.DB) UpdateQuestionStatusAndParticipantInfoFn {
	return func(ctx context.Context, code string, name string, phone string, id int, byPass bool) error {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}

		if !byPass {
			resultQ, err := tx.ExecContext(ctx, `
				UPDATE db.question
				SET	status = 'pending',
					recent_time = ?
				WHERE id = ?
			;`, time.Now().Format(util.DateTimeFormat), id)
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
		}

		resultP, err := tx.ExecContext(ctx, `
			UPDATE db.participant
			SET	name = ?,
				phone_number = ?,
				question_id = ?,
				registered_time = ?
			WHERE code = ?
		;`, name, phone, id, time.Now().Format(util.DateTimeFormat), code)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}
		rowP, err := resultP.RowsAffected()
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}
		if rowP != 1 {
			err := errors.Errorf("Expected to affect 1 row but got affected %d row", rowP)
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

type QueryParticipantAndAnswerFn func(ctx context.Context, code string) (*ParticipantWAnswer, error)

func NewQueryParticipantAndAnswerFn(db *sql.DB) QueryParticipantAndAnswerFn {
	return func(ctx context.Context, code string) (*ParticipantWAnswer, error) {
		var participantWAnswer ParticipantWAnswer
		err := db.QueryRowContext(ctx, `
			SELECT	x.question_id,
					x.answer,
					y.correct_answer,
					x.registered_time
			FROM db.participant x
			LEFT JOIN db.question y ON x.question_id = y.id
			WHERE x.code = ?
		;`, code).Scan(
			&participantWAnswer.QuestionId,
			&participantWAnswer.Answer,
			&participantWAnswer.CorrectAnswer,
			&participantWAnswer.RegisteredTime,
		)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return &participantWAnswer, nil
		}
	}
}

type UpdateParticipantAnswerAndStatusFn func(ctx context.Context, code string, answer int, id int, status string) error

func NewUpdateParticipantAnswerAndStatusFn(db *sql.DB) UpdateParticipantAnswerAndStatusFn {
	return func(ctx context.Context, code string, answer int, id int, status string) error {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}

		resultQ, err := tx.ExecContext(ctx, `
			UPDATE db.question
			SET status = ?,
				recent_time = ?
			WHERE id = ?
		;`, status, time.Now().Format(util.DateTimeFormat), id)
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

		resultP, err := tx.ExecContext(ctx, `
			UPDATE db.participant
			SET	answer = ?
			WHERE code = ?
		;`, answer, code)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}
		rowP, err := resultP.RowsAffected()
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}
			return err
		}
		if rowP != 1 {
			err := errors.Errorf("Expected to affect 1 row but got affected %d row", rowP)
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
