package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/posts"
	"golang-stepik-2022q1/reditclone/pkg/posts/usecase"
	"golang-stepik-2022q1/reditclone/pkg/session"
	"golang-stepik-2022q1/reditclone/pkg/utils/http_utils"
	"io/ioutil"
	"net/http"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, ok := vars["postId"]
	if !ok {
		log.Clog(ctx).Info("Improper request params", log.Fields{"id": id})
		http_utils.HttpError(w, "Wrong id provided", http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	in := &posts.CommentIn{}
	err := json.Unmarshal(body, in)
	if err != nil {
		log.Clog(ctx).Warn("Cant unmarshal comment", log.Fields{"err": err.Error()})
		http_utils.HttpError(w, "Unmarshaling error", http.StatusInternalServerError)
		return
	}

	sess := session.FromCtx(r.Context())
	if sess == nil {
		log.Rlog(r).Warn("Cant load session from request")
		http_utils.HttpError(w, "Cant get creator info from request", http.StatusInternalServerError)
		return
	}
	in.Author.ID = sess.User.Id
	in.Author.Username = sess.User.Username

	post, err := h.manager.CreateComment(ctx, id, in)
	if err != nil {
		if err == usecase.ItemNotFound {
			http_utils.HttpError(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Clog(ctx).Error("Cant create comment", log.Fields{"error": err.Error()})
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http_utils.JsonResp(w, post, http.StatusCreated)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	postId := vars["postId"]
	commentId := vars["commentId"]

	if postId == "" || commentId == "" {
		log.Clog(ctx).Info("Improper request params")
		http_utils.HttpError(w, "Wrong request params", http.StatusBadRequest)
		return
	}
	post, err := h.manager.DeleteComment(ctx, postId, commentId)
	if err != nil {
		if err == usecase.ItemNotFound {
			http_utils.HttpError(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Clog(ctx).Error("Cant delete comment", log.Fields{"error": err.Error()})
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http_utils.JsonResp(w, post, http.StatusCreated)
}
