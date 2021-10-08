package main

import (
	"crypto/tls"
	"flag"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	logrus *logrus.Logger
	notify   interface {
		CreateNotifyTable() error
		InsertNotify(models.Notify) error
		Get(string) ([]models.Notify, error)
		UpdateNotifyStatus(string) error
	}
	posts interface {
		CreatePostsTable() error
		InsertPostIntoDB(models.Post) (int, error)
		Get() ([]models.Post, error)
		GetPostById(string) (models.Post, error)
		GetPostsByCategory(string) ([]models.Post, error)
		GetPostsByAuthor(int) ([]models.Post, error)
		GetFavoritePosts(int) ([]models.Post, error)
		GetPostsByComments(int) ([]models.Post, error)
		DeletePost(int) error
		UpdatePost(models.Post) error
	}
	categoryPostLinks interface {
		CreateCategoryPostLinksTable() error
		InsertCategoryPostLinkIntoDB(int, []string) error
		Get(int) ([]string, error)
		DeleteLinks(int) error
	}
	comments interface {
		CreateCommentsTable() error
		InsertCommentIntoDB(models.Comment) error
		DeleteComment(int) error
		UpdateComment(int, string) error
		GetCommentById(int) (models.Comment, error)
	}
	ratings interface {
		CreateRatingsTable() error
		InsertPostRating(int, int, int) error
		InsertCommentRating(int, int, int) error
		GetRatingById(int, string) (int, error)
		GetCommentsByPostId(int) ([]models.Comment, error)
		IsRatingExists(int, int, string) (bool, int, error)
	}
	users interface {
		GetUserIdByLogin(string) (int, error)
		Get(int) (models.User, error)
		CreateUsersTable() error
		CreateUser(models.User) error
		Authenticate(string, []byte) (int, error)
		GetUsernameById(string) (string, error)
	}
}

func main() {
	addr := flag.String("addr", ":9000", "HTTP network address")
	dsn := flag.String("dsn", "./forum.db", "Sqlite3 data source name")
	flag.Parse()

	*addr = os.Getenv("PORT")

	infoLog := log.New(os.Stdout, Green+"INFO\t"+Reset, log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, Red+"ERROR\t"+Reset, log.Ldate|log.Ltime|log.Lshortfile)
	var log = logrus.New()

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(5)

	err = db.Ping()
	if err != nil {
		errorLog.Println(err)
	}
	defer db.Close()

	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		logrus: 		   log,
		notify:            &sqlite3.NotifyModel{DB: db},
		posts:             &sqlite3.PostModel{DB: db},
		categoryPostLinks: &sqlite3.CategoryPostLinkModel{DB: db},
		comments:          &sqlite3.CommentModel{DB: db},
		ratings:           &sqlite3.RatingModel{DB: db},
		users:             &sqlite3.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:           ":" + *addr,
		MaxHeaderBytes: 524288,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Infof("Server run on http://127.0.0.1:%s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
