package sqlite3

import "database/sql"

type RatingModel struct {
	DB *sql.DB
}

func (m *RatingModel) CreateRatingsTable() error {
	ratingPostsTable, err := m.DB.Prepare(CreateRatingPostSQL)
	if err != nil {
		return err
	}
	_, err = ratingPostsTable.Exec()
	if err != nil {
		return err
	}

	ratingCommentsTable, err := m.DB.Prepare(CreateRatingCommentSQL)
	if err != nil {
		return err
	}
	_, err = ratingCommentsTable.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) InsertPostRating(userId, postId, value int) error {
	var currentRatingValue int
	found, currentRatingValue, err := m.IsRatingExists(userId, postId, "post")
	if err != nil {
		return err
	}
	if !found {
		stmt, err := m.DB.Prepare(InsertRatingPostSQL)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(postId, userId, value)
		if err != nil {
			return err
		}
		return nil
	}
	if value != currentRatingValue {
		currentRatingValue += value
	}
	err = m.UpdatePostRating(userId, postId, currentRatingValue)
	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) InsertCommentRating(userId, commentId, value int) error {
	var currentRatingValue int
	found, currentRatingValue, err := m.IsRatingExists(userId, commentId, "comment")
	if err != nil {
		return err
	}
	if !found {
		stmt, err := m.DB.Prepare(InsertRatingCommentSQL)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(commentId, userId, value)
		if err != nil {
			return err
		}
		return nil
	}
	if value != currentRatingValue {
		currentRatingValue += value
	}
	err = m.UpdateCommentRating(userId, commentId, currentRatingValue)
	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) UpdateCommentRating(userId, commentId, value int) error {
	_, err := m.DB.Exec(UpdateRatingCommentSQL, value, userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) UpdatePostRating(userId, postId, value int) error {
	_, err := m.DB.Exec(UpdateRatingPostSQL, value, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func (m *RatingModel) IsRatingExists(userId, id int, flag string) (bool, int, error) {
	var value int
	switch flag {
	case "post":
		rows, err := m.DB.Query(SelectPostRatingByID, userId, id)
		if err != nil {
			return false, 0, err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&value)
			return true, value, err
		}
	case "comment":
		rows, err := m.DB.Query(SelectCommentRatingByID, userId, id)
		if err != nil {
			return false, 0, err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&value)
			return true, value, err
		}
	default:
		return false, 0, nil
	}
	return false, 0, nil
}

func (m *RatingModel) GetRatingById(id int, flag string) (int, error) {
	var value int
	var result int
	switch flag {
	case "post":
		rows, err := m.DB.Query("SELECT value FROM ratingPosts WHERE post_id = ?", id)
		if err != nil {
			return 0, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&value)
			if err != nil {
				return 0, err
			}
			result += value
		}
		return result, nil
	case "comment":
		rows, err := m.DB.Query("SELECT value FROM ratingComments WHERE comment_id = ?", id)
		if err != nil {
			return 0, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&value)
			if err != nil {
				return 0, err
			}
			result += value
		}
		return result, nil
	default:
		return 0, nil
	}
	return value, nil
}
