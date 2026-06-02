package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/d3v-sh/sec_agent/config"
	"github.com/d3v-sh/sec_agent/logger"
)

func RememberTarget(target, info string) {
	logger.DB().Exec("INSERT OR REPLACE INTO targets (host, info, updated) VALUES (?, ?, ?)",
		target, info, time.Now())
}

// agent can query its own memory
func RecallTarget(target string) string {
	var info string
	logger.DB().QueryRow("SELECT info FROM targets WHERE host = ?", target).Scan(&info)
	return info
}

func runNmap(target, flags string) string {
	args := strings.Fields(flags)
	args = append(args, target)
	out, err := exec.Command("nmap", args...).CombinedOutput()
	if err != nil {
		return "nmap error: " + err.Error()
	}
	return string(out)
}

func runCurl(url, flags string) string {
	args := strings.Fields(flags)
	args = append(args, url)
	out, err := exec.Command("curl", args...).CombinedOutput()
	if err != nil {
		return "curl error: " + err.Error()
	}
	return string(out)
}
func runGobuster(mode, target, wordlist, flags string) string {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	args := []string{mode, "-u", target, "-w", wordlist}
	if flags != "" {
		args = append(args, strings.Fields(flags)...)
	}

	out, err := exec.Command("gobuster", args...).CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "gobuster timed out after 5 minutes\nPartial output:\n" + string(out)
		}
		return "gobuster error: " + err.Error()
	}
	return string(out)
}

func runWhois(target string) string {
	out, err := exec.Command("whois", target).CombinedOutput()
	if err != nil {
		return "whois error: " + err.Error()
	}
	return string(out)
}
func runCVESearch(product, version string) string {
	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=%s+%s", product, version)
	out, err := exec.Command("curl", "-s", url).CombinedOutput()
	if err != nil {
		return "CVE search error: " + err.Error()
	}
	return string(out)
}

func DispatchTool(name, args string, cfg *config.Config) string {
	result := dispatch(name, args, cfg)
	logger.LogTool(name, args, result)
	return result
}

