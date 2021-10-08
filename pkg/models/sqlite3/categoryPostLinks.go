package sqlite3

import "database/sql"

type CategoryPostLinkModel struct {
	DB *sql.DB
}

func (m *CategoryPostLinkModel) CreateCategoryPostLinksTable() error {
	categoryPostLinkTable, err := m.DB.Prepare(CreateCategoryPostLinkSQL)
	if err != nil {
		return err
	}
	_, err = categoryPostLinkTable.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *CategoryPostLinkModel) DeleteLinks(postId int) error {
	_, err := m.DB.Exec(DeleteLinkSQL, postId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CategoryPostLinkModel) InsertCategoryPostLinkIntoDB(postId int, categoryName []string) error {
	stmt, err := m.DB.Prepare(InsertCategoryPostLinkSQL)
	if err != nil {
		return err
	}
	for _, v := range categoryName {
		_, err = stmt.Exec(postId, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *CategoryPostLinkModel) Get(postId int) ([]string, error) {
	var id string
	var categoryName string
	var categories []string
	rows, err := m.DB.Query("SELECT * FROM categoryPostLink WHERE post_id = ?", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &categoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categoryName)
	}
	return categories, nil
}
