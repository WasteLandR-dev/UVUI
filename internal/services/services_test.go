package services

// mockCommandExecutor is a mock implementation of the CommandExecutorInterface.
type mockCommandExecutor struct {
	IsUVAvailableFunc func() bool
	ExecuteFunc       func(command string, args ...string) ([]byte, error)
	RunCommandFunc    func(command string, args ...string) ([]byte, error)
}

func (m *mockCommandExecutor) IsUVAvailable() bool {
	if m.IsUVAvailableFunc != nil {
		return m.IsUVAvailableFunc()
	}
	return true
}

func (m *mockCommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command, args...)
	}
	return nil, nil
}

func (m *mockCommandExecutor) RunCommand(command string, args ...string) ([]byte, error) {
	if m.RunCommandFunc != nil {
		return m.RunCommandFunc(command, args...)
	}
	return nil, nil
}
