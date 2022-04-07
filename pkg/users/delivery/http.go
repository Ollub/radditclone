package delivery

import (
	"encoding/json"
	"golang-stepik-2022q1/reditclone/pkg/log"
	sessionUC "golang-stepik-2022q1/reditclone/pkg/session/usecase"
	"golang-stepik-2022q1/reditclone/pkg/users"
	"golang-stepik-2022q1/reditclone/pkg/users/usecase"
	"golang-stepik-2022q1/reditclone/pkg/utils/http_utils"
	"net/http"
)

type Handler struct {
	manager        *usecase.Manager
	sessionManager *sessionUC.Manager
}

func NewHandler(manager *usecase.Manager, sessionManager *sessionUC.Manager) *Handler {
	return &Handler{
		manager:        manager,
		sessionManager: sessionManager,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIn, err := http_utils.FromBody[users.UserIn](r)
	if err != nil {
		log.Clog(ctx).Error("Marshaling error")
		http_utils.HttpError(w, "Marshaling error", http.StatusInternalServerError)
		return
	}

	user, err := h.manager.Create(ctx, userIn)
	if err == usecase.UserExistsError {
		http_utils.HttpError(w, "UserId exist", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Clog(ctx).Error("Unexpected error during user creation", log.Fields{"err": err.Error()})
		http_utils.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.sessionManager.IssueToken(r.Context(), user)
	if err != nil {
		log.Clog(ctx).Error("Cant issue token", log.Fields{"err": err.Error()})
		http_utils.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := &LoginResp{token}
	resp, err := json.Marshal(out)
	if err != nil {
		log.Clog(ctx).Error("Marshaling error")
		http_utils.HttpError(w, "Marshaling error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Clog(ctx).Error("Cant write response")
		http_utils.HttpError(w, "Cant write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	in, err := http_utils.FromBody[LoginReq](r)
	if err != nil {
		log.Rlog(r).Error("Marshaling error")
		http_utils.HttpError(w, "Marshaling error", http.StatusInternalServerError)
		return
	}

	user, err := h.manager.GetByName(r.Context(), in.Username)

	if user == nil {
		http_utils.HttpError(w, "UserId not found", http.StatusUnauthorized)
		return
	}

	passValid := usecase.CheckPassword(user.PassHash, in.Password)
	if !passValid {
		http_utils.HttpError(w, "invalid password", http.StatusUnauthorized)
		return
	}

	token, err := h.sessionManager.IssueToken(r.Context(), user)
	if err != nil {
		http_utils.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := &LoginResp{token}
	resp, err := json.Marshal(out)
	if err != nil {
		http_utils.HttpError(w, "Marshaling error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		http_utils.HttpError(w, "Cant write response", http.StatusInternalServerError)
	}
}
