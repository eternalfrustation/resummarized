CREATE TABLE IF NOT EXISTS article (
    article_id SERIAL PRIMARY KEY,

    headline_title VARCHAR(255) NOT NULL,
    lead_paragraph TEXT NOT NULL,

    background_context TEXT NOT NULL,
    research_question TEXT NOT NULL,

    simplified_methods TEXT NOT NULL,

    core_findings TEXT NOT NULL,
    surprise_finding TEXT,

    future_implications TEXT NOT NULL,
    study_limitations TEXT NOT NULL,
    next_steps TEXT NOT NULL,

    -- Metadata
    date_published DATE DEFAULT CURRENT_DATE
);
