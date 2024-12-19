package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"WeChatOpenDevTools-AppleSilicon/pkg/frida"
)

// Commons handles common operations for the application
type Commons struct {
	WeChatUtils *WeChatUtils
	Injector    *frida.Injector
}

// NewCommons creates a new Commons instance
func NewCommons() *Commons {
	injector, err := frida.NewInjector()
	if err != nil {
		fmt.Printf("Warning: failed to create Frida injector: %v\n", err)
		return &Commons{
			WeChatUtils: NewWeChatUtils(),
		}
	}

	return &Commons{
		WeChatUtils: NewWeChatUtils(),
		Injector:    injector,
	}
}

// LoadWeChatConfigs loads WeChat configurations and injects the code
func (c *Commons) LoadWeChatConfigs() error {
	pid, version, err := c.WeChatUtils.GetWeChatPIDAndVersion()
	if err != nil {
		return fmt.Errorf("failed to get WeChat process info: %v", err)
	}

	return c.loadAndInject(pid, version)
}

// LoadWeChatConfigsPID loads WeChat configurations for a specific PID
func (c *Commons) LoadWeChatConfigsPID(pid int) error {
	version, err := c.WeChatUtils.GetWeChatVersion()
	if err != nil {
		return fmt.Errorf("failed to get WeChat version: %v", err)
	}

	return c.loadAndInject(pid, version)
}

// loadAndInject loads configuration and injects code into WeChat process
func (c *Commons) loadAndInject(pid int, version string) error {
	// Read hook.js
	hookPath := filepath.Join(filepath.Dir(c.WeChatUtils.ConfigsPath), "scripts", "hook.js")
	hookCode, err := ioutil.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("failed to read hook script: %v", err)
	}

	// Read address configuration
	addressPath := filepath.Join(c.WeChatUtils.ConfigsPath, fmt.Sprintf("address_%s_arm.json", version))
	addresses, err := ioutil.ReadFile(addressPath)
	if err != nil {
		return fmt.Errorf("failed to read address configuration: %v", err)
	}

	// Combine code
	finalCode := fmt.Sprintf("var address=%s\n%s", addresses, hookCode)

	// Inject code using Frida
	if c.Injector == nil {
		return fmt.Errorf("Frida injector not initialized")
	}

	fmt.Printf("Injecting code into process %d with version %s\n", pid, version)
	return c.Injector.InjectScript(pid, finalCode)
}