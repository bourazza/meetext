package ai

type Task struct {
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Priority        *string `json:"priority"` // low | medium | high | null
	Assignee        *string `json:"assignee"` // string | null
	DueDate         *string `json:"due_date"` // YYYY-MM-DD | null
	Status          *string `json:"status"`   // todo | in_progress | done | null
	ConfidenceScore float64 `json:"confidence_score"`
}

type Ticket struct {
	Title           string  `json:"title"`
	Type            string  `json:"type"`   // bug | feature | enhancement | task
	Description     string  `json:"description"`
	Status          string  `json:"status"` // todo | in_progress | done
	ConfidenceScore float64 `json:"confidence_score"`
}

type Decision struct {
	Decision        string  `json:"decision"`
	MadeBy          *string `json:"made_by"` // string | null
	ConfidenceScore float64 `json:"confidence_score"`
}

type Risk struct {
	Risk            string  `json:"risk"`
	Severity        *string `json:"severity"`   // low | medium | high | null
	Mitigation      *string `json:"mitigation"` // string | null
	ConfidenceScore float64 `json:"confidence_score"`
}

type Blocker struct {
	Description     string  `json:"description"`
	ConfidenceScore float64 `json:"confidence_score"`
}

type ActionItem struct {
	Description     string  `json:"description"`
	Owner           *string `json:"owner"`
	Deadline        *string `json:"deadline"`
	ConfidenceScore float64 `json:"confidence_score"`
}

type TechnicalNote struct {
	Topic   string `json:"topic"`
	Details string `json:"details"`
}

type AIResult struct {
	Summary                      string          `json:"summary"`
	Tasks                        []Task          `json:"tasks"`
	Tickets                      []Ticket        `json:"tickets"`
	Decisions                    []Decision      `json:"decisions"`
	Risks                        []Risk          `json:"risks"`
	Blockers                     []Blocker       `json:"blockers"`
	TechnicalNotes               []TechnicalNote `json:"technical_notes"`
	ActionItems                  []ActionItem    `json:"action_items"`
	ProjectDocumentationMarkdown string          `json:"project_documentation_markdown"`
}
