package sqlite3

/*------------------------------------------------------*/
/*                                                      */
/*                    NOTIFY STATEMENTS                 */
/*                                                      */
/*------------------------------------------------------*/

const CreateNotifyTableSQL = `
	CREATE TABLE IF NOT EXISTS notify (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		author_id TEXT NOT NULL,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		action_type TEXT NOT NULL,
		is_active BOOLEAN NOT NULL
	);
`

const InsertNotifySQL = `
	INSERT INTO notify (
		author_id, post_id, user_id, action_type, is_active
	) VALUES (?, ?, ?, ?, ?);
`

const UpdateNotifyStatusSQL = `
	UPDATE notify SET is_active = false
	WHERE author_id = ?;
`

/*------------------------------------------------------*/
/*                                                      */
/*                    POST STATEMENTS                   */
/*                                                      */
/*------------------------------------------------------*/

const CreatePostsTableSQL = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		human_date TEXT NOT NULL,
		image_url TEXT
	);
`
const InsertPostSQL = `
	INSERT INTO posts (
		title, content, user_id, created_at, human_date, image_url
	) VALUES (?, ?, ?, ?, ?, ?);
`

const DeletePostSQL = `
	DELETE FROM posts WHERE id = ?;
`

const UpdatePostSQL = `
	UPDATE posts SET title = ?, content = ? WHERE id = ?;
`

/*------------------------------------------------------*/
/*                                                      */
/*                    USER STATEMENTS                   */
/*                                                      */
/*------------------------------------------------------*/

const CreateUsersTableSQL = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		login TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TEXT NOT NULL
	);
`
const InsertUserSQL = `
	INSERT INTO users (
		login, email, password, created_at
	) VALUES (?, ?, ?, ?);
`

/*------------------------------------------------------*/
/*                                                      */
/*                   COMMENT STATEMENTS                 */
/*                                                      */
/*------------------------------------------------------*/

const CreateCommentsTableSQL = `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
`

const InsertCommentSQL = `
	INSERT INTO comments (
		post_id, user_id, content, created_at
	) VALUES (?, ?, ?, ?);
`

const DeleteCommentSQL = `
	DELETE FROM comments WHERE id = ?;
`

const UpdateCommentSQL = `
	UPDATE comments SET content = ? WHERE id = ?;
`

/*------------------------------------------------------*/
/*                                                      */
/*                   CATEGORY STATEMENTS                */
/*                                                      */
/*------------------------------------------------------*/

const CreateCategoryTableSQL = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
`

const CreateCategoryPostLinkSQL = `
	CREATE TABLE IF NOT EXISTS categoryPostLink (
		post_id INTEGER NOT NULL,
		category_name STRING NOT NULL
	);
`

const DeleteLinkSQL = `
	DELETE FROM categoryPostLink WHERE post_id = ?;
`

const InsertCategoriesSQL = `
	INSERT INTO categories (
		name
	) VALUES (?);
`

const InsertCategoryPostLinkSQL = `
	INSERT INTO categoryPostLink (
		post_id, category_name
	) VALUES (?, ?);
`

/*------------------------------------------------------*/
/*                                                      */
/*                    RATING STATEMENTS                 */
/*                                                      */
/*------------------------------------------------------*/

const CreateRatingPostSQL = `
	CREATE TABLE IF NOT EXISTS ratingPosts (
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		value INTEGER NOT NULL
	);
`

const CreateRatingCommentSQL = `
	CREATE TABLE IF NOT EXISTS ratingComments (
		comment_id INTEGER NOT NULL,
		user_Id INTEGER NOT NULL,
		value INTEGER NOT NULL
	);
`

const InsertRatingPostSQL = `
	INSERT INTO ratingPosts (
		post_id, user_id, value
	) VALUES (?, ?, ?);
`

const InsertRatingCommentSQL = `
	INSERT INTO ratingComments (
		comment_id, user_id, value
	) VALUES (?, ?, ?);
`

const UpdateRatingPostSQL = `
	UPDATE ratingPosts SET value = ?
	WHERE user_id = ? AND post_id = ?;
`

const UpdateRatingCommentSQL = `
	UPDATE ratingComments SET value = ?
	WHERE user_id = ? AND comment_id = ?;
`

const SelectPostRatingByID = `
	SELECT value FROM ratingPosts where user_id = ? AND post_id = ?;
`

const SelectCommentRatingByID = `
	SELECT value FROM ratingComments where user_id = ? AND comment_id = ?;
`
