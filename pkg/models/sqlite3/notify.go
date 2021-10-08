package sqlite3

import (
	"database/sql"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
)

type NotifyModel struct {
	DB *sql.DB
}

func (m *NotifyModel) CreateNotifyTable() error {
	notifyTable, err := m.DB.Prepare(CreateNotifyTableSQL)
	if err != nil {
		return err
	}
	_, err = notifyTable.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *NotifyModel) InsertNotify(notify models.Notify) error {
	if notify.AuthorId == notify.UserId {
		return nil
	}
	insertNotify, err := m.DB.Prepare(InsertNotifySQL)
	if err != nil {
		return err
	}
	_, err = insertNotify.Exec(
		notify.AuthorId,
		notify.PostId,
		notify.UserId,
		notify.ActionType,
		notify.IsActive,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *NotifyModel) Get(authorId string) ([]models.Notify, error) {
	var currentNotify models.Notify
	var notifies []models.Notify

	rows, err := m.DB.Query("SELECT * FROM notify WHERE author_id = ?", authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&currentNotify.Id, &currentNotify.AuthorId, &currentNotify.PostId, &currentNotify.UserId, &currentNotify.ActionType, &currentNotify.IsActive)
		if err != nil {
			return nil, err
		}
		switch currentNotify.ActionType {
		case "comment":
			currentNotify.IsCommented = true
		case "like":
			currentNotify.IsLiked = true
		case "dislike":
			currentNotify.IsDisliked = true
		default:
		}
		notifies = append(notifies, currentNotify)
	}
	return notifies, nil
}

func (m *NotifyModel) UpdateNotifyStatus(authorId string) error {
	_, err := m.DB.Exec(UpdateNotifyStatusSQL, authorId)
	if err != nil {
		return err
	}
	return nil
}