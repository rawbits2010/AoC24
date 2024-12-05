package outputhandler

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type RunningEnvironment struct {
	Name                 string
	ExeName              string
	AddsCSICursorSupport bool
	AddsCSIColorSupport  bool
	AddsEmojiSupport     bool
}

var knowEnvironments = []RunningEnvironment{
	{
		Name:                 "File Explorer",
		ExeName:              "explorer.exe",
		AddsCSICursorSupport: false,
		AddsCSIColorSupport:  false,
		AddsEmojiSupport:     false,
	},
	{
		Name:                 "VS Code",
		ExeName:              "Code.exe",
		AddsCSICursorSupport: true,
		AddsCSIColorSupport:  true,
		AddsEmojiSupport:     true,
	},
}

type TerminalInfo struct {
	Name             string
	ExeName          string
	CSICursorSupport bool
	CSIColorSupport  bool
	EmojiSupport     bool
}

var knownTerminals = []TerminalInfo{
	{ // Go debugger
		Name:             "Delve",
		ExeName:          "dlv.exe",
		CSICursorSupport: false,
		CSIColorSupport:  true,
		EmojiSupport:     true,
	},
	{ // I use Git bash so...
		Name:             "bash",
		ExeName:          "bash.exe",
		CSICursorSupport: true,
		CSIColorSupport:  true,
		EmojiSupport:     true,
	},
	{
		Name:             "Command Prompt",
		ExeName:          "cmd.exe",
		CSICursorSupport: true,
		CSIColorSupport:  true,
		EmojiSupport:     false,
	},
	{
		Name:             "PowerShell",
		ExeName:          "powershell.exe",
		CSICursorSupport: true,
		CSIColorSupport:  true,
		EmojiSupport:     false,
	},
	{ // Win11 thing, no clue about this one
		Name:             "Windows Terminal",
		ExeName:          "wt.exe",
		CSICursorSupport: true,
		CSIColorSupport:  true,
		EmojiSupport:     true,
	},
}

// GetTerminalInfo detects the terminal and the runner of the terminal
// to provide some hand tested info on features.
//
// NOTE: This method is not even close to accurate, but it's good enough
// for what it's used for here. :)
func GetTerminalInfo() (*TerminalInfo, *RunningEnvironment, error) {

	processes, err := GetProcesses()
	if err != nil {
		return nil, nil, fmt.Errorf("error enumerating processes: %w", err)
	}
	currPID := uint32(os.Getppid())

	var terminal TerminalInfo
	var terminalFound = false
	var env RunningEnvironment
	var envFound bool = false
	for {
		if procInfo, ok := (*processes)[currPID]; ok {
			//fmt.Printf("%s\n", procInfo.ExeName)

			if !terminalFound {
				for _, ti := range knownTerminals {
					if ti.ExeName == procInfo.ExeName {
						terminal = ti
						terminalFound = true
						break
					}
				}
			}

			if !envFound {
				for _, e := range knowEnvironments {
					if e.ExeName == procInfo.ExeName {
						env = e
						envFound = true
						break
					}
				}
			}

			currPID = procInfo.ParentPID
		} else {
			break
		}
	}

	//fmt.Printf("T:%s,E:%s", terminal.Name, env.Name)
	return &terminal, &env, nil
}

type ProcessInfo struct {
	ExeName   string
	ProcessId uint32
	ParentPID uint32
}

// GetProcesses enumerates all running processes - at least browsing MSDN gives the impression.
func GetProcesses() (*map[uint32]ProcessInfo, error) {

	hSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("error in CreateToolhelp32Snapshot: %w", err)
	}
	defer windows.CloseHandle(hSnapshot)

	pe := windows.ProcessEntry32{}
	pe.Size = uint32(unsafe.Sizeof(pe))

	processInfo := make(map[uint32]ProcessInfo, 0)
	for {

		exeName := windows.UTF16ToString(pe.ExeFile[:])
		processInfo[pe.ProcessID] = ProcessInfo{
			ExeName:   exeName,
			ProcessId: pe.ProcessID,
			ParentPID: pe.ParentProcessID,
		}

		err := windows.Process32Next(hSnapshot, &pe)
		if err == windows.ERROR_NO_MORE_FILES {
			return &processInfo, nil
		} else if err != nil {
			return nil, fmt.Errorf("error in Process32Next: %w", err)
		}
	}
}
