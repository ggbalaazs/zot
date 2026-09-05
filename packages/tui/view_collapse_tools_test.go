package tui

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/patriceckhart/zot/packages/provider"
)

func TestCollapseToolCallShowsLastBoxPreviewLine(t *testing.T) {
	args := json.RawMessage(`{"command":"echo preview"}`)
	v := View{Theme: Dark, CollapseToolCall: true, ToolCalls: []ToolCallView{{
		ID: "toolu_1", Name: "bash", Args: ShortArgs("bash", args), Result: "first-output\nlast-output\n",
	}}}
	plain := stripANSI(strings.Join(v.BuildLive(80), "\n"))
	if !strings.Contains(plain, "bash") {
		t.Fatalf("collapsed render lost tool header:\n%s", plain)
	}
	if !strings.Contains(plain, "last-output") || strings.Contains(plain, "first-output") {
		t.Fatalf("collapsed render did not show only the last result line:\n%s", plain)
	}
	if !strings.ContainsAny(plain, "┌┐") || !strings.ContainsAny(plain, "└┘") || !strings.Contains(plain, "│") {
		t.Fatalf("collapsed box render did not retain the complete empty frame:\n%s", plain)
	}
	v.ExpandAll = true
	expanded := stripANSI(strings.Join(v.BuildLive(80), "\n"))
	if !strings.Contains(expanded, "first-output") {
		t.Fatalf("ctrl-o expansion did not reveal non-compact collapsed tool result:\n%s", expanded)
	}
}

func TestCollapseToolCallShowsLastCompactPreviewLine(t *testing.T) {
	args := json.RawMessage(`{"command":"echo preview"}`)
	v := View{
		Theme:            Dark,
		CompactMode:      true,
		CollapseToolCall: true,
		Messages: []provider.Message{
			{Role: provider.RoleAssistant, Content: []provider.Content{
				provider.ToolCallBlock{ID: "toolu_1", Name: "bash", Arguments: args},
			}},
			{Role: provider.RoleTool, Content: []provider.Content{
				provider.ToolResultBlock{CallID: "toolu_1", Content: []provider.Content{provider.TextBlock{Text: "first-output\nlast-output"}}},
			}},
		},
	}
	plain := stripANSI(strings.Join(v.Build(80), "\n"))
	if !strings.Contains(plain, "bash") {
		t.Fatalf("collapsed compact render lost tool header:\n%s", plain)
	}
	if !strings.Contains(plain, "last-output") || strings.Contains(plain, "first-output") {
		t.Fatalf("collapsed compact render did not show only the last result line:\n%s", plain)
	}
	v.ExpandAll = true
	expanded := stripANSI(strings.Join(v.Build(80), "\n"))
	if !strings.Contains(expanded, "first-output") {
		t.Fatalf("ctrl-o expansion did not reveal collapsed tool result:\n%s", expanded)
	}
}

func TestCollapseToolCallShowsLastLivePreviewLine(t *testing.T) {
	args := json.RawMessage(`{"command":"echo preview"}`)
	v := View{Theme: Dark, CollapseToolCall: true, ToolCalls: []ToolCallView{{
		ID: "toolu_1", Name: "bash", Args: ShortArgs("bash", args), RawJSONBuf: string(args),
	}}}
	plain := stripANSI(strings.Join(v.BuildLive(80), "\n"))
	if !strings.Contains(plain, "$ echo") {
		t.Fatalf("collapsed live render lost its one-line preview:\n%s", plain)
	}
	if !strings.Contains(plain, "bash") {
		t.Fatalf("collapsed live render lost tool header:\n%s", plain)
	}
}
