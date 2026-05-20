package ai

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Priority    string `json:"priority"` // low | medium | high
	DueDate     string `json:"due_date"`
}

type Ticket struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"` // bug | feature | enhancement
	Status      string `json:"status"` // todo | in_progress | done
}

type Decision struct {
	Description string `json:"description"`
	MadeBy      string `json:"made_by"`
}

type Risk struct {
	Description string `json:"description"`
	Severity    string `json:"severity"` // low | medium | high
	Mitigation  string `json:"mitigation"`
}

type ActionItem struct {
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Deadline    string `json:"deadline"`
}

type ClientRequest struct {
	Request     string `json:"request"`
	Priority    string `json:"priority"`
	IsCommitted bool   `json:"is_committed"`
}

type TechnicalNote struct {
	Topic   string `json:"topic"`
	Details string `json:"details"`
}

type AIResult struct {
	Summary              string          `json:"summary"`
	Tasks                []Task          `json:"tasks"`
	Tickets              []Ticket        `json:"tickets"`
	Decisions            []Decision      `json:"decisions"`
	Risks                []Risk          `json:"risks"`
	ActionItems          []ActionItem    `json:"action_items"`
	TechnicalNotes       []TechnicalNote `json:"technical_notes"`
	ClientRequests       []ClientRequest `json:"client_requests"`
	ProjectDocumentation string          `json:"project_documentation"`
}
