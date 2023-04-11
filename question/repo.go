package question

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/susinl/coolkids-trivia-game/util"
)

type QueryParticipantByCodeFn func(ctx context.Context, code string) (*Participant, error)

func NewQueryParticipantByCodeFn(db *sql.DB) QueryParticipantByCodeFn {
	return func(ctx context.Context, code string) (*Participant, error) {
		var participant Participant
		err := db.QueryRowContext(ctx, `
			SELECT	game_code,
					name,
					email,
					phone_number,
					question_id,
					answer,
					registered_time
			FROM db.participant
			WHERE game_code = ?
		;`, code).Scan(
			&participant.GameCode,
			&participant.Name,
			&participant.Email,
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
					status
			FROM db.question
			WHERE status = ?
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

type UpdateQuestionStatusAndParticipantInfoFn func(ctx context.Context, code string, name string, email string, phone string, id int) error

func NewUpdateQuestionStatusAndParticipantInfoFn(db *sql.DB) UpdateQuestionStatusAndParticipantInfoFn {
	return func(ctx context.Context, code string, name string, email string, phone string, id int) error {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}

		resultQ, err := tx.ExecContext(ctx, `
			UPDATE db.question
			SET	status = 'pending'
			WHERE id = ?
		;`, id)
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
			SET	name = ?,
				email = ?,
				phone_number = ?,
				question_id = ?,
				registered_time = ?
			WHERE game_code = ?
		;`, name, email, phone, id, time.Now().Format(util.DateTimeFormat), code)
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
			WHERE x.game_code = ?
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
			SET status = ?
			WHERE id = ?
		;`, status, id)
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
			WHERE game_code = ?
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