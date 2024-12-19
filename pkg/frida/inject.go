package frida

import (
	"bufio"
	"fmt"
	"log"
	"os"

	frida "github.com/frida/frida-go/frida"
)

// Injector handles Frida injection operations
type Injector struct {
	manager *frida.DeviceManager
	device  frida.DeviceInt
}

// NewInjector creates a new Frida injector
func NewInjector() (*Injector, error) {
	manager := frida.NewDeviceManager()
	
	device, err := manager.LocalDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get local device: %v", err)
	}

	return &Injector{
		manager: manager,
		device:  device,
	}, nil
}

// InjectScript injects JavaScript code into the target process
func (i *Injector) InjectScript(pid int, code string) error {
	// Attach to the target process
	session, err := i.device.Attach(pid, nil)
	if err != nil {
		return fmt.Errorf("failed to attach to process %d: %v", pid, err)
	}
	defer session.Detach()

	// Create script
	script, err := session.CreateScript(code)
	if err != nil {
		return fmt.Errorf("failed to create script: %v", err)
	}
	defer script.Unload()

	// Set message handler
	script.On("message", func(msg string) {
		log.Printf("[*] Received: %s\n", msg)
	})

	// Load script
	if err := script.Load(); err != nil {
		return fmt.Errorf("failed to load script: %v", err)
	}

	// Keep the script running until user input
	fmt.Println("Press Enter to stop...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	return nil
}