// Helper functions
func dispatch(name, args string, cfg *config.Config) string {
	switch name {
	case "search_cve":
		var a struct {
			Product string `json:"product"`
			Version string `json:"version"`
		}
		json.Unmarshal([]byte(args), &a)
		return runCVESearch(a.Product, a.Version)
	case "run_nmap":
		var a struct {
			Target string `json:"target"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		return runNmap(a.Target, a.Flags)

	case "run_whois":
		var a struct {
			Target string `json:"target"`
		}
		json.Unmarshal([]byte(args), &a)
		return runWhois(a.Target)

	case "read_file":
		var a struct {
			Path string `json:"path"`
		}
		json.Unmarshal([]byte(args), &a)

		fmt.Printf("[debug] reading path: %s\n", a.Path)

		if cfg.SafeMode && !safePath(a.Path, cfg.AllowedPaths) {
			return fmt.Sprintf("unsafe path: allowed paths are %v", cfg.AllowedPaths)
		}

		content, err := os.ReadFile(a.Path)
		if err != nil {
			return "error: " + err.Error()
		}
		return string(content)
	case "write_file":
		var a struct {
			Path   string `json:"path"`
			Data   string `json:"data"`
			Append bool   `json:"append"`
		}
		json.Unmarshal([]byte(args), &a)
		if cfg.SafeMode && !safePath(a.Path, cfg.AllowedPaths) {
			cwd, _ := os.Getwd()
			return fmt.Sprintf("unsafe path: allowed paths are %v or %s", cfg.AllowedPaths, cwd)
		}
		if safePath(a.Path, cfg.AllowedPaths) {
			if a.Append {
				f, err := os.OpenFile(a.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					return "error: " + err.Error()
				}
				defer f.Close()
				f.WriteString(a.Data)
			} else {
				if err := os.WriteFile(a.Path, []byte(a.Data), 0666); err != nil {
					return "error: " + err.Error()
				}
			}
			return "Data written successfully"
		} else {
			fmt.Println("unsafe path")
			return "unsafe path"
		}

	case "run_curl":
		var a struct {
			URL   string `json:"url"`
			Flags string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		return runCurl(a.URL, a.Flags)
	case "run_gobuster":
		var a struct {
			Mode     string `json:"mode"`
			Target   string `json:"target"`
			Wordlist string `json:"wordlist"`
			Flags    string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		return runGobuster(a.Mode, a.Target, a.Wordlist, a.Flags)
	case "run_amass":
		var a struct {
			Domain string `json:"domain"`
			Mode   string `json:"mode"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{a.Mode, "-d", a.Domain}
		if a.Flags != "" {
			arguments = append(arguments, strings.Fields(a.Flags)...)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, "amass", arguments...).CombinedOutput()
		if err != nil {
			return "amass error: " + err.Error()
		}
		return string(out)

	case "run_recon_ng":
		var a struct {
			Module string `json:"module"`
			Target string `json:"target"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		script := fmt.Sprintf("modules load %s\noptions set SOURCE %s\nrun\nexit\n", a.Module, a.Target)
		cmd := exec.Command("recon-ng", "--no-check")
		cmd.Stdin = strings.NewReader(script)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "recon-ng error: " + err.Error()
		}
		return string(out)

	case "run_theharvester":
		var a struct {
			Domain string `json:"domain"`
			Source string `json:"source"`
			Limit  string `json:"limit"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{"-d", a.Domain, "-b", a.Source}
		if a.Limit != "" {
			arguments = append(arguments, "-l", a.Limit)
		}
		out, err := exec.Command("theHarvester", arguments...).CombinedOutput()
		if err != nil {
			return "theHarvester error: " + err.Error()
		}
		return string(out)

	case "run_sherlock":
		var a struct {
			Username string `json:"username"`
			Flags    string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := strings.Fields(a.Flags)
		arguments = append(arguments, a.Username)
		out, err := exec.Command("sherlock", arguments...).CombinedOutput()
		if err != nil {
			return "sherlock error: " + err.Error()
		}
		return string(out)

	case "run_spiderfoot":
		var a struct {
			Target  string `json:"target"`
			Modules string `json:"modules"`
			Output  string `json:"output"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{"-s", a.Target, "-q"}
		if a.Modules != "" {
			arguments = append(arguments, "-m", a.Modules)
		}
		if a.Output != "" {
			arguments = append(arguments, "-o", a.Output)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, "amass", arguments...).CombinedOutput()
		if err != nil {
			return "spiderfoot error: " + err.Error()
		}
		return string(out)

	case "run_eyewitness":
		var a struct {
			Target string `json:"target"`
			Output string `json:"output"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{"--web", "-u", a.Target, "-d", a.Output}
		if a.Flags != "" {
			arguments = append(arguments, strings.Fields(a.Flags)...)
		}
		out, err := exec.Command("eyewitness", arguments...).CombinedOutput()
		if err != nil {
			return "eyewitness error: " + err.Error()
		}
		return string(out)

	case "run_ffuf":
		var a struct {
			URL      string `json:"url"`
			Wordlist string `json:"wordlist"`
			Flags    string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{"-u", a.URL, "-w", a.Wordlist}
		if a.Flags != "" {
			arguments = append(arguments, strings.Fields(a.Flags)...)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, "ffuf", arguments...).CombinedOutput()
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return "ffuf timed out\nPartial output:\n" + string(out)
			}
			return "ffuf error: " + err.Error()
		}
		return string(out)

	case "run_crtsh":
		var a struct {
			Domain string `json:"domain"`
		}
		json.Unmarshal([]byte(args), &a)
		url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", a.Domain)
		out, err := exec.Command("curl", "-s", url).CombinedOutput()
		if err != nil {
			return "crt.sh error: " + err.Error()
		}
		return string(out)

	case "run_sslscan":
		var a struct {
			Target string `json:"target"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := strings.Fields(a.Flags)
		arguments = append(arguments, a.Target)
		out, err := exec.Command("sslscan", arguments...).CombinedOutput()
		if err != nil {
			return "sslscan error: " + err.Error()
		}
		return string(out)

	case "run_testssl":
		var a struct {
			Target string `json:"target"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := strings.Fields(a.Flags)
		arguments = append(arguments, a.Target)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, "testssl.sh", arguments...).CombinedOutput()
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return "testssl timed out\nPartial output:\n" + string(out)
			}
			return "testssl error: " + err.Error()
		}
		return string(out)

	case "run_gitleaks":
		var a struct {
			Path  string `json:"path"`
			Flags string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{"detect", "--source", a.Path}
		if a.Flags != "" {
			arguments = append(arguments, strings.Fields(a.Flags)...)
		}
		out, err := exec.Command("gitleaks", arguments...).CombinedOutput()
		if err != nil {
			return "gitleaks error: " + err.Error()
		}
		return string(out)

	case "run_sqlmap":
		var a struct {
			Target string `json:"target"`
			Flags  string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{"-u", a.Target, "--batch"}
		if a.Flags != "" {
			arguments = append(arguments, strings.Fields(a.Flags)...)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, "sqlmap", arguments...).CombinedOutput()
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return "sqlmap timed out\nPartial output:\n" + string(out)
			}
			return "sqlmap error: " + err.Error()
		}
		return string(out)
	case "remember":
		var a struct {
			Target string `json:"target"`
			Info   string `json:"info"`
		}
		json.Unmarshal([]byte(args), &a)
		RememberTarget(a.Target, a.Info)
		return fmt.Sprintf("Remembered: %s -> %s", a.Target, a.Info)
	case "recall":
		var a struct {
			Target string `json:"target"`
		}
		json.Unmarshal([]byte(args), &a)
		info := RecallTarget(a.Target)
		if info == "" {
			return fmt.Sprintf("No information saved for %s", a.Target)
		}
		return fmt.Sprintf("Known info about %s: %s", a.Target, info)
	case "run_dradis":
		var a struct {
			Command string `json:"command"`
			File    string `json:"file"`
			Flags   string `json:"flags"`
		}
		json.Unmarshal([]byte(args), &a)
		arguments := []string{a.Command}
		if a.File != "" {
			arguments = append(arguments, a.File)
		}
		if a.Flags != "" {
			arguments = append(arguments, strings.Fields(a.Flags)...)
		}
		out, err := exec.Command("dradis", arguments...).CombinedOutput()
		if err != nil {
			return "dradis error: " + err.Error()
		}
		return string(out)
	default:
		return "unknown tool: " + name
	}
}

func safePath(path string, allowedPaths []string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	for _, allowed := range allowedPaths {
		allowedAbs, err := filepath.Abs(allowed)
		if err != nil {
			continue
		}
		if strings.HasPrefix(abs, allowedAbs) {
			return true
		}
	}
	return false
}
