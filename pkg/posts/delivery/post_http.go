package delivery

import (
	"github.com/gorilla/mux"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/posts"
	"golang-stepik-2022q1/reditclone/pkg/posts/usecase"
	"golang-stepik-2022q1/reditclone/pkg/session"
	"golang-stepik-2022q1/reditclone/pkg/utils/http_utils"
	"net/http"
)

type Handler struct {
	manager *usecase.Manager
}

func NewHandler(manager *usecase.Manager) *Handler {
	return &Handler{manager: manager}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.manager.GetAll(r.Context())
	if err != nil {
		http_utils.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, items, http.StatusOK)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item, err := h.manager.Get(r.Context(), vars["id"])
	if err != nil {
		http_utils.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		http_utils.HttpError(w, "Item not found", http.StatusNotFound)
		return
	}
	http_utils.JsonResp(w, item, http.StatusOK)
}

func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		log.Clog(ctx).Info("Improper request params")
		http_utils.HttpError(w, "Wrong id provided", http.StatusBadRequest)
		return
	}

	items, err := h.manager.FilterByUser(r.Context(), username)
	if err != nil {
		http_utils.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, items, http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	postIn, err := http_utils.FromBody[posts.PostIn](r)
	if err != nil {
		http_utils.HttpError(w, "Unmarshalling error", http.StatusInternalServerError)
		return
	}
	err = http_utils.Validate(postIn)
	if err != nil {
		http_utils.HttpError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	sess := session.FromCtx(r.Context())
	if sess == nil {
		log.Rlog(r).Warn("Cant load session from request")
		http_utils.HttpError(w, "Cant get creator info from request", http.StatusInternalServerError)
		return
	}
	postIn.Author.ID = sess.User.Id
	postIn.Author.Username = sess.User.Username
	log.Rlog(r).Info("postIn data collected", log.Fields{"postId": postIn, "sess": sess})
	post, err := h.manager.Create(r.Context(), postIn)
	if err != nil {
		log.Rlog(r).Warn("Cant create users")
		http_utils.HttpError(w, "Cant create users", http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, post, http.StatusCreated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	postId := vars["postId"]

	if postId == "" {
		log.Clog(ctx).Info("Improper request params")
		http_utils.HttpError(w, "Wrong request params", http.StatusBadRequest)
		return
	}
	post, err := h.manager.DeletePost(ctx, postId)
	if err != nil {
		if err == usecase.ItemNotFound {
			http_utils.HttpError(w, "Post not found", http.StatusNotFound)
			return
		}
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http_utils.JsonResp(w, post, http.StatusOK)
}
