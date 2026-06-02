package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/d3v-sh/d3m0n-agent/config"
	"github.com/d3v-sh/d3m0n-agent/logger"
	"github.com/d3v-sh/d3m0n-agent/tools"
	"github.com/d3v-sh/d3m0n-agent/ui"
	"github.com/fatih/color"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func trimHistory(messages []openai.ChatCompletionMessageParamUnion, max int) []openai.ChatCompletionMessageParamUnion {
	if max == 0 || len(messages) <= max {
		return messages
	}
	system := messages[0]
	trimmed := messages[len(messages)-(max-1):]
	return append([]openai.ChatCompletionMessageParamUnion{system}, trimmed...)
}

func printMarkdown(content string) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Println(content)
		return
	}
	out, err := renderer.Render(content)
	if err != nil {
		fmt.Println(content)
		return
	}
	fmt.Print(out)

}

func runAgent(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessageParamUnion, cfg *config.Config) []openai.ChatCompletionMessageParamUnion {
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond, spinner.WithWriter(os.Stderr))
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Timeout)*time.Second)
	defer cancel()
	for {
		params := openai.ChatCompletionNewParams{
			Model:    cfg.Model,
			Messages: messages,
			Tools:    tools.Tools,
		}

		s.Suffix = " Thinking...."
		s.Start()

		resp, err := client.Chat.Completions.New(ctx, params)
		s.Stop()

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Request time out.")
			} else {
				fmt.Println("API error: ", err)
			}
			return messages
		}
		choice := resp.Choices[0]

		// No tool call
		if choice.FinishReason == "stop" {
			ui.AssistantLabel()
			printMarkdown(choice.Message.Content)
			fullcontent, err := streamResponse(ctx, client, params)
			if err != nil {
				printMarkdown(choice.Message.Content)
				messages = append(messages, openai.AssistantMessage(choice.Message.Content))
			} else {
				messages = append(messages, openai.AssistantMessage(fullcontent))
			}
			messages = trimHistory(messages, cfg.MaxHistory)
			return messages
		}

		// Tool call
		if choice.FinishReason == "tool_calls" {
			messages = append(messages, choice.Message.ToParam())

			for _, toolCall := range choice.Message.ToolCalls {
				ui.PrintToolCall(toolCall.Function.Name)
				s.Suffix = fmt.Sprintf("Running: %s\n", toolCall.Function.Name)
				s.Start()
				result := tools.DispatchTool(toolCall.Function.Name, toolCall.Function.Arguments, cfg)
				s.Stop()
				ui.PrintToolResult(toolCall.Function.Name, result)
				messages = append(messages, openai.ToolMessage(toolCall.ID, result))
			}

		}
		// safeguard for free models
		if choice.FinishReason == "stop" && len(choice.Message.ToolCalls) == 0 {
			// model answered without calling any tools
			fmt.Println(color.YellowString("\n⚠ Note: This response was generated directly by the model without calling any tools."))
		}
	}
}
func streamResponse(ctx context.Context, client *openai.Client, params openai.ChatCompletionNewParams) (string, error) {
	stream := client.Chat.Completions.NewStreaming(ctx, params)

	var fullContent strings.Builder
	var displayed strings.Builder

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if delta != "" {
				fullContent.WriteString(delta)
				displayed.WriteString(delta)

				fmt.Print(color.YellowString(delta))
			}
		}
	}
	if err := stream.Err(); err != nil {
		return "", err
	}

	content := fullContent.String()
	lineCount := strings.Count(displayed.String(), "\n") + 2
	fmt.Printf("\033[%dA\033[J", lineCount)
	printMarkdown(content)
	return content, nil

}
func extractContent(message openai.ChatCompletionMessageParamUnion) string {
	// marshal to map and extract content field
	data, err := json.Marshal(message)
	if err != nil {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}
	if content, ok := m["content"].(string); ok {
		return content
	}
	return ""
}

