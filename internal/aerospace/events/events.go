package events

type Event = string

const (
	WorkspaceChange Event = "aerospace_workspace_change"
)

type WorkspaceChangeEventInfo struct {
	Focused string `json:"focused"`
	Prev    string `json:"prev"`
}
