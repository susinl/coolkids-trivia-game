package mysql

import (
	"database/sql"

	"github.com/susinl/coolkids-trivia-game/application/port"
	"github.com/susinl/coolkids-trivia-game/domain/model"
)

type GameCodeRepository struct {
	db *sql.DB
}

func NewGameCodeRepository(db *sql.DB) port.GameCodeRepository {
	return &GameCodeRepository{
		db: db,
	}
}

func (r *GameCodeRepository) FindByCode(code string) (*model.GameCode, error) {
	row := r.db.QueryRow("SELECT game_code, name, email, phone_number, question_id, status, registered_time FROM game_codes WHERE game_code = ?", code)
	gameCode := &model.GameCode{}

	if err := row.Scan(&gameCode.Code, &gameCode.Name, &gameCode.Email, &gameCode.PhoneNumber, &gameCode.QuestionID, &gameCode.Status, &gameCode.RegisteredTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, port.ErrGameCodeNotFound
		}
		return nil, err
	}

	return gameCode, nil
}

func (r *GameCodeRepository) Save(gameCode *model.GameCode) error {
	_, err := r.db.Exec("INSERT INTO game_codes (game_code, name, email, phone_number, question_id, status, registered_time) VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=?, email=?, phone_number=?, question_id=?, status=?, registered_time=?", gameCode.Code, gameCode.Name, gameCode.Email, gameCode.PhoneNumber, gameCode.QuestionID, gameCode.Status, gameCode.RegisteredTime, gameCode.Name, gameCode.Email, gameCode.PhoneNumber, gameCode.QuestionID, gameCode.Status, gameCode.RegisteredTime)
	return err
}

func (r *GameCodeRepository) CountWinners() (int, error) {
	var count int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM game_codes WHERE status = 'won'").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
