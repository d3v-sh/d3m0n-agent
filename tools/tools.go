package tools

import "github.com/openai/openai-go/v3"

var Tools = []openai.ChatCompletionToolUnionParam{
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_nmap",
		Description: openai.String("Runs an nmap scan on a target IP or domain"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{"type": "string", "description": "IP or hostname"},
				"flags":  map[string]any{"type": "string", "description": "nmap flags e.g. -sV -p 80,443"},
			},
			"required": []string{"target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_whois",
		Description: openai.String("Run whois lookup on a domain or IP"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{"type": "string"},
			},
			"required": []string{"target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "read_file",
		Description: openai.String("Reads a file from disk and returns its contents"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Absolute or relative file path",
				},
			},
			"required": []string{"path"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "write_file",
		Description: openai.String("Create a file or append data to the file"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Absolute or relative path",
				},
				"data": map[string]any{
					"type":        "string",
					"description": "Content to write to the file",
				},
				"append": map[string]any{
					"type":        "boolean",
					"description": "If true, appends to file instead of overwriting",
				},
			},
			"required": []string{"path"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_curl",
		Description: openai.String("Runs curl to probe HTTP endpoints, inspect headers, or fetch content"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{
					"type":        "string",
					"description": "Target URL",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "curl flags e.g. -I for headers, -L to follow redirects",
				},
			},
			"required": []string{"url"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_gobuster",
		Description: openai.String("Runs gobuster for directory or subdomain brute forcing"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"mode": map[string]any{
					"type":        "string",
					"description": "Mode: dir dns or vhost",
				},
				"target": map[string]any{
					"type":        "string",
					"description": "Target URL or domain",
				},
				"wordlist": map[string]any{
					"type":        "string",
					"description": "Path to wordlist file",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra gobuster flags e.g. -t 50 -x php,html",
				},
			},

			"required": []string{"mode", "target", "wordlist"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_amass",
		Description: openai.String("Runs amass for DNS enumeration and subdomain discovery"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"domain": map[string]any{
					"type":        "string",
					"description": "Target domain",
				},
				"mode": map[string]any{
					"type":        "string",
					"description": "Mode: enum, intel, viz, track, db",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. -passive -brute",
				},
			},
			"required": []string{"domain", "mode"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_recon_ng",
		Description: openai.String("Runs recon-ng for web reconnaissance"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"module": map[string]any{
					"type":        "string",
					"description": "Module to run e.g. recon/domains-hosts/hackertarget",
				},
				"target": map[string]any{
					"type":        "string",
					"description": "Target domain or IP",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags",
				},
			},
			"required": []string{"module", "target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_theharvester",
		Description: openai.String("Runs theHarvester for email, subdomain, and name gathering from public sources"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"domain": map[string]any{
					"type":        "string",
					"description": "Target domain",
				},
				"source": map[string]any{
					"type":        "string",
					"description": "Data source e.g. google, bing, linkedin, all",
				},
				"limit": map[string]any{
					"type":        "string",
					"description": "Number of results to retrieve",
				},
			},
			"required": []string{"domain", "source"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_sherlock",
		Description: openai.String("Runs sherlock to find social media accounts by username"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"username": map[string]any{
					"type":        "string",
					"description": "Username to search across platforms",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. --timeout 10 --print-found",
				},
			},
			"required": []string{"username"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_spiderfoot",
		Description: openai.String("Runs spiderfoot for automated OSINT collection"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "Target domain, IP, email, or username",
				},
				"modules": map[string]any{
					"type":        "string",
					"description": "Comma separated modules to run e.g. sfp_dns,sfp_whois",
				},
				"output": map[string]any{
					"type":        "string",
					"description": "Output file path",
				},
			},
			"required": []string{"target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_eyewitness",
		Description: openai.String("Runs eyewitness to take screenshots of web services"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "URL or file containing list of URLs",
				},
				"output": map[string]any{
					"type":        "string",
					"description": "Output directory path",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. --web --timeout 30",
				},
			},
			"required": []string{"target", "output"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_ffuf",
		Description: openai.String("Runs ffuf for fast web fuzzing of directories, parameters, and vhosts"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{
					"type":        "string",
					"description": "Target URL with FUZZ keyword e.g. https://target.com/FUZZ",
				},
				"wordlist": map[string]any{
					"type":        "string",
					"description": "Path to wordlist",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. -mc 200 -t 50 -e .php,.html",
				},
			},
			"required": []string{"url", "wordlist"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_crtsh",
		Description: openai.String("Queries crt.sh certificate transparency logs for subdomain discovery"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"domain": map[string]any{
					"type":        "string",
					"description": "Target domain to query",
				},
			},
			"required": []string{"domain"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_sslscan",
		Description: openai.String("Runs sslscan to test SSL/TLS configuration of a target"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "Target host:port e.g. example.com:443",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. --no-colour --bugs",
				},
			},
			"required": []string{"target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_testssl",
		Description: openai.String("Runs testssl.sh for thorough SSL/TLS testing"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "Target host:port e.g. example.com:443",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. --severity HIGH --json",
				},
			},
			"required": []string{"target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_gitleaks",
		Description: openai.String("Runs gitleaks to detect secrets and sensitive data in git repos or directories"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to git repo or directory to scan",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. --report-format json --report-path report.json",
				},
			},
			"required": []string{"path"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_sqlmap",
		Description: openai.String("Runs sqlmap for automated SQL injection detection and exploitation"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "Target URL e.g. http://target.com/page?id=1",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags e.g. --dbs --level 3 --risk 2 --batch",
				},
			},
			"required": []string{"target"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "run_dradis",
		Description: openai.String("Runs dradis CLI to import scan results into dradis reporting framework"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "Dradis command e.g. import, export, create",
				},
				"file": map[string]any{
					"type":        "string",
					"description": "File to import or export",
				},
				"flags": map[string]any{
					"type":        "string",
					"description": "Extra flags",
				},
			},
			"required": []string{"command"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "search_cve",
		Description: openai.String("Search NVD database for CVEs affecting a product or version"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"product": map[string]any{"type": "string"},
				"version": map[string]any{"type": "string"},
			},
			"required": []string{"product"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "remember",
		Description: openai.String("Save important information about a target for future reference"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{"type": "string"},
				"info":   map[string]any{"type": "string"},
			},
			"required": []string{"target", "info"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "recall",
		Description: openai.String("Recall saved information about a target"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{"type": "string"},
			},
			"required": []string{"target"},
		},
	}),
}
