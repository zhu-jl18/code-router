package main

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestRunCodexTaskWithContext_Opencode_ErrorEventIsSurfaced(t *testing.T) {
	defer resetTestHooks()

	var waitErr error
	if runtime.GOOS == "windows" {
		waitErr = exec.Command("cmd", "/c", "exit", "7").Run()
	} else {
		waitErr = exec.Command("sh", "-c", "exit 7").Run()
	}
	exitErr, ok := waitErr.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 7 {
		t.Fatalf("failed to build ExitError with code 7, got %T: %v", waitErr, waitErr)
	}

	opencodeErr := `{"type":"error","timestamp":1,"sessionID":"ses_err","error":{"name":"APIError","data":{"message":"Unauthorized: {\"type\":\"error\",\"error\":{\"type\":\"CreditsError\",\"message\":\"No payment method\"}}","statusCode":401,"responseBody":"{\"type\":\"error\",\"error\":{\"type\":\"CreditsError\",\"message\":\"No payment method\"}}","metadata":{"url":"https://opencode.ai/zen/v1/responses"}}}}` + "\n"

	fake := newFakeCmd(fakeCmdConfig{
		StdoutPlan: []fakeStdoutEvent{
			{Data: opencodeErr},
		},
		WaitErr: exitErr,
	})

	newCommandRunner = func(ctx context.Context, name string, args ...string) commandRunner {
		return fake
	}

	res := runCodexTaskWithContext(
		context.Background(),
		TaskSpec{ID: "opencode-err", Task: "hello", WorkDir: ".", Mode: "new"},
		OpencodeBackend{},
		nil,
		false,
		true,
		5,
	)

	if res.ExitCode != 7 {
		t.Fatalf("ExitCode=%d, want 7 (res=%+v)", res.ExitCode, res)
	}
	if res.SessionID != "ses_err" {
		t.Fatalf("SessionID=%q, want %q (res=%+v)", res.SessionID, "ses_err", res)
	}
	if !strings.Contains(res.Error, "CreditsError: No payment method") {
		t.Fatalf("Error=%q, want to contain %q (res=%+v)", res.Error, "CreditsError: No payment method", res)
	}
	if res.Message != "CreditsError: No payment method" {
		t.Fatalf("Message=%q, want %q (res=%+v)", res.Message, "CreditsError: No payment method", res)
	}
}

