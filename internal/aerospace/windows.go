package aerospace

type Window struct {
	ID  WindowID
	App string
}

type FullWindow struct {
	ID          WindowID
	App         string
	WorkspaceID WorkspaceID
	MonitorID   MonitorID
}
