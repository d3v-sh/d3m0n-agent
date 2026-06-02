package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	// colors
	primary   = lipgloss.Color("#00ff99")
	secondary = lipgloss.Color("#0099ff")
	warning   = lipgloss.Color("#ffaa00")
	danger    = lipgloss.Color("#ff4444")
	muted     = lipgloss.Color("#666666")
	white     = lipgloss.Color("#ffffff")

	// styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Padding(0, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(secondary).
			Bold(true)

	toolStyle = lipgloss.NewStyle().
			Foreground(warning).
			Bold(true).
			PaddingLeft(2)

	resultStyle = lipgloss.NewStyle().
			Foreground(muted).
			PaddingLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(danger).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(muted)

	infoStyle = lipgloss.NewStyle().
			Foreground(secondary).
			PaddingLeft(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(0, 1)
)

func Banner(version, model string) {
	banner := `██████╗ ██████╗ ███╗   ███╗ ██████╗ ███╗   ██╗
 ██╔══██╗╚════██╗████╗ ████║██╔═████╗████╗  ██║
 ██║  ██║ █████╔╝██╔████╔██║██║██╔██║██╔██╗ ██║
 ██║  ██║ ╚═══██╗██║╚██╔╝██║████╔╝██║██║╚██╗██║
 ██████╔╝██████╔╝██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
 ╚═════╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
          [ AI-Powered Cybersecurity Agent ]
               [ by d3v-sh | v1.0.0 ]`

	fmt.Println(lipgloss.NewStyle().Foreground(primary).Bold(true).Render(banner))

	info := lipgloss.JoinHorizontal(lipgloss.Left,
		infoStyle.Render("Version: "+version),
		infoStyle.Foreground(muted).Render("  │  "),
		infoStyle.Render("Model: "+model),
		infoStyle.Foreground(muted).Render("  │  "),
		infoStyle.Render("Time: "+time.Now().Format("2006-01-02 15:04:05")),
	)
	fmt.Println(info)
	fmt.Println(separatorStyle.Render(strings.Repeat("─", 80)))
	fmt.Println()
}

func PrintHelp() {
	commands := [][]string{
		{"/help", "Show this message"},
		{"/tools", "List available tools"},
		{"/session", "Show current session info"},
		{"/history", "Show past sessions"},
		{"/targets", "Show remembered targets"},
		{"/config", "Show current configuration"},
		{"/report", "Generate pentest report"},
		{"clear", "Clear screen and reset history"},
		{"exit", "Quit d3m0n"},
	}

	title := headerStyle.Render("Available Commands")
	fmt.Println(title)
	fmt.Println(separatorStyle.Render(strings.Repeat("─", 40)))

	for _, cmd := range commands {
		key := lipgloss.NewStyle().Foreground(primary).Bold(true).Width(12).Render(cmd[0])
		desc := lipgloss.NewStyle().Foreground(white).Render(cmd[1])
		fmt.Printf("  %s  %s\n", key, desc)
	}
	fmt.Println()
}

func PrintTools() {
	tools := [][]string{
		{"run_nmap", "Port scanning and service detection"},
		{"run_whois", "Domain/IP registration lookup"},
		{"run_curl", "HTTP probing and header inspection"},
		{"run_gobuster", "Directory and subdomain brute forcing"},
		{"run_ffuf", "Fast web fuzzing"},
		{"run_amass", "Subdomain enumeration"},
		{"run_theharvester", "OSINT gathering"},
		{"run_sherlock", "Username search across platforms"},
		{"run_sslscan", "SSL/TLS configuration testing"},
		{"run_testssl", "Thorough SSL/TLS testing"},
		{"run_gitleaks", "Secret detection in repos"},
		{"run_sqlmap", "SQL injection testing"},
		{"run_crtsh", "Certificate transparency lookup"},
		{"run_amass", "DNS enumeration"},
		{"read_file", "Read file from disk"},
		{"write_file", "Write file to disk"},
		{"remember", "Save target information"},
		{"recall", "Recall saved target info"},
	}

	title := headerStyle.Render("Available Tools")
	fmt.Println(title)
	fmt.Println(separatorStyle.Render(strings.Repeat("─", 50)))

	for _, tool := range tools {
		key := lipgloss.NewStyle().Foreground(warning).Bold(true).Width(20).Render(tool[0])
		desc := lipgloss.NewStyle().Foreground(white).Render(tool[1])
		fmt.Printf("  %s  %s\n", key, desc)
	}
	fmt.Println()
}

func PrintConfig(model, baseURL string, maxHistory, timeout int, safeMode bool, allowedPaths []string) {
	content := fmt.Sprintf(
		"%s %s\n%s %s\n%s %d\n%s %ds\n%s %v\n%s %v",
		lipgloss.NewStyle().Foreground(muted).Render("Model:        "), lipgloss.NewStyle().Foreground(white).Render(model),
		lipgloss.NewStyle().Foreground(muted).Render("Base URL:     "), lipgloss.NewStyle().Foreground(white).Render(baseURL),
		lipgloss.NewStyle().Foreground(muted).Render("Max History:  "), maxHistory,
		lipgloss.NewStyle().Foreground(muted).Render("Timeout:      "), timeout,
		lipgloss.NewStyle().Foreground(muted).Render("Safe Mode:    "), safeMode,
		lipgloss.NewStyle().Foreground(muted).Render("Allowed Paths:"), allowedPaths,
	)
	fmt.Println(boxStyle.Render(content))
	fmt.Println()
}

func PrintToolCall(toolName string) {
	icon := "⚙"
	msg := fmt.Sprintf("%s  calling → %s", icon, toolName)
	fmt.Println(toolStyle.Render(msg))
}

func PrintToolResult(toolName, result string) {
	// truncate long results in display
	display := result
	if len(display) > 200 {
		display = display[:200] + "..."
	}
	_ = display
	fmt.Println(resultStyle.Render(fmt.Sprintf("✓  %s completed", toolName)))
}

func PrintError(msg string) {
	fmt.Println(errorStyle.Render("✗  " + msg))
}

func PrintWarning(msg string) {
	fmt.Println(lipgloss.NewStyle().Foreground(warning).Bold(true).Render("⚠  " + msg))
}

func PrintSuccess(msg string) {
	fmt.Println(lipgloss.NewStyle().Foreground(primary).Bold(true).Render("✓  " + msg))
}

func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render("ℹ  " + msg))
}

func Prompt() {
	fmt.Print(promptStyle.Render("\n❯ "))
}

func Separator() {
	fmt.Println(separatorStyle.Render(strings.Repeat("─", 80)))
}

func PrintSession(id string, toolCalls int) {
	content := fmt.Sprintf(
		"%s %s\n%s %d",
		lipgloss.NewStyle().Foreground(muted).Render("Session ID:  "), lipgloss.NewStyle().Foreground(white).Render(id),
		lipgloss.NewStyle().Foreground(muted).Render("Tool Calls:  "), toolCalls,
	)
	fmt.Println(boxStyle.Render(content))
}

func AssistantLabel() {
	label := lipgloss.NewStyle().
		Foreground(secondary).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(secondary).
		PaddingLeft(1).
		Render("Assistant")
	fmt.Println("\n" + label)
}

// for color package compatibility
func RedString(format string, a ...any) string {
	return color.New(color.FgRed).Sprintf(format, a...)
}
