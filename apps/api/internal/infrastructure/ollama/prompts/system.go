package prompts

// PromptMode defines the extraction mode
type PromptMode string

const (
	PromptModeFast     PromptMode = "fast"
	PromptModeBalanced PromptMode = "balanced"
	PromptModeStrict   PromptMode = "strict"
)

// ============================================================================
// CHUNK SUMMARY PROMPTS
// ============================================================================

const chunkSummaryPromptFast = `Summarize this meeting section. Focus on key points.

You MUST return ONLY valid JSON in this exact format:
{"summary": "<2-3 sentences>"}

Do not include any text before or after the JSON.

Meeting section:
`

const chunkSummaryPromptBalanced = `Summarize this meeting section concisely. Focus on key points, decisions, and action items.

You MUST return ONLY valid JSON in this exact format:
{"summary": "<your 2-3 sentence summary>"}

Do not include any text before or after the JSON.

Meeting section:
`

const chunkSummaryPromptStrict = `You are Meetext AI.

Extract factual meeting information from this section.

RULES:
- Never invent information
- Never infer missing details
- Only extract explicitly stated information
- If unclear, omit it
- Remove duplicate information
- Ignore filler conversation

Focus on:
- action items with explicit owners
- finalized decisions
- technical discussions (APIs, databases, frameworks, infrastructure)
- blockers and risks explicitly mentioned
- ticket references

You MUST return ONLY valid JSON in this exact format:
{"summary": "<factual 2-4 sentence summary of this section>"}

Do not include any text before or after the JSON.

Meeting section:
`

// ============================================================================
// STRUCTURED EXTRACTION PROMPTS
// ============================================================================

const structuredExtractionPromptFast = `Extract structured information from these meeting summaries.

Rules:
1. Use null for missing data
2. Only extract explicitly stated information

Return valid JSON:
{
  "summary": "Overall meeting summary (2-3 sentences)",
  "tasks": [{"title": "", "description": "", "priority": null, "assignee": null, "due_date": null, "status": null, "confidence_score": 0.8}],
  "tickets": [{"title": "", "type": "task", "description": "", "status": "todo", "confidence_score": 0.8}],
  "decisions": [{"decision": "", "made_by": null, "confidence_score": 0.8}],
  "risks": [{"risk": "", "severity": null, "mitigation": null, "confidence_score": 0.8}],
  "blockers": [{"description": "", "confidence_score": 0.8}],
  "technical_notes": [{"topic": "", "details": ""}],
  "action_items": [{"description": "", "owner": null, "deadline": null, "confidence_score": 0.8}],
  "project_documentation_markdown": "# Meeting Documentation\n\n## Overview\n...\n\n## Decisions\n...\n\n## Next Steps\n..."
}

Meeting summaries:
`

const structuredExtractionPromptBalanced = `You are an expert meeting analyst. Extract structured information from these meeting summaries.

RULES:
1. Never invent information
2. Use null for missing data
3. Only extract explicitly stated information
4. Separate decisions from tasks
5. Separate blockers from risks

Return ONLY valid JSON matching this schema:
{
  "summary": "Overall meeting summary (2-3 sentences)",
  "tasks": [
    {
      "title": "Task name",
      "description": "What needs to be done",
      "priority": "low" | "medium" | "high" | null,
      "assignee": "Person name" | null,
      "due_date": "YYYY-MM-DD" | null,
      "status": "todo" | "in_progress" | "done" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "tickets": [
    {
      "title": "Ticket title",
      "type": "bug" | "feature" | "enhancement" | "task",
      "description": "Detailed description",
      "status": "todo" | "in_progress" | "done",
      "confidence_score": 0.0-1.0
    }
  ],
  "decisions": [
    {
      "decision": "What was decided",
      "made_by": "Person name" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "risks": [
    {
      "risk": "Potential future risk",
      "severity": "low" | "medium" | "high" | null,
      "mitigation": "Proposed solution" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "blockers": [
    {
      "description": "Current active blocker",
      "confidence_score": 0.0-1.0
    }
  ],
  "technical_notes": [
    {
      "topic": "Technical area",
      "details": "Technical discussion"
    }
  ],
  "action_items": [
    {
      "description": "Action to take",
      "owner": "Person name" | null,
      "deadline": "YYYY-MM-DD" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "project_documentation_markdown": "# Meeting Documentation\n\n## Overview\n...\n\n## Decisions\n...\n\n## Next Steps\n..."
}

Meeting summaries:
`

