package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.index)

	mux.HandleFunc("/user/register", app.register)
	mux.HandleFunc("/user/profile", app.profile)
	mux.HandleFunc("/user/login", app.login)
	mux.HandleFunc("/user/logout", app.logout)

	mux.HandleFunc("/post/create", app.createPost)
	mux.HandleFunc("/post/", app.postById)
	mux.HandleFunc("/post/like/", app.likePost)
	mux.HandleFunc("/post/dislike/", app.dislikePost)
	mux.HandleFunc("/post/remove/", app.removePost)
	mux.HandleFunc("/edit/", app.editPost)

	mux.HandleFunc("/category/", app.postByCategory)
	mux.HandleFunc("/favorite/", app.favoritePosts)
	mux.HandleFunc("/my/", app.myPosts)
	mux.HandleFunc("/postsByComment/", app.postsByComments)
	mux.HandleFunc("/comment/remove/", app.removeComment)
	mux.HandleFunc("/commentEdit/", app.editComment)

	mux.HandleFunc("/comment/create", app.createComment)
	mux.HandleFunc("/comment/like/", app.likeComment)
	mux.HandleFunc("/comment/dislike/", app.dislikeComment)

	mux.HandleFunc("/notify", app.showNotify)

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	return app.logRequest(app.limit(app.secureHeaders(app.recoverPanic(mux))))
}
