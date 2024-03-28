package store

import (
	"context"
	"marketplace/internal/model"
	"marketplace/util"
)

func (db *DB) CreateAd(ad *model.Ad) error {
	return db.QueryRow("INSERT INTO ads(title, description, image_url, price, author) VALUES($1, $2, $3, $4, $5) RETURNING id",
		ad.Title, ad.Description, ad.ImageURL, ad.Price, ad.Author).Scan(&ad.ID)
}

func (db *DB) GetAds(ctx context.Context, limit, offset int, sortType, sortDirection string, priceMin, priceMax float64) ([]model.Ad, int, error) {
	totalQuery := `SELECT COUNT(*) FROM ads WHERE price >= $1 AND price <= $2`
	var total int
	err := db.QueryRowContext(ctx, totalQuery, priceMin, priceMax).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	adsQuery := `SELECT id, title, description, image_url, price, author FROM ads WHERE price >= $3 AND price <= $4 ORDER BY ` + sortType + ` ` + sortDirection + ` LIMIT $1 OFFSET $2`
	rows, err := db.QueryContext(ctx, adsQuery, limit, offset, priceMin, priceMax)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ads []model.Ad
	for rows.Next() {
		var ad model.Ad
		if err := rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.Author); err != nil {
			return nil, 0, err
		}

		username := ctx.Value(util.ContextKey("username"))
		if username != nil {
			if ad.Author == ctx.Value(util.ContextKey("username")).(string) {
				ad.Author = "You"
			}
		}

		ads = append(ads, ad)
	}

	return ads, total, nil
}
