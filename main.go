package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func setupLogger() *os.File {
	logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Unable to open log file:", err)
	}
	log.SetOutput(logFile)
	return logFile
}

func out() *log.Logger {
	Stdlog := log.New(os.Stderr)
	log.ErrorLevelStyle = lipgloss.NewStyle().
		SetString("ERROR!!").
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.AdaptiveColor{
			Light: "203",
			Dark:  "204",
		}).
		Foreground(lipgloss.Color("#211f26"))
	return Stdlog
}

func closeLogger(logFile *os.File) {
	err := logFile.Close()
	if err != nil {
		log.Fatal("Error closing log file:", err)
	}
}

func checkerr(err error, Stdlog *log.Logger) {
	if err != nil {
		Stdlog.Error(err)
		log.Fatal(err)
		return
	}
}

var welcomeStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#81efc5")).
	Padding(0, 1, 0, 1).
	BorderBottom(true).
	BorderTop(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottomForeground(lipgloss.Color("#3c4056")).
	BorderTopForeground(lipgloss.Color("#3c4056"))

var InputPromptStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#6aaa96"))

var arrStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ba9af8"))

var BorderNotif = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderTop(true).
	BorderForeground(lipgloss.Color("#3c4056")).
	Padding(0, 1, 0, 1)

var infoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#f19da5"))

func main() {
	Stdlog := out()
	logFile := setupLogger()
	defer closeLogger(logFile)
	fmt.Println(welcomeStyle.Render("Welcome to Organizer !!"))
	fmt.Println(InputPromptStyle.Render("Enter the directory to clean"))
	fmt.Print(arrStyle.Render("> "))

	in := bufio.NewReader(os.Stdin)

	filePathDir, err := in.ReadString('\n')
	checkerr(err, Stdlog)
	filePathDir = strings.TrimSuffix(filePathDir, "\n")
	filePathDir = strings.TrimSuffix(filePathDir, "\r")

	files, err := os.ReadDir(filePathDir)
	checkerr(err, Stdlog)
	fmt.Println(BorderNotif.Render("Below files are found"))

	for _, file := range files {
		fmt.Println(infoStyle.Render(fmt.Sprintf("* %s", file.Name())))
		extensionName := filepath.Ext(file.Name())
		extensionName = strings.TrimPrefix(extensionName, ".")

		if file.Name() == ".DS_Store" || extensionName == "" {
			continue
		}
		destinationDir := filepath.Join(filePathDir, extensionName)
		err := os.MkdirAll(destinationDir, 0755)
		checkerr(err, Stdlog)

		sourceFilePath := filepath.Join(filePathDir, file.Name())
		destinationFilePath := filepath.Join(destinationDir, file.Name())
		err = os.Rename(sourceFilePath, destinationFilePath)
		checkerr(err, Stdlog)
	}
	log.Info(filePathDir)

	fmt.Println(InputPromptStyle.Render("Press Enter key to close"))
	_, err = in.ReadString('\n')
	checkerr(err, Stdlog)
}
