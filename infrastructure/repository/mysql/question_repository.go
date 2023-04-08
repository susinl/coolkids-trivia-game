package mysql

import (
	"database/sql"

	"github.com/susinl/coolkids-trivia-game/application/port"
	"github.com/susinl/coolkids-trivia-game/domain/model"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) port.QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) FindAvailable() (*model.Question, error) {
	row := r.db.QueryRow("SELECT id, question_text, choice_A, choice_B, choice_C, choice_D, choice_E, choice_F, correct_answer, status FROM questions WHERE status = 'ready' LIMIT 1")
	question := &model.Question{}

	if err := row.Scan(&question.ID, &question.QuestionText, &question.Choices[0], &question.Choices[1], &question.Choices[2], &question.Choices[3], &question.Choices[4], &question.Choices[5], &question.CorrectAnswer, &question.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, port.ErrNoAvailableQuestions
		}
		return nil, err
	}

	return question, nil
}

func (r *QuestionRepository) Save(question *model.Question) error {
	_, err := r.db.Exec("UPDATE questions SET status=? WHERE id=?", question.Status, question.ID)
	return err
}
