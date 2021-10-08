package sqlite3

import (
	"database/sql"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) GetCommentById(commentId int) (models.Comment, error) {
	var comment models.Comment
	rows, err := m.DB.Query("SELECT * FROM comments WHERE id = ?", commentId)
	if err != nil {
		return comment, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return comment, err
		}
	}
	return comment, err
}

func (m *CommentModel) DeleteComment(commentId int) error {
	_, err := m.DB.Exec(DeleteCommentSQL, commentId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) UpdateComment(commentId int, content string) error {
	_, err := m.DB.Exec(UpdateCommentSQL, content, commentId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) CreateCommentsTable() error {
	commentsTable, err := m.DB.Prepare(CreateCommentsTableSQL)
	if err != nil {
		return err
	}
	_, err = commentsTable.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) InsertCommentIntoDB(commentData models.Comment) error {
	insertComment, err := m.DB.Prepare(InsertCommentSQL)
	if err != nil {
		return err
	}
	_, err = insertComment.Exec(
		commentData.PostId,
		commentData.UserId,
		commentData.Content,
		commentData.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) GetCommentsByPostId(postId int) ([]models.Comment, error) {
	var comment models.Comment
	var comments []models.Comment
	rows, err := m.DB.Query("SELECT * FROM comments WHERE post_id = ?", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.Rating, err = m.GetRatingById(comment.Id, "comment")
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, err
}