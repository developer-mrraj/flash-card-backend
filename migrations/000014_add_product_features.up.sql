-- Add features JSON column to products table
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS features JSONB DEFAULT '[]'::jsonb;

-- Update features for Shivaji Maharaj deck
UPDATE products SET features = '[
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M4 19.5A2.5 2.5 0 0 1 6.5 17H20\"></path><path d=\"M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z\"></path></svg>",
    "title": "Chronological Journey",
    "description": "From the birth at Shivneri to the Grand Coronation at Raigad, every major life event is captured with vivid illustrations and key facts."
  },
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z\"></path></svg>",
    "title": "Warfare Tactics",
    "description": "Visual guides to the Ganimi Kava (Guerrilla tactics), naval strategy, and the architectural brilliance of the hill forts."
  },
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polygon points=\"12 2 2 7 12 12 22 7 12 2\"></polygon><polyline points=\"2 17 12 22 22 17\"></polyline><polyline points=\"2 12 12 17 22 12\"></polyline></svg>",
    "title": "Civic Governance",
    "description": "Detailed breakdown of the Ashta Pradhan Mandal and the administration policies that built a sovereign Maratha state."
  }
]'::jsonb WHERE slug = 'shivaji-maratha';

-- Update features for Kids A to Z deck
UPDATE products SET features = '[
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z\"></path><path d=\"M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z\"></path></svg>",
    "title": "Alphabet Adventure",
    "description": "Learn A to Z with beautifully illustrated cards that connect letters to Indian heritage, culture, and iconic monuments."
  },
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><circle cx=\"12\" cy=\"12\" r=\"10\"></circle><path d=\"M8 14s1.5 2 4 2 4-2 4-2\"></path><line x1=\"9\" y1=\"9\" x2=\"9.01\" y2=\"9\"></line><line x1=\"15\" y1=\"9\" x2=\"15.01\" y2=\"9\"></line></svg>",
    "title": "Engaging Visuals",
    "description": "Keep kids focused with vibrant, high-contrast imagery designed specifically for early childhood cognitive development."
  },
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M12 2v20\"></path><path d=\"M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6\"></path></svg>",
    "title": "Fun Fact Learning",
    "description": "Each card includes bite-sized, easy-to-read facts on the back, turning playtime into a fun and educational experience."
  }
]'::jsonb WHERE slug = 'a-to-z-indian-heritage';

-- Update features for NEET Toppers Blueprint deck
UPDATE products SET features = '[
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M2 12h4l3-9 5 18 3-9h5\"></path></svg>",
    "title": "High-Yield Topics",
    "description": "Focus strictly on NCERT-aligned content, covering the most frequently tested concepts in Biology, Physics, and Chemistry."
  },
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><rect x=\"3\" y=\"3\" width=\"18\" height=\"18\" rx=\"2\" ry=\"2\"></rect><line x1=\"3\" y1=\"9\" x2=\"21\" y2=\"9\"></line><line x1=\"9\" y1=\"21\" x2=\"9\" y2=\"9\"></line></svg>",
    "title": "Visual Mnemonics",
    "description": "Memorize complex pathways, formulas, and diagrams instantly using proven visual memory techniques."
  },
  {
    "icon": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polygon points=\"12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2\"></polygon></svg>",
    "title": "Topper''s Strategies",
    "description": "Learn how top rankers connect dots between chapters with specially curated revision and interlinking cards."
  }
]'::jsonb WHERE slug = 'neet-toppers-blueprint';
