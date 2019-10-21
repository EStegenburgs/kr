// +build windows

package main

import (
	"fmt"
	"github.com/kryptco/kr"
	"github.com/pkg/browser"
	"github.com/urfave/cli"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
)

func initTerminal() {
	var m uint32
	windows.GetConsoleMode(windows.Stdout, &m)
	windows.SetConsoleMode(windows.Stdout, m|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func restartCommandOptions(c *cli.Context, isUserInitiated bool) (err error) {
	if isUserInitiated {
		kr.Analytics{}.PostEventUsingPersistedTrackingID("kr", "restart", nil, nil)
	}

	_ = migrateSSHConfig()

	kr.KillKrd()
	startKrd()

	if isUserInitiated {
		PrintErr(os.Stderr, "Restarted Krypton daemon.")
	}
	return
}

func upgradeCommand(c *cli.Context) (err error) {
	return fmt.Errorf("Upgrade not supported")
}

func uninstallCommand(c *cli.Context) (err error) {
	go func() {
		kr.Analytics{}.PostEventUsingPersistedTrackingID("kr", "uninstall", nil, nil)
	}()
	confirmOrFatal(os.Stderr, "Uninstall Krypton from this workstation?")

	cleanSSHConfig()

	kr.KillKrd()

	//uninstallCodesigning()
	PrintErr(os.Stderr, "Krypton uninstalled. If you experience any issues, please refer to https://krypt.co/docs/start/installation.html#uninstalling-kr")
	return
}

func startKrd() (err error) {
	cmd := exec.Command("cmd.exe", "/C", "start", "/b", `krd.exe`)
	return cmd.Run()
}

func openBrowser(url string) {
	err := browser.OpenURL(url)
	if err != nil {
		os.Stderr.WriteString("Unable to open browser, please visit " + url + "\r\n")
	}
}

func killKrd() {
	kr.KillKrd()
}
