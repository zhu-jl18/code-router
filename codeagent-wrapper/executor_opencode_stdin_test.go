package main

import (
	"context"
	"testing"
)

func TestRunCodexTaskWithContext_Opencode_DoesNotUseStdinForPrompt(t *testing.T) {
	defer resetTestHooks()

	task := "line1\nline2"
	opencodeOut := `{"type":"text","timestamp":1,"sessionID":"ses_ok","part":{"type":"text","text":"OK"}}` + "\n" +
		`{"type":"step_finish","timestamp":2,"sessionID":"ses_ok","part":{"type":"step-finish","reason":"stop"}}` + "\n"

	fake := newFakeCmd(fakeCmdConfig{
		StdoutPlan: []fakeStdoutEvent{{Data: opencodeOut}},
	})

	var gotName string
	var gotArgs []string

	newCommandRunner = func(ctx context.Context, name string, args ...string) commandRunner {
		gotName = name
		gotArgs = append([]string(nil), args...)
		return fake
	}

	res := runCodexTaskWithContext(
		context.Background(),
		TaskSpec{ID: "opencode", Task: task, WorkDir: ".", Mode: "new", UseStdin: true},
		OpencodeBackend{},
		nil,
		false,
		true,
		5,
	)

	if res.ExitCode != 0 || res.Error != "" {
		t.Fatalf("ExitCode=%d Error=%q, want success (res=%+v)", res.ExitCode, res.Error, res)
	}
	if res.Message != "OK" {
		t.Fatalf("Message=%q, want %q (res=%+v)", res.Message, "OK", res)
	}

	if gotName != "opencode" {
		t.Fatalf("command name=%q, want %q", gotName, "opencode")
	}
	if len(gotArgs) == 0 {
		t.Fatalf("missing args")
	}
	// If stdin mode is used, targetArg becomes "-" and OpencodeBackend would omit the prompt entirely.
	// We must still pass the prompt as positional args.
	if gotArgs[len(gotArgs)-1] != task {
		t.Fatalf("last arg=%q, want task=%q (args=%v)", gotArgs[len(gotArgs)-1], task, gotArgs)
	}
	if fake.StdinContents() != "" {
		t.Fatalf("unexpected stdin write: %q", fake.StdinContents())
	}
}

