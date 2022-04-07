package delivery

import (
	"github.com/gorilla/mux"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/posts/usecase"
	"golang-stepik-2022q1/reditclone/pkg/session"
	"golang-stepik-2022q1/reditclone/pkg/utils/http_utils"
	"net/http"
)

func (h *Handler) Upvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	postId, ok := vars["postId"]
	if !ok {
		log.Clog(ctx).Info("Improper request params")
		http_utils.HttpError(w, "Wrong id provided", http.StatusBadRequest)
		return
	}

	sess := session.FromCtx(r.Context())
	if sess == nil {
		log.Rlog(r).Warn("Cant load session from request")
		http_utils.HttpError(w, "Cant get creator info from request", http.StatusInternalServerError)
		return
	}

	post, err := h.manager.Upvote(postId, sess.User.Id)
	if err != nil {
		if err == usecase.ItemNotFound {
			http_utils.HttpError(w, "Post not found", http.StatusNotFound)
			return
		}
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http_utils.JsonResp(w, post, http.StatusCreated)
}

func (h *Handler) Downvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	postId, ok := vars["postId"]
	if !ok {
		log.Clog(ctx).Info("Improper request params")
		http_utils.HttpError(w, "Wrong id provided", http.StatusBadRequest)
		return
	}

	sess := session.FromCtx(r.Context())
	if sess == nil {
		log.Rlog(r).Warn("Cant load session from request")
		http_utils.HttpError(w, "Cant get creator info from request", http.StatusInternalServerError)
		return
	}

	post, err := h.manager.Downvote(postId, sess.User.Id)
	if err != nil {
		if err == usecase.ItemNotFound {
			http_utils.HttpError(w, "Post not found", http.StatusNotFound)
			return
		}
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http_utils.JsonResp(w, post, http.StatusCreated)
}

func (h *Handler) Unvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	postId, ok := vars["postId"]
	if !ok {
		log.Clog(ctx).Info("Improper request params")
		http_utils.HttpError(w, "Wrong id provided", http.StatusBadRequest)
		return
	}

	sess := session.FromCtx(r.Context())
	if sess == nil {
		log.Rlog(r).Warn("Cant load session from request")
		http_utils.HttpError(w, "Cant get creator info from request", http.StatusInternalServerError)
		return
	}

	post, err := h.manager.Unvote(postId, sess.User.Id)
	if err != nil {
		if err == usecase.ItemNotFound {
			http_utils.HttpError(w, "Post not found", http.StatusNotFound)
			return
		}
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http_utils.JsonResp(w, post, http.StatusCreated)
}
