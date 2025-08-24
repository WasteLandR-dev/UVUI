package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"uvui/internal/types"
	"uvui/internal/ui"
)

// Mock services for testing commands
type MockUVInstaller struct {
	mock.Mock
}

func (m *MockUVInstaller) IsInstalled() (bool, string, error) {
	args := m.Called()
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockUVInstaller) Install() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUVInstaller) GetInstallCommand() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockPythonManager struct {
	mock.Mock
}

func (m *MockPythonManager) ListAvailable() ([]types.PythonVersion, error) {
	args := m.Called()
	return args.Get(0).([]types.PythonVersion), args.Error(1)
}

func (m *MockPythonManager) ListInstalled() ([]types.PythonVersion, error) {
	args := m.Called()
	return args.Get(0).([]types.PythonVersion), args.Error(1)
}

func (m *MockPythonManager) Install(version string) error {
	args := m.Called(version)
	return args.Error(0)
}

func (m *MockPythonManager) Uninstall(version string) error {
	args := m.Called(version)
	return args.Error(0)
}

func (m *MockPythonManager) Pin(version string) error {
	args := m.Called(version)
	return args.Error(0)
}

func (m *MockPythonManager) Find(version string) (*types.PythonVersion, error) {
	args := m.Called(version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.PythonVersion), args.Error(1)
}

func TestCheckUVStatus(t *testing.T) {
	mockInstaller := &MockUVInstaller{}

	// Test successful check
	mockInstaller.On("IsInstalled").Return(true, "1.0.0", nil)

	cmd := CheckUVStatus(mockInstaller)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	uvMsg, ok := msg.(ui.UVInstalledMsg)

	assert.True(t, ok)
	assert.True(t, uvMsg.Success)
	assert.Equal(t, "1.0.0", uvMsg.Version)
	assert.Nil(t, uvMsg.Error)

	mockInstaller.AssertExpectations(t)
}

func TestCheckUVStatus_Error(t *testing.T) {
	mockInstaller := &MockUVInstaller{}

	// Test error case
	mockInstaller.On("IsInstalled").Return(false, "", assert.AnError)

	cmd := CheckUVStatus(mockInstaller)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	uvMsg, ok := msg.(ui.UVInstalledMsg)

	assert.True(t, ok)
	assert.False(t, uvMsg.Success)
	assert.Equal(t, "", uvMsg.Version)
	assert.NotNil(t, uvMsg.Error)

	mockInstaller.AssertExpectations(t)
}

func TestInstallUV(t *testing.T) {
	mockInstaller := &MockUVInstaller{}

	// Test successful installation
	mockInstaller.On("Install").Return(nil)

	cmd := InstallUV(mockInstaller)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	uvMsg, ok := msg.(ui.UVInstalledMsg)

	assert.True(t, ok)
	assert.True(t, uvMsg.Success)
	assert.Equal(t, "", uvMsg.Version)
	assert.Nil(t, uvMsg.Error)

	mockInstaller.AssertExpectations(t)
}

func TestInstallUV_Error(t *testing.T) {
	mockInstaller := &MockUVInstaller{}

	// Test installation error
	mockInstaller.On("Install").Return(assert.AnError)

	cmd := InstallUV(mockInstaller)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	uvMsg, ok := msg.(ui.UVInstalledMsg)

	assert.True(t, ok)
	assert.False(t, uvMsg.Success)
	assert.Equal(t, "", uvMsg.Version)
	assert.NotNil(t, uvMsg.Error)

	mockInstaller.AssertExpectations(t)
}

