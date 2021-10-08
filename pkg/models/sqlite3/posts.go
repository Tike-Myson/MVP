package sqlite3

import (
	"database/sql"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) DeletePost(postId int) error {
	_, err := m.DB.Exec(DeletePostSQL, postId)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) UpdatePost(post models.Post) error {
	_, err := m.DB.Exec(UpdatePostSQL, post.Title, post.Content, post.Id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) CreatePostsTable() error {
	postsTable, err := m.DB.Prepare(CreatePostsTableSQL)
	if err != nil {
		return err
	}
	_, err = postsTable.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) InsertPostIntoDB(postData models.Post) (int, error) {
	insertPost, err := m.DB.Prepare(InsertPostSQL)
	if err != nil {
		return 0, err
	}
	res, err := insertPost.Exec(
		postData.Title,
		postData.Content,
		postData.UserId,
		postData.CreatedAt,
		postData.HumanDate,
		postData.ImageURL,
	)
	if err != nil {
		return 0, err
	}
	lid, err := res.LastInsertId()
	return int(lid), nil
}

func (m *PostModel) Get() ([]models.Post, error) {
	var CurrentPost models.Post
	var Posts []models.Post

	rows, err := m.DB.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&CurrentPost.Id, &CurrentPost.Title, &CurrentPost.Content, &CurrentPost.UserId, &CurrentPost.CreatedAt, &CurrentPost.HumanDate, &CurrentPost.ImageURL)
		if err != nil {
			return nil, err
		}
		Posts = append(Posts, CurrentPost)
	}
	return Posts, nil
}

func (m *PostModel) GetPostsByCategory(categoryName string) ([]models.Post, error) {
	var postsArr []string
	var postId string
	var post models.Post
	var posts []models.Post

	rows, err := m.DB.Query("SELECT post_id FROM categoryPostLink WHERE category_name = ?", categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&postId)
		if err != nil {
			return nil, err
		}
		postsArr = append(postsArr, postId)
	}

	for _, v := range postsArr {
		rows, err := m.DB.Query("SELECT * FROM posts WHERE id = ?", v)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.HumanDate, &post.ImageURL)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (m *PostModel) GetPostsByAuthor(userId int) ([]models.Post, error) {
	var post models.Post
	var posts []models.Post

	rows, err := m.DB.Query("SELECT * FROM posts WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.HumanDate, &post.ImageURL)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *PostModel) GetFavoritePosts(userId int) ([]models.Post, error) {
	var postsArr []string
	var postId string
	var post models.Post
	var posts []models.Post
	rows, err := m.DB.Query("SELECT post_id FROM ratingPosts WHERE value = 1 AND user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&postId)
		if err != nil {
			return nil, err
		}
		postsArr = append(postsArr, postId)
	}

	for _, v := range postsArr {
		rows, err := m.DB.Query("SELECT * FROM posts WHERE id = ?", v)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.HumanDate, &post.ImageURL)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (m *PostModel) GetPostById(id string) (models.Post, error) {
	var post models.Post

	rows, err := m.DB.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return post, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.HumanDate, &post.ImageURL)
		if err != nil {
			return post, err
		}
	}
	if post.Id > 0 {
		return post, err
	}
	return post, models.ErrNoRecord
}

func (m *PostModel) GetPostsByComments(userId int) ([]models.Post, error) {
	var postId string
	var postIdArr []string
	var post models.Post
	var posts []models.Post
	rows, err := m.DB.Query("SELECT post_id FROM comments WHERE user_id", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&postId)
		if err != nil {
			return nil, err
		}
		ok := IsPostIdExists(postId, postIdArr)
		if !ok {
			postIdArr = append(postIdArr, postId)
		}
	}
	for _, v := range postIdArr {
		rows, err := m.DB.Query("SELECT * FROM posts WHERE id = ?", v)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.HumanDate, &post.ImageURL)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func IsPostIdExists(postId string, postIdArr []string) bool {
	for i := range postIdArr {
		if postIdArr[i] == postId {
			return true
		}
	}
	return false
}
