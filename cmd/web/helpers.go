package main

import (
	"fmt"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
	"net/http"
	"runtime/debug"
	"strconv"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.logrus.Errorln(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}
	_, ok := IsTokenExists(cookie.Value)
	if !ok {
		return false
	}
	return true
}

func (app *application) GetCommentsStruct(postId int) ([]models.Comment, error) {
	var comments []models.Comment
	comments, err := app.ratings.GetCommentsByPostId(postId)
	if err != nil {
		return nil, err
	}
	for i := range comments {
		comments[i].Username, err = app.users.GetUsernameById(comments[i].UserId)
		if err != nil {
			return nil, err
		}
		comments[i].Rating, err = app.ratings.GetRatingById(comments[i].Id, "comment")
		if err != nil {
			return nil, err
		}
	}
	return comments, nil
}

func (app *application) GetPostsStruct(id, category string, userId int) ([]models.Post, error) {
	var posts []models.Post
	var err error
	switch id {
	case "comments":
		posts, err = app.posts.GetPostsByComments(userId)
		if err != nil {
			return nil, err
		}
	case "favorite":
		posts, err = app.posts.GetFavoritePosts(userId)
		if err != nil {
			return nil, err
		}
	case "category":
		posts, err = app.posts.GetPostsByCategory(category)
		if err != nil {
			return nil, err
		}
	case "my":
		posts, err = app.posts.GetPostsByAuthor(userId)
		if err != nil {
			return nil, err
		}
	case "":
		posts, err = app.posts.Get()
		if err != nil {
			return nil, err
		}
	default:
		post, err := app.posts.GetPostById(id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	for i, _ := range posts {
		posts[i].Comments, err = app.GetCommentsStruct(posts[i].Id)
		if err != nil {
			return nil, err
		}
		posts[i].Rating, err = app.ratings.GetRatingById(posts[i].Id, "post")
		if err != nil {
			return nil, err
		}
		posts[i].Category, err = app.categoryPostLinks.Get(posts[i].Id)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (app *application) GetAuthorIdByPostId(postId int) (string, error) {
	post, err  := app.posts.GetPostById(strconv.Itoa(postId))
	if err != nil {
		return "", err
	}
	return post.UserId, nil
}

func (app *application) GetUserId(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, err
	}
	userId, err := app.GetUserIdByCookie(cookie)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (app *application) GetRatingMarkupPost(userId int, resp models.Resp) ([]models.Post, error) {
	if userId == 0 {
		return resp.Posts, nil
	}
	for i := range resp.Posts {
		ok, res, err := app.ratings.IsRatingExists(userId, resp.Posts[i].Id, "post")
		if err != nil {
			return nil, err
		}
		if ok {
			if res == 1 {
				resp.Posts[i].IsLiked = true
			}
			if res == -1 {
				resp.Posts[i].IsDisliked = true
			}
		}
	}
	return resp.Posts, nil
}

func (app *application) GetRatingMarkupComment(userId int, resp models.Resp) ([]models.Post, error) {
	if userId == 0 {
		return resp.Posts, nil
	}
	for i := range resp.Posts {
		for j := range resp.Posts[i].Comments {
			ok, res, err := app.ratings.IsRatingExists(userId, resp.Posts[i].Comments[j].Id, "comment")
			if err != nil {
				return nil, err
			}
			if ok {
				if res == 1 {
					resp.Posts[i].Comments[j].IsLiked = true
				}
				if res == -1 {
					resp.Posts[i].Comments[j].IsDisliked = true
				}
			}
		}
	}
	return resp.Posts, nil
}

func (app *application) GetNotifyStruct(userId string) ([]models.Notify, error) {
	notifies, err := app.notify.Get(userId)
	if err != nil {
		return nil, err
	}
	for i := range notifies {
		notifies[i].UserLogin, err = app.users.GetUsernameById(notifies[i].UserId)
		if err != nil {
			return nil, err
		}
	}
	return notifies, nil
}

func (app *application) NewNotifyCounter(notifies []models.Notify) int {
	var counter int
	for _, notify := range notifies {
		if notify.IsActive {
			counter++
		}
	}
	return counter
}