const structuredExtractionPromptStrict = `You are Meetext AI.

Extract factual meeting information from transcripts.

RULES:
- Never invent information
- Never infer missing details
- Only extract explicitly stated information
- If unclear, omit it
- Remove duplicate information
- Ignore filler conversation

IMPORTANT:
Meetings may contain MANY:
- action items
- decisions
- blockers
- tickets
- technical discussions

Extract ALL valid items.
Do NOT stop after finding one item.

EXTRACTION RULES:

Action Items:
- Include only explicitly assigned tasks (task + owner required)
- Optional: deadline, priority

Decisions:
- Include only finalized decisions
- Do NOT include suggestions, debates, or possibilities

Risks / Blockers:
- Include only explicitly mentioned blockers or issues
- Do NOT infer risks from context

Technical Discussions:
- Extract APIs, databases, infrastructure, AI systems, frontend/backend discussions

Tickets:
- Extract ticket IDs and descriptions if mentioned

DEDUPLICATION:
- The transcript may contain repeated OCR/ASR segments
- Remove duplicates, merge repeated ideas, ignore filler

Extract MULTIPLE items per section if they exist.
Scan the FULL transcript before generating output.

Return ONLY valid JSON. Do NOT include any text before or after the JSON:

{
  "summary": "Concise overall meeting summary (2-4 sentences)",
  "tasks": [
    {
      "title": "Task name",
      "description": "What needs to be done",
      "priority": "low" | "medium" | "high" | null,
      "assignee": "Person name" | null,
      "due_date": "YYYY-MM-DD" | null,
      "status": "todo" | "in_progress" | "done" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "tickets": [
    {
      "title": "Ticket ID or title",
      "type": "bug" | "feature" | "enhancement" | "task",
      "description": "Detailed description",
      "status": "todo" | "in_progress" | "done",
      "confidence_score": 0.0-1.0
    }
  ],
  "decisions": [
    {
      "decision": "What was decided",
      "made_by": "Person name or Team" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "risks": [
    {
      "risk": "Explicitly mentioned risk or concern",
      "severity": "low" | "medium" | "high" | null,
      "mitigation": "Proposed solution" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "blockers": [
    {
      "description": "Explicitly mentioned blocker",
      "confidence_score": 0.0-1.0
    }
  ],
  "technical_notes": [
    {
      "topic": "Technical area",
      "details": "Concrete technical discussion"
    }
  ],
  "action_items": [
    {
      "description": "Action to take",
      "owner": "Person name" | null,
      "deadline": "YYYY-MM-DD" | null,
      "confidence_score": 0.0-1.0
    }
  ],
  "project_documentation_markdown": "# Meeting Summary\n\n## Main Topics\n- item\n\n## Technical Discussions\n- item\n\n## Decisions Made\n- item\n\n## Action Items\n- Task — Owner — Due Date\n\n## Product Tickets\n- Ticket — Description\n\n## Risks / Blockers\n- item\n\n## Important Notes\n- item"
}

Meeting summaries:
`

// ============================================================================
// PROMPT BUILDERS
// ============================================================================

// BuildChunkSummaryPrompt creates a compact prompt for chunk summarization
func BuildChunkSummaryPrompt(chunk string, mode PromptMode) string {
	switch mode {
	case PromptModeFast:
		return chunkSummaryPromptFast + chunk
	case PromptModeStrict:
		return chunkSummaryPromptStrict + chunk
	default:
		return chunkSummaryPromptBalanced + chunk
	}
}

// BuildStructuredExtractionPrompt creates a prompt for final structured extraction
func BuildStructuredExtractionPrompt(mergedSummaries string, mode PromptMode) string {
	switch mode {
	case PromptModeFast:
		return structuredExtractionPromptFast + mergedSummaries
	case PromptModeStrict:
		return structuredExtractionPromptStrict + mergedSummaries
	default:
		return structuredExtractionPromptBalanced + mergedSummaries
	}
}
