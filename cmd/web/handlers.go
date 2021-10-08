package main

import (
	"errors"
	"git.01.alem.school/Nurtilek_Asankhan/forum-authentication/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) removeComment(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	switch r.Method {
	case "POST":
		isAuth := app.isAuthenticated(r)
		if !isAuth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/comment/remove/"):]
		commentId, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		comment, err := app.comments.GetCommentById(commentId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		if comment.UserId != strconv.Itoa(userId) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		err = app.comments.DeleteComment(commentId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) editComment(w http.ResponseWriter, r *http.Request) {
	var resp models.Resp
	t, err := template.ParseFiles("./ui/html/updateComment.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	switch r.Method {
	case "GET":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/commentEdit/"):]
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		commentId, err := strconv.Atoi(id)
		comment, err := app.comments.GetCommentById(commentId)
		if comment.UserId != strconv.Itoa(userId) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		resp.Comment = comment
		err = t.Execute(w, resp)
		if err != nil {
			return
		}
	case "POST":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusUnauthorized)
			return
		}
		id := r.URL.Path[len("/commentEdit/"):]
		commentId, err := strconv.Atoi(id)
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		comment, err := app.comments.GetCommentById(commentId)
		if comment.UserId != strconv.Itoa(userId) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		comment.Content = r.FormValue("content")
		err = app.comments.UpdateComment(comment.Id, comment.Content)
		if err != nil {
			app.serverError(w, err)
			return
		}
		url := "/post/" + comment.PostId
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) removePost(w http.ResponseWriter, r *http.Request) {
	//referer := r.Header.Get("Referer")
	switch r.Method {
	case "POST":
		isAuth := app.isAuthenticated(r)
		if !isAuth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/post/remove/"):]
		postId, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post, err := app.posts.GetPostById(id)
		if post.UserId != strconv.Itoa(userId) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		err = app.posts.DeletePost(postId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) editPost(w http.ResponseWriter, r *http.Request) {
	//referer := r.Header.Get("Referer")
	var resp models.Resp
	t, err := template.ParseFiles("./ui/html/updatePost.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	switch r.Method {
	case "GET":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/edit/"):]
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		post, err := app.posts.GetPostById(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		posts, err := app.GetPostsStruct(id, "", 0)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post = posts[0]
		post.UpdatedCategory = strings.Join(post.Category, " ")
		if post.UserId != strconv.Itoa(userId) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		resp.Posts = append(resp.Posts, post)
		err = t.Execute(w, resp)
		if err != nil {
			return
		}
	case "POST":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusUnauthorized)
			return
		}
		var post models.Post
		id := r.URL.Path[len("/edit/"):]
		post, err = app.posts.GetPostById(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")
		post.Category = strings.Fields(r.FormValue("tags"))
		err = app.posts.UpdatePost(post)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.categoryPostLinks.DeleteLinks(post.Id)
		err = app.categoryPostLinks.InsertCategoryPostLinkIntoDB(post.Id, post.Category)
		if err != nil {
			app.serverError(w, err)
			return
		}
		url := "/post/" + strconv.Itoa(post.Id)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) showNotify(w http.ResponseWriter, r *http.Request) {
	var resp models.Resp
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./ui/html/notify.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r , "/user/login", http.StatusSeeOther)
			return
		}
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.notify.UpdateNotifyStatus(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) postsByComments(w http.ResponseWriter, r *http.Request) {
	var resp models.Resp
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./ui/html/index.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetPostsStruct("comments", "", userId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupPost(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupComment(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

//func (app *application) googleAuth(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/auth/google" {
//		app.clientError(w, http.StatusNotFound)
//		return
//	}
//	switch r.Method {
//	case "GET":
//	default:
//		app.clientError(w, http.StatusMethodNotAllowed)
//		return
//	}
//}
//
//func (app *application) githubAuth(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/auth/github" {
//		app.clientError(w, http.StatusNotFound)
//		return
//	}
//}
//
//func (app *application) facebookAuth(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/auth/facebook" {
//		app.clientError(w, http.StatusNotFound)
//		return
//	}
//}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	if !firstEnterFlag {
		err := app.createAllTables()
		if err != nil {
			app.serverError(w, err)
			return
		}
		firstEnterFlag = true
	}
	if r.URL.Path != "/" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	var resp models.Resp

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./ui/html/index.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		resp.Posts, err = app.GetPostsStruct("", "", 0)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupPost(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) postByCategory(w http.ResponseWriter, r *http.Request) {
	var resp models.Resp
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./ui/html/index.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		categoryName := r.URL.Path[len("/category/"):]
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetPostsStruct("category", categoryName, userId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupPost(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) favoritePosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/favorite/" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	var resp models.Resp
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./ui/html/index.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetPostsStruct("favorite", "", userId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupPost(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) myPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/my/" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	var resp models.Resp
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./ui/html/index.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetPostsStruct("my", "", userId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupPost(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	var resp models.ErrResp
	resp.IsAuthenticated = app.isAuthenticated(r)
	switch r.Method {
	case "GET":
		if resp.IsAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		t, err := template.ParseFiles("./ui/html/login.html")
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	case "POST":
		if resp.IsAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		login := r.FormValue("inputLogin")
		pass := r.FormValue("inputPassword")
		_, err := app.users.Authenticate(login, []byte(pass))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				resp.IsInvalidCredentials = true
				t, err := template.ParseFiles("./ui/html/login.html")
				if err != nil {
					app.serverError(w, err)
					return
				}
				err = t.Execute(w, resp)
				if err != nil {
					app.serverError(w, err)
					return
				}
				return
			} else {
				app.serverError(w, err)
				return
			}
		}
		token, ok := IsSessionExists(login)
		if ok {
			cookie := DeleteCookie(token)
			http.SetCookie(w, &cookie)
		}
		cookie := MakeCookie(login)
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/register" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	var resp models.ErrResp
	t, err := template.ParseFiles("./ui/html/signup.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	switch r.Method {
	case "GET":
		if resp.IsAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	case "POST":
		if resp.IsAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		login := r.FormValue("inputLogin")
		email := r.FormValue("inputEmail")
		pass := r.FormValue("inputPassword")
		rPass := r.FormValue("repeatInputPassword")
		if pass != rPass {
			resp.IsDifferentPasswords = true
			err = t.Execute(w, resp)
			return
		}
		password, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
		if err != nil {
			app.serverError(w, err)
			return
		}
		user := models.User{
			Login:    login,
			Email:    email,
			Password: password,
		}
		err = app.users.CreateUser(user)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.email" {
				resp.IsDuplicateEmail = true
				t.Execute(w, resp)
				return
			}
			if err.Error() == "UNIQUE constraint failed: users.nickname" {
				resp.IsDuplicateUsername = true
				t.Execute(w, resp)
				return
			}
			app.serverError(w, err)
			return
		}
		cookie := MakeCookie(email)
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/logout" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		for _, cookie := range r.Cookies() {
			dCookie := DeleteCookie(cookie.Value)
			http.SetCookie(w, &dCookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	var resp models.Resp
	t, err := template.ParseFiles("./ui/html/createPost.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	switch r.Method {
	case "GET":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	case "POST":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusUnauthorized)
			return
		}
		var post models.Post
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")
		post.Category = strings.Fields(r.FormValue("tags"))
		post.CreatedAt = time.Now()
		post.HumanDate = time.Now().Format("January 2, 2006")
		post.UserId = strconv.Itoa(userId)
		postId, err := app.posts.InsertPostIntoDB(post)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.categoryPostLinks.InsertCategoryPostLinkIntoDB(postId, post.Category)
		if err != nil {
			app.serverError(w, err)
			return
		}
		url := "/post/" + strconv.Itoa(postId)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	var resp models.ErrResp
	switch r.Method {
	case "POST":
		resp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusUnauthorized)
			return
		}
		var comment models.Comment
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		comment.Content = r.FormValue("content")
		comment.CreatedAt = time.Now()
		comment.UserId = strconv.Itoa(userId)
		comment.PostId = r.FormValue("postId")
		err = app.comments.InsertCommentIntoDB(comment)
		if err != nil {
			app.serverError(w, err)
			return
		}
		var notify models.Notify
		postId, err := strconv.Atoi(comment.PostId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		notify.AuthorId, err = app.GetAuthorIdByPostId(postId)
		notify.UserId = strconv.Itoa(userId)
		notify.ActionType = "comment"
		notify.IsActive = true
		err = app.notify.InsertNotify(notify)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/post/"+comment.PostId, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) postById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/post/"):]
	var resp models.Resp
	t, err := template.ParseFiles("./ui/html/post.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	switch r.Method {
	case "GET":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		resp.Posts, err = app.GetPostsStruct(id, "", 0)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != http.ErrNoCookie && err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupPost(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Posts, err = app.GetRatingMarkupComment(userId, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) likePost(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	switch r.Method {
	case "POST":
		isAuth := app.isAuthenticated(r)
		if !isAuth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/post/like/"):]
		postId, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.ratings.InsertPostRating(userId, postId, 1)
		if err != nil {
			app.serverError(w, err)
			return
		}
		var notify models.Notify
		notify.AuthorId, err = app.GetAuthorIdByPostId(postId)
		notify.UserId = strconv.Itoa(userId)
		notify.PostId = strconv.Itoa(postId)
		notify.ActionType = "like"
		notify.IsActive = true
		err = app.notify.InsertNotify(notify)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) dislikePost(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	switch r.Method {
	case "POST":
		isAuth := app.isAuthenticated(r)
		if !isAuth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/post/dislike/"):]
		postId, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.ratings.InsertPostRating(userId, postId, -1)
		if err != nil {
			app.serverError(w, err)
			return
		}
		var notify models.Notify
		notify.AuthorId, err = app.GetAuthorIdByPostId(postId)
		notify.UserId = strconv.Itoa(userId)
		notify.ActionType = "dislike"
		notify.IsActive = true
		err = app.notify.InsertNotify(notify)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	switch r.Method {
	case "POST":
		isAuth := app.isAuthenticated(r)
		if !isAuth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/comment/like/"):]
		commentId, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.ratings.InsertCommentRating(userId, commentId, 1)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) dislikeComment(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	switch r.Method {
	case "POST":
		isAuth := app.isAuthenticated(r)
		if !isAuth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		id := r.URL.Path[len("/comment/dislike/"):]
		commentId, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.ratings.InsertCommentRating(userId, commentId, -1)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/profile" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	var resp models.Resp
	t, err := template.ParseFiles("./ui/html/profile.html")
	if err != nil {
		app.serverError(w, err)
		return
	}
	switch r.Method {
	case "GET":
		resp.ErrResp.IsAuthenticated = app.isAuthenticated(r)
		if !resp.ErrResp.IsAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		userId, err := app.GetUserId(r)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.User, err = app.users.Get(userId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.Notify, err = app.GetNotifyStruct(strconv.Itoa(userId))
		if err != nil {
			app.serverError(w, err)
			return
		}
		resp.NewNotifyCount = app.NewNotifyCounter(resp.Notify)
		err = t.Execute(w, resp)
		if err != nil {
			app.serverError(w, err)
			return
		}
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}
