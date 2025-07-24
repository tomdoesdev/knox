package workspace

import "os"

// WithLocalWorkspace handles getting the current working directory and finding
// the workspace, then calls the provided handler with the workspace.
// Returns raw errors without wrapping for maximum flexibility.
func WithLocalWorkspace(handler func(*Workspace) error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ws, err := FindWorkspace(cwd)
	if err != nil {
		return err
	}

	return handler(ws)
}

func WithEnsuredLocalWorkspace(handler func(*Workspace, InitResult) error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	result, err := EnsureWorkspace(cwd)
	if err != nil {
		return err
	}

	ws, err := FindWorkspace(cwd)
	if err != nil {
		return err
	}
	return handler(ws, result)
}
