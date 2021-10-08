package sqlite3

import (
	"database/sql"
	"errors"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetUsernameById(id string) (string, error) {
	var login string
	stmt := "SELECT login FROM users WHERE id = ?"
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNoRecord
		} else {
			return "", err
		}
	}
	return login, nil
}

func (m *UserModel) GetUserIdByLogin(username string) (int, error) {
	var id int
	stmt := "SELECT id FROM users WHERE email = ?"
	row := m.DB.QueryRow(stmt, username)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) CreateUsersTable() error {
	usersTable, err := m.DB.Prepare(CreateUsersTableSQL)
	if err != nil {
		return err
	}
	_, err = usersTable.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) CreateUser(user models.User) error {
	insertStmt, err := m.DB.Prepare(InsertUserSQL)
	if err != nil {
		return err
	}
	_, err = insertStmt.Exec(
		user.Login,
		user.Email,
		user.Password,
		time.Now().Format("2006-January-02"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email string, password []byte) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, password FROM users WHERE email = ?"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Get(userId int) (models.User, error) {
	var user models.User
	stmt := "SELECT login, email, created_at FROM users WHERE id = ?"
	row := m.DB.QueryRow(stmt, userId)
	err := row.Scan(&user.Login, &user.Email, &user.CreatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