func TestLoadPythonVersions(t *testing.T) {
	mockManager := &MockPythonManager{}

	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	installed := []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true},
	}

	// Test successful loading
	mockManager.On("ListAvailable").Return(available, nil)
	mockManager.On("ListInstalled").Return(installed, nil)

	cmd := LoadPythonVersions(mockManager)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonMsg, ok := msg.(ui.PythonVersionsLoadedMsg)

	assert.True(t, ok)
	assert.Equal(t, available, pythonMsg.Available)
	assert.Equal(t, installed, pythonMsg.Installed)
	assert.Nil(t, pythonMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestLoadPythonVersions_AvailableError(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test error in available versions
	mockManager.On("ListAvailable").Return([]types.PythonVersion{}, assert.AnError)

	cmd := LoadPythonVersions(mockManager)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonMsg, ok := msg.(ui.PythonVersionsLoadedMsg)

	assert.True(t, ok)
	assert.Nil(t, pythonMsg.Available)
	assert.Nil(t, pythonMsg.Installed)
	assert.NotNil(t, pythonMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestLoadPythonVersions_InstalledError(t *testing.T) {
	mockManager := &MockPythonManager{}

	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
	}

	// Test error in installed versions
	mockManager.On("ListAvailable").Return(available, nil)
	mockManager.On("ListInstalled").Return([]types.PythonVersion{}, assert.AnError)

	cmd := LoadPythonVersions(mockManager)
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonMsg, ok := msg.(ui.PythonVersionsLoadedMsg)

	assert.True(t, ok)
	assert.Equal(t, available, pythonMsg.Available)
	assert.Nil(t, pythonMsg.Installed)
	assert.NotNil(t, pythonMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestInstallPythonVersion(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test successful installation
	mockManager.On("Install", "3.12.0").Return(nil)

	cmd := InstallPythonVersion(mockManager, "3.12.0")
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonOpMsg, ok := msg.(ui.PythonOperationMsg)

	assert.True(t, ok)
	assert.Equal(t, "install", pythonOpMsg.Operation)
	assert.Equal(t, "3.12.0", pythonOpMsg.Target)
	assert.True(t, pythonOpMsg.Success)
	assert.Nil(t, pythonOpMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestInstallPythonVersion_Error(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test installation error
	mockManager.On("Install", "3.12.0").Return(assert.AnError)

	cmd := InstallPythonVersion(mockManager, "3.12.0")
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonOpMsg, ok := msg.(ui.PythonOperationMsg)

	assert.True(t, ok)
	assert.Equal(t, "install", pythonOpMsg.Operation)
	assert.Equal(t, "3.12.0", pythonOpMsg.Target)
	assert.False(t, pythonOpMsg.Success)
	assert.NotNil(t, pythonOpMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestUninstallPythonVersion(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test successful uninstallation
	mockManager.On("Uninstall", "3.11.0").Return(nil)

	cmd := UninstallPythonVersion(mockManager, "3.11.0")
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonOpMsg, ok := msg.(ui.PythonOperationMsg)

	assert.True(t, ok)
	assert.Equal(t, "uninstall", pythonOpMsg.Operation)
	assert.Equal(t, "3.11.0", pythonOpMsg.Target)
	assert.True(t, pythonOpMsg.Success)
	assert.Nil(t, pythonOpMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestUninstallPythonVersion_Error(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test uninstallation error
	mockManager.On("Uninstall", "3.11.0").Return(assert.AnError)

	cmd := UninstallPythonVersion(mockManager, "3.11.0")
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonOpMsg, ok := msg.(ui.PythonOperationMsg)

	assert.True(t, ok)
	assert.Equal(t, "uninstall", pythonOpMsg.Operation)
	assert.Equal(t, "3.11.0", pythonOpMsg.Target)
	assert.False(t, pythonOpMsg.Success)
	assert.NotNil(t, pythonOpMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestPinPythonVersion(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test successful pinning
	mockManager.On("Pin", "3.11.0").Return(nil)

	cmd := PinPythonVersion(mockManager, "3.11.0")
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonOpMsg, ok := msg.(ui.PythonOperationMsg)

	assert.True(t, ok)
	assert.Equal(t, "pin", pythonOpMsg.Operation)
	assert.Equal(t, "3.11.0", pythonOpMsg.Target)
	assert.True(t, pythonOpMsg.Success)
	assert.Nil(t, pythonOpMsg.Error)

	mockManager.AssertExpectations(t)
}

func TestPinPythonVersion_Error(t *testing.T) {
	mockManager := &MockPythonManager{}

	// Test pinning error
	mockManager.On("Pin", "3.11.0").Return(assert.AnError)

	cmd := PinPythonVersion(mockManager, "3.11.0")
	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()
	pythonOpMsg, ok := msg.(ui.PythonOperationMsg)

	assert.True(t, ok)
	assert.Equal(t, "pin", pythonOpMsg.Operation)
	assert.Equal(t, "3.11.0", pythonOpMsg.Target)
	assert.False(t, pythonOpMsg.Success)
	assert.NotNil(t, pythonOpMsg.Error)

	mockManager.AssertExpectations(t)
}
