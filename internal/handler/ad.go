package handler

import (
	"encoding/json"
	"log"
	"marketplace/internal/model"
	"marketplace/internal/store"
	"marketplace/util"
	"net/http"
	"strconv"
)

func NewHandler(store *store.DB) *Handler {
	return &Handler{store: store}
}

func (h *Handler) PostAd(w http.ResponseWriter, r *http.Request) {
	var ad model.Ad
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		log.Printf("Error decoding ad: %v", err)
		util.SendJSONError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Context().Value(util.ContextKey("username")) == nil {
		util.SendJSONError(w, r, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ad.Author = r.Context().Value(util.ContextKey("username")).(string)

	switch util.IsValidImage(ad.ImageURL) {
	case -1:
		util.SendJSONError(w, r, "Invalid image URL", http.StatusBadRequest)
		return
	case -2:
		util.SendJSONError(w, r, "Invalid image format", http.StatusBadRequest)
		return
	case -3:
		util.SendJSONError(w, r, "Invalid image format", http.StatusBadRequest)
		return
	case -4:
		util.SendJSONError(w, r, "Image too large", http.StatusBadRequest)
		return
	}

	if ad.Title == "" || ad.Description == "" || ad.ImageURL == "" || len([]rune(ad.ImageURL)) > 255 || len([]rune(ad.Title)) > 255 || len([]rune(ad.Description)) > 255 || ad.Price <= 0 || ad.Price > 1000000 {
		util.SendJSONError(w, r, "Invalid ad", http.StatusBadRequest)
		return
	}

	if err := h.store.CreateAd(&ad); err != nil {
		util.SendJSONError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendJSONResponse(w, r, ad, ad.Author, http.StatusCreated)
}

func (h *Handler) GetAds(w http.ResponseWriter, r *http.Request) {
	limit, offset, sortType, sortDirection := 10, 0, "created_at", "desc"
	priceMin, priceMax := 0.0, 1000000.0

	queryParams := r.URL.Query()

	if l, ok := queryParams["limit"]; ok && len(l[0]) > 0 {
		if newLimit, err := strconv.Atoi(l[0]); err == nil {
			limit = newLimit
		}
	}

	if o, ok := queryParams["offset"]; ok && len(o[0]) > 0 {
		if newOffset, err := strconv.Atoi(o[0]); err == nil {
			offset = newOffset
		}
	}

	if st, ok := queryParams["sortType"]; ok && len(st[0]) > 0 {
		sortType = st[0]
	}

	if sd, ok := queryParams["sortDirection"]; ok && len(sd[0]) > 0 {
		sortDirection = sd[0]
	}

	if pmin, ok := queryParams["priceMin"]; ok && len(pmin[0]) > 0 {
		if newPriceMin, err := strconv.ParseFloat(pmin[0], 64); err == nil {
			priceMin = newPriceMin
		}
	}

	if pmax, ok := queryParams["priceMax"]; ok && len(pmax[0]) > 0 {
		if newPriceMax, err := strconv.ParseFloat(pmax[0], 64); err == nil {
			priceMax = newPriceMax
		}
	}

	ads, total, err := h.store.GetAds(r.Context(), limit, offset, sortType, sortDirection, priceMin, priceMax)
	if err != nil {
		util.SendJSONError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Ads   []model.Ad `json:"ads"`
		Total int        `json:"total"`
	}{
		Ads:   ads,
		Total: total,
	}

	if r.Context().Value(util.ContextKey("username")) == nil {
		util.SendJSONResponse(w, r, response, "Unauthorized", http.StatusOK)
		return
	}

	util.SendJSONResponse(w, r, response, r.Context().Value(util.ContextKey("username")).(string), http.StatusOK)
}
