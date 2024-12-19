package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// WeChatUtils handles WeChat process and version related operations
type WeChatUtils struct {
	ConfigsPath string
	VersionList []string
}

// NewWeChatUtils creates a new WeChatUtils instance
func NewWeChatUtils() *WeChatUtils {
	wu := &WeChatUtils{}
	wu.ConfigsPath = wu.GetConfigsPath()
	wu.VersionList = wu.GetVersionList()
	return wu
}

// GetConfigsPath returns the path to configs directory
func (w *WeChatUtils) GetConfigsPath() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(execPath), "configs")
}

// GetVersionList returns list of supported WeChat versions
func (w *WeChatUtils) GetVersionList() []string {
	files, err := os.ReadDir(w.ConfigsPath)
	if err != nil {
		return nil
	}

	var versions []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "address_") {
			version := strings.Split(file.Name(), "_")[1]
			versions = append(versions, version)
		}
	}
	return versions
}

// GetWeChatPIDAndVersion returns the PID and version of running WeChat process
func (w *WeChatUtils) GetWeChatPIDAndVersion() (int, string, error) {
	pidCmd := exec.Command("bash", "-c", 
		"ps aux | grep 'WeChatAppEx' | grep -v 'grep' | grep ' --client_version' | grep '-user-agent=' | awk '{print $2}' | tail -n 1")
	
	versionCmd := exec.Command("bash", "-c",
		"ps aux | grep 'WeChatAppEx' | grep -v 'grep' | grep ' --client_version' | grep '-user-agent=' | grep -oE 'MacWechat/([0-9]+\\.)+[0-9]+\\(0x\\d+\\)' | grep -oE '(0x\\d+)' | sed 's/0x//g' | head -n 1")

	pidOutput, err := pidCmd.Output()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get PID: %v", err)
	}

	versionOutput, err := versionCmd.Output()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get version: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidOutput)))
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse PID: %v", err)
	}

	version := strings.TrimSpace(string(versionOutput))
	return pid, version, nil
}

// GetWeChatVersion returns only the version of running WeChat process
func (w *WeChatUtils) GetWeChatVersion() (string, error) {
	cmd := exec.Command("bash", "-c",
		"ps aux | grep 'WeChatAppEx' | grep -v 'grep' | grep ' --client_version' | grep '-user-agent=' | grep -oE 'MacWechat/([0-9]+\\.)+[0-9]+\\(0x\\d+\\)' | grep -oE '(0x\\d+)' | sed 's/0x//g' | head -n 1")
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}