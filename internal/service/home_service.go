package service

import (
	"backend/internal/dto"
	"backend/internal/repository"
)

type HomeService interface {
	GetHomeContent() (*dto.HomeResponse, error)
}

type homeService struct {
	productRepo            repository.ProductRepository
	featuredCollectionRepo repository.FeaturedCollectionRepository
}

func NewHomeService(productRepo repository.ProductRepository, featuredCollectionRepo repository.FeaturedCollectionRepository) HomeService {
	return &homeService{
		productRepo:            productRepo,
		featuredCollectionRepo: featuredCollectionRepo,
	}
}

func (s *homeService) GetHomeContent() (*dto.HomeResponse, error) {
	// 1. Hardcoded Hero Cards
	heroCards := []dto.HeroCard{
		{
			Image:       "/images/banner_249.png",
			Badge:       "OFFER",
			Title:       "Complete Learning Flash Cards",
			Description: "Build Knowledge, Strengthen Memory & Excel in Every Step! Covers History, Heritage & Education Topics. Perfect for kids, students, and competitive exams.",
			Stats:       "🎴 3 Powerful Sets",
			Rating:      "⭐ 5.0/5",
		},
		{
			Image:       "/images/hero_shivaji.png",
			Badge:       "HISTORY",
			Title:       "Shivaji Maharaj Era",
			Description: "Master the legend of the Maratha Empire. Our premium visual flashcards cover the Grand Coronation, Guerrilla Warfare tactics, Administration, and the architectural brilliance of the Sea Forts.",
			Stats:       "🎴 52 Cards",
			Rating:      "⭐ 4.9/5",
		},
		{
			Image:       "/images/hero_ramayana.png",
			Badge:       "CULTURE",
			Title:       "Ramayana Visuals",
			Description: "Visual storytelling through flashcards covering all 7 Kandas of the Ramayana. Experience the epic journey of Lord Rama, Sita, and Hanuman with stunning artwork.",
			Stats:       "🎴 108 Cards",
			Rating:      "⭐ 4.8/5",
		},
		{
			Image:       "/images/hero_polity.png",
			Badge:       "EXAM PRO",
			Title:       "UPSC Polity Series",
			Description: "Complete coverage of Laxmikanth, Constitutional Articles, and Amendments. Perfect for quick revisions before the Prelims and Mains examinations.",
			Stats:       "🎴 250 Cards",
			Rating:      "⭐ 4.7/5",
		},
	}

	// 2. Fetch Featured Collections from DB (non-fatal if table not yet migrated)
	var featuredDTOs []dto.FeaturedCollectionDTO
	if collections, fcErr := s.featuredCollectionRepo.FindAll(); fcErr == nil {
		for _, c := range collections {
			featuredDTOs = append(featuredDTOs, dto.FeaturedCollectionDTO{
				ID:    c.ID.String(),
				Title: c.Title,
				Image: c.Image,
			})
		}
	}
	if featuredDTOs == nil {
		featuredDTOs = []dto.FeaturedCollectionDTO{}
	}

	// 3. Fetch Bestselling Products from DB
	products, err := s.productRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var productDTOs []dto.ProductResponse
	for _, p := range products {
		productDTOs = append(productDTOs, mapToProductResponse(&p))
	}

	// Assemble Response
	return &dto.HomeResponse{
		HeroCards:           heroCards,
		FeaturedCollections: featuredDTOs,
		BestsellingProducts: productDTOs,
	}, nil
}
