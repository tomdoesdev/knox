package workspace

type (
	State struct {
		ActiveProject string `json:"active_project"`
	}
)

func NewStateDefault() *State {
	return &State{
		ActiveProject: defaultProjectName,
	}
}
