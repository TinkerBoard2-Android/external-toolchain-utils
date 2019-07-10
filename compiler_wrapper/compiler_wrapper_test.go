package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestAddCommonFlags(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		ctx.cfg.commonFlags = []string{"-someflag"}
		cmd := ctx.must(callCompiler(ctx, ctx.cfg,
			ctx.newCommand(gccX86_64, mainCc)))
		if err := verifyArgOrder(cmd, "-someflag", mainCc); err != nil {
			t.Error(err)
		}
	})
}

func TestAddGccConfigFlags(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		ctx.cfg.gccFlags = []string{"-someflag"}
		cmd := ctx.must(callCompiler(ctx, ctx.cfg,
			ctx.newCommand(gccX86_64, mainCc)))
		if err := verifyArgOrder(cmd, "-someflag", mainCc); err != nil {
			t.Error(err)
		}
	})
}

func TestAddClangConfigFlags(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		ctx.cfg.clangFlags = []string{"-someflag"}
		cmd := ctx.must(callCompiler(ctx, ctx.cfg,
			ctx.newCommand(clangX86_64, mainCc)))
		if err := verifyArgOrder(cmd, "-someflag", mainCc); err != nil {
			t.Error(err)
		}
	})
}

func TestLogGeneralExecError(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		testOldWrapperPaths := []string{
			"",
			filepath.Join(ctx.tempDir, "fakewrapper"),
		}
		for _, testOldWrapperPath := range testOldWrapperPaths {
			ctx.cfg.oldWrapperPath = testOldWrapperPath
			// Note: No need to write the old wrapper as we don't execute
			// it due to the general error from the new error.
			ctx.cmdMock = func(cmd *command, stdout io.Writer, stderr io.Writer) error {
				return errors.New("someerror")
			}
			stderr := ctx.mustFail(callCompiler(ctx, ctx.cfg, ctx.newCommand(gccX86_64, mainCc)))
			if err := verifyInternalError(stderr); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(stderr, gccX86_64) {
				t.Errorf("could not find compiler path on stderr. Got: %s", stderr)
			}
			if !strings.Contains(stderr, "someerror") {
				t.Errorf("could not find original error on stderr. Got: %s", stderr)
			}
		}
	})
}

func TestLogMissingCCacheExecError(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		ctx.cfg.useCCache = true

		testOldWrapperPaths := []string{
			"",
			filepath.Join(ctx.tempDir, "fakewrapper"),
		}
		for _, testOldWrapperPath := range testOldWrapperPaths {
			ctx.cfg.oldWrapperPath = testOldWrapperPath
			// Note: No need to write the old wrapper as we don't execute
			// it due to the general error from the new error.
			ctx.cmdMock = func(cmd *command, stdout io.Writer, stderr io.Writer) error {
				return syscall.ENOENT
			}
			ctx.stderrBuffer.Reset()
			stderr := ctx.mustFail(callCompiler(ctx, ctx.cfg, ctx.newCommand(gccX86_64, mainCc)))
			if err := verifyNonInternalError(stderr, "ccache not found under .*. Please install it"); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestLogExitCodeErrorWhenComparingToOldWrapper(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		ctx.cfg.mockOldWrapperCmds = false
		ctx.cfg.oldWrapperPath = filepath.Join(ctx.tempDir, "fakewrapper")

		ctx.cmdMock = func(cmd *command, stdout io.Writer, stderr io.Writer) error {
			writeMockWrapper(ctx, &mockWrapperConfig{
				Cmds: []*mockWrapperCmd{
					{
						Path:     cmd.path,
						Args:     cmd.args,
						ExitCode: 2,
					},
				},
			})
			fmt.Fprint(stderr, "someerror")
			return newExitCodeError(2)
		}

		exitCode := callCompiler(ctx, ctx.cfg, ctx.newCommand(gccX86_64, mainCc))
		if exitCode != 2 {
			t.Fatalf("Expected exit code 2. Got: %d", exitCode)
		}
		if err := verifyNonInternalError(ctx.stderrString(), "someerror"); err != nil {
			t.Fatal(err)
		}
	})
}

func TestErrorOnLogRusageAndForceDisableWError(t *testing.T) {
	withTestContext(t, func(ctx *testContext) {
		ctx.env = []string{
			"FORCE_DISABLE_WERROR=1",
			"GETRUSAGE=" + filepath.Join(ctx.tempDir, "rusage.log"),
		}
		stderr := ctx.mustFail(callCompiler(ctx, ctx.cfg, ctx.newCommand(gccX86_64, mainCc)))
		if err := verifyNonInternalError(stderr, "GETRUSAGE is meaningless with FORCE_DISABLE_WERROR"); err != nil {
			t.Error(err)
		}
	})
}

func TestPrintUserCompilerError(t *testing.T) {
	buffer := bytes.Buffer{}
	printCompilerError(&buffer, newUserErrorf("abcd"))
	if buffer.String() != "abcd\n" {
		t.Errorf("Unexpected string. Got: %s", buffer.String())
	}
}

func TestPrintOtherCompilerError(t *testing.T) {
	buffer := bytes.Buffer{}
	printCompilerError(&buffer, errors.New("abcd"))
	if buffer.String() != "Internal error. Please report to chromeos-toolchain@google.com.\nabcd\n" {
		t.Errorf("Unexpected string. Got: %s", buffer.String())
	}
}
