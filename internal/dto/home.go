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

type HomeResponse struct {
	HeroCards           []HeroCard              `json:"hero_cards"`
	FeaturedCollections []FeaturedCollectionDTO `json:"featured_collections"`
	BestsellingProducts []ProductResponse       `json:"bestselling_products"`
}
