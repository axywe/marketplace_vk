package handler

import (
	"encoding/json"
	"marketplace/internal/model"
	"marketplace/internal/store"
	"marketplace/pkg/hash"
	"marketplace/pkg/jwt"
	"marketplace/util"
	"net/http"
	"regexp"
)

type Handler struct {
	store *store.DB
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		util.SendJSONError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Username == "You" || u.Username == "you" || len([]rune(u.Username)) > 64 || len([]rune(u.Password)) > 64 || len([]rune(u.Username)) < 3 || len([]rune(u.Password)) < 3 || !regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(u.Username) || !regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(u.Password) {
		util.SendJSONError(w, r, "Invalid user", http.StatusBadRequest)
		return
	}

	if err := h.store.CreateUser(&u); err != nil {
		if err.Error() == "User already exists" {
			util.SendJSONError(w, r, err.Error(), http.StatusConflict)
			return
		}
		util.SendJSONError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	u.Password = "" // Don't return the password
	util.SendJSONResponse(w, r, u, u.Username, http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		util.SendJSONError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	storedUser, err := h.store.GetUser(u.Username)
	if err != nil {
		util.SendJSONError(w, r, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !hash.CheckPasswordHash(u.Password, storedUser.Password) {
		util.SendJSONError(w, r, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(storedUser.Username, h.store.Config.JWTSecret)
	if err != nil {
		util.SendJSONError(w, r, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"token":    token,
		"username": storedUser.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	util.SendJSONResponse(w, r, resp, storedUser.Username, http.StatusOK)
}
