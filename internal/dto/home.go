package dto

type HeroCard struct {
	Image       string `json:"image"`
	Badge       string `json:"badge"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Stats       string `json:"stats"`
	Rating      string `json:"rating"`
}

type FeaturedCollectionDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

// BannerDTO represents a promotional ad banner returned to the frontend.
type BannerDTO struct {
	ID        string `json:"id"`
	Slot      string `json:"slot"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	CtaText   string `json:"cta_text"`
	CtaLink   string `json:"cta_link"`
	ImageURL  string `json:"image_url"`
	BgColor   string `json:"bg_color"`
	TextColor string `json:"text_color"`
	BadgeText string `json:"badge_text"`
	IsActive  bool   `json:"is_active"`
	SortOrder int    `json:"sort_order"`
}

type HomeResponse struct {
	HeroCards           []HeroCard              `json:"hero_cards"`
	FeaturedCollections []FeaturedCollectionDTO `json:"featured_collections"`
	BestsellingProducts []ProductResponse       `json:"bestselling_products"`
	Banners             []BannerDTO             `json:"banners"`
}