func generateHTMLReport(markdownContent, outputPath string) error {
	glamour.NewTermRenderer(glamour.WithStylePath("nooty"))

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownContent), &buf); err != nil {
		return err
	}
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Penetration Testing Report</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; 
               max-width: 1000px; margin: 40px auto; padding: 0 20px; 
               background: #0d1117; color: #e6edf3; }
        h1 { color: #58a6ff; border-bottom: 2px solid #21262d; padding-bottom: 10px; }
        h2 { color: #79c0ff; border-bottom: 1px solid #21262d; padding-bottom: 6px; }
        h3 { color: #d2a8ff; }
        table { border-collapse: collapse; width: 100%%; margin: 16px 0; }
        th { background: #161b22; color: #58a6ff; padding: 10px; 
             border: 1px solid #30363d; text-align: left; }
        td { padding: 10px; border: 1px solid #30363d; }
        tr:nth-child(even) { background: #161b22; }
        code { background: #161b22; padding: 2px 6px; border-radius: 4px; 
               color: #ff7b72; font-family: monospace; }
        pre { background: #161b22; padding: 16px; border-radius: 8px; 
              border: 1px solid #30363d; overflow-x: auto; }
        pre code { color: #e6edf3; background: none; padding: 0; }
        strong { color: #ffa657; }
        .header { background: #161b22; padding: 20px; border-radius: 8px; 
                  border-left: 4px solid #58a6ff; margin-bottom: 30px; }
        hr { border: none; border-top: 1px solid #21262d; margin: 30px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔒 Penetration Testing Report</h1>
        <p>Generated by <strong>d3m0n</strong> on %s</p>
    </div>
    %s
</body>
</html>`, time.Now().Format("2006-01-02 15:04:05"), buf.String())

	return os.WriteFile(outputPath, []byte(html), 0644)
}

func main() {

	cfg, err := config.Load("config.yaml")

	ui.Banner("1.0.0", cfg.Model)
	ui.Prompt()

	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	if err := logger.Start(); err != nil {
		log.Fatalf("Failed to start logger: %v", err)
	}
	defer logger.End()
	ctx := context.Background()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if cfg.APIKey != "" {
		apiKey = cfg.APIKey
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(cfg.BaseURL),
	)

	reader := bufio.NewReader(os.Stdin)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(cfg.SystemPrompt),
	}

	fmt.Println("Cybersecurity Agent ready. Type 'exit' to quit, 'clear' to reset.")

	for {
		fmt.Print("\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}
		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}
		if input == "clear" {
			fmt.Print("\033[H\033[2J")
			messages = []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(cfg.SystemPrompt),
			}
			continue
		}
		if input == "/session" {
			ui.PrintConfig(cfg.Model, cfg.BaseURL, cfg.MaxHistory, cfg.Timeout, cfg.SafeMode, cfg.AllowedPaths)
			continue
		}
		if input == "/config" {
			ui.PrintConfig(cfg.Model, cfg.BaseURL, cfg.MaxHistory, cfg.Timeout, cfg.SafeMode, cfg.AllowedPaths)
			continue
		}
		if input == "/tools" {
			ui.PrintTools()
			continue
		}
		if input == "/help" {
			ui.PrintHelp()
			continue
		}
		if input == "/history" {
			rows, err := logger.DB().Query(`
        SELECT s.id, s.start_time, COUNT(t.id) as tool_calls
        FROM sessions s
        LEFT JOIN tool_logs t ON s.id = t.session_id
        GROUP BY s.id
        ORDER BY s.start_time DESC
        LIMIT 10
    `)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			defer rows.Close()
			fmt.Println("\nLast 10 sessions:")
			for rows.Next() {
				var id, start string
				var count int
				rows.Scan(&id, &start, &count)
				fmt.Printf("  %s | %s | %d tool calls\n", id, start, count)
			}
			continue
		}
		if input == "/targets" {
			rows, err := logger.DB().Query("SELECT host, info, updated FROM targets ORDER BY updated DESC")
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			defer rows.Close()
			fmt.Println("\nKnown targets:")
			for rows.Next() {
				var host, info, updated string
				rows.Scan(&host, &info, &updated)
				fmt.Printf("  %s | %s | %s\n", host, updated, info)
			}
			continue
		}
		if input == "/report" {
			messages = append(messages, openai.UserMessage(`
        Based on all tool outputs in this conversation, generate a professional 
        penetration testing report with:
        - Executive Summary
        - Scope  
        - Findings (Critical/High/Medium/Low)
        - Evidence
        - Recommendations
        - Conclusion
        Format in Markdown.
    `))
			messages = runAgent(ctx, &client, messages, cfg)

			content := extractContent(messages[len(messages)-1])
			timestamp := time.Now().Format("20060102-150405")
			os.MkdirAll("output", 0755)

			mdPath := fmt.Sprintf("output/report-%s.md", timestamp)
			os.WriteFile(mdPath, []byte(content), 0644)

			// save html
			htmlPath := fmt.Sprintf("output/report-%s.html", timestamp)
			if err := generateHTMLReport(content, htmlPath); err != nil {
				fmt.Println(color.RedString("HTML generation failed: %v", err))
			} else {
				fmt.Printf(color.GreenString("\nReports saved:\n  MD:   %s\n  HTML: %s\n"), mdPath, htmlPath)
			}

		}

		messages = append(messages, openai.UserMessage(input))
		messages = runAgent(ctx, &client, messages, cfg)

	}

}
