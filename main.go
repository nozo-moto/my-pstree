package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Process struct {
	PID     int
	PPID    int
	Command string
}

func getProcesses() ([]Process, error) {
	// Read the /proc filesystem to get process information
	// This is a simplified example, and it assumes all processes are in /proc
	processes := []Process{}
	dirs, _ := os.ReadDir("/proc")
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(dir.Name())
		if err != nil {
			continue
			// ignore not pid in proc dir
		}
		if pid > 0 {
			cmdline, err := os.ReadFile("/proc/" + dir.Name() + "/comm")
			if err != nil {
			  return nil, err
			}
			cmd := string(cmdline[:len(cmdline)-1]) // Remove trailing null byte
			ppid, _ := readPPID("/proc/" + dir.Name() + "/stat")
			processes = append(processes, Process{PID: pid, PPID: ppid, Command: cmd})
		}
	}
	return processes, nil
}

func readPPID(statPath string) (int, error) {
	// Read the parent process ID from the /proc/[pid]/stat file
	stat, err := os.ReadFile(statPath)
	if err != nil {
		return 0, err
	}
	// The parent process ID is the fourth field in the stat file
	fields := splitFields(string(stat))
	ppid, _ := strconv.Atoi(fields[3])
	return ppid, nil
}

func splitFields(line string) []string {
	// Split a line into fields, handling quoted strings
	fields := []string{}
	field := ""
	inQuote := false
	for _, char := range line {
		switch char {
		case ' ':
			if inQuote {
				field += string(char)
			} else {
				fields = append(fields, field)
				field = ""
			}
		case '"':
			inQuote = !inQuote
			field += string(char)
		default:
			field += string(char)
		}
	}
	if field != "" {
		fields = append(fields, field)
	}
	return fields
}

func printProcessTree(processes []Process, ppid int, indent string) {
	// Recursive function to print the process tree
	for _, process := range processes {
		if process.PPID == ppid {
			fmt.Printf("%s%d %s\n", indent, process.PID, process.Command)
			printProcessTree(processes, process.PID, indent+"  ")
		}
	}
}

func main() {
	processes, err := getProcesses()
	if err != nil {
	  log.Fatal(err)
	}
	printProcessTree(processes, 0, "")
}
