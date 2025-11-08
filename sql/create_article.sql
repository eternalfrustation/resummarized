INSERT INTO article (
    headline_title,
    lead_paragraph,
    background_context,
    research_question,
    simplified_methods,
    core_findings,
    surprise_finding,
    future_implications,
    study_limitations,
    next_steps
) VALUES (
    $1, -- headline_title
    $2, -- lead_paragraph
    $3, -- background_context
    $4, -- research_question
    $5, -- simplified_methods
    $6, -- core_findings
    $7, -- surprise_finding (optional TEXT, can be NULL)
    $8, -- future_implications
    $9, -- study_limitations
    $10  -- next_steps
)
RETURNING *;
