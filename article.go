package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/UniquityVentures/resummarized/pages"
)

type Article struct {
	ArticleID int `json:"article_id" comment:"Database ID for internal tracking. 0 when c"`

	HeadlineTitle string `json:"headline_title" comment:"The catchy, non-technical title/headline."`
	LeadParagraph string `json:"lead_paragraph" comment:"The opening paragraph that states the main finding and its real-world impact (the 'Why You Care')."`

	BackgroundContext string `json:"background_context" comment:"Simplified background on the large problem or field of study."`
	ResearchQuestion  string `json:"research_question" comment:"The focused question the specific research paper set out to answer, framed as a mystery or challenge."`

	SimplifiedMethods string `json:"simplified_methods" comment:"A non-expert friendly description of the methodology (what they did, not how the instruments work)."`

	CoreFindings    string `json:"core_findings" comment:"The primary result/conclusion, stated simply and focused on magnitude or effect."`
	SurpriseFinding string `json:"surprise_finding,omitempty" comment:"Any unexpected or counter-intuitive results found (optional)."`

	FutureImplications string `json:"future_implications" comment:"What the discovery means for humanity, technology, or the field (the future potential)."`
	StudyLimitations   string `json:"study_limitations" comment:"The most critical limitation or caveat of the study (e.g., 'only conducted in mice')."`
	NextSteps          string `json:"next_steps" comment:"The immediate next phase of research or development."`

	Author        string    `json:"author" comment:"The name of the article writer (not necessarily the scientist)."`
	DatePublished time.Time `json:"date_published" comment:"Date the summary article was published."`
}

func (app *App) CreateArticle(ctx context.Context, article pages.PostCreateForm) {
	create_post, err := os.ReadFile("sql/create_post.sql")
	if err != nil {
		log.Printf("Unable to find sql for creating posts table: %v\n", err)
		os.Exit(1)
	}
	app.Db.QueryRow(ctx, string(create_post),
		article.HeadlineTitle,      // $1
		article.LeadParagraph,      // $2
		article.BackgroundContext,  // $3
		article.ResearchQuestion,   // $4
		article.SimplifiedMethods,  // $5
		article.CoreFindings,       // $6
		article.SurpriseFindings,   // $7 (If empty string is passed, it stores "" not NULL)
		article.FutureImplications, // $8
		article.StudyLimitations,   // $9
		article.NextSteps,          // $10
	)
}
