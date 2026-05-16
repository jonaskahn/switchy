//go:build windows

package registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	clientName   = "Switchy"
	progIDSuffix = "SwitchyHTML"
)

func Register(exePath string) error {
	if exePath == "" {
		var err error
		exePath, err = os.Executable()
		if err != nil {
			return fmt.Errorf("registry: resolve exe: %w", err)
		}
		exePath, _ = filepath.Abs(exePath)
	}

	base := `SOFTWARE\Clients\StartMenuInternet\` + clientName
	capPath := base + `\Capabilities`

	steps := []func() error{
		func() error { return writeClient() },
		func() error { return writeProgID(exePath) },
		func() error { return writeCapabilities(capPath) },
		func() error { return writeURLAssociations(capPath) },
		func() error { return writeShellOpenCommand(base, exePath) },
		func() error { return writeRegisteredApplication(capPath) },
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func Unregister() error {
	base := `SOFTWARE\Clients\StartMenuInternet\` + clientName
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, base+`\Capabilities\URLAssociations`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, base+`\Capabilities`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, base+`\shell\open\command`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, base+`\shell\open`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, base+`\shell`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, base)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\`+progIDSuffix+`\shell\open\command`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\`+progIDSuffix+`\shell\open`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\`+progIDSuffix+`\shell`)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\`+progIDSuffix)

	ra, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\RegisteredApplications`, registry.ALL_ACCESS)
	if err == nil {
		_ = ra.DeleteValue(clientName)
		ra.Close()
	}
	return nil
}

func OpenDefaultAppsSettings() error {
	return exec.Command("cmd", "/c", "start", "ms-settings:defaultapps").Start()
}

func createKey(path, label string) (registry.Key, error) {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return 0, fmt.Errorf("registry: create %s key: %w", label, err)
	}
	return k, nil
}

func writeClient() error {
	base := `SOFTWARE\Clients\StartMenuInternet\` + clientName
	k, err := createKey(base, "client")
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue("", clientName)
}

func writeProgID(exePath string) error {
	k, err := createKey(`SOFTWARE\Classes\`+progIDSuffix+`\shell\open\command`, "ProgID")
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue("", fmt.Sprintf(`"%s" "%%1"`, exePath))
}

func writeCapabilities(capPath string) error {
	k, err := createKey(capPath, "Capabilities")
	if err != nil {
		return err
	}
	defer k.Close()
	_ = k.SetStringValue("ApplicationDescription", "Switchy — choose which browser opens your links")
	_ = k.SetStringValue("ApplicationName", clientName)
	return nil
}

func writeURLAssociations(capPath string) error {
	k, err := createKey(capPath+`\URLAssociations`, "URLAssociations")
	if err != nil {
		return err
	}
	defer k.Close()
	_ = k.SetStringValue("http", progIDSuffix)
	_ = k.SetStringValue("https", progIDSuffix)
	return nil
}

func writeShellOpenCommand(base, exePath string) error {
	k, err := createKey(base+`\shell\open\command`, "shell open command")
	if err != nil {
		return err
	}
	defer k.Close()
	_ = k.SetStringValue("", fmt.Sprintf(`"%s" "%%1"`, exePath))
	return nil
}

func writeRegisteredApplication(capPath string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\RegisteredApplications`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("registry: open RegisteredApplications: %w", err)
	}
	defer k.Close()
	return k.SetStringValue(clientName, capPath)
}
