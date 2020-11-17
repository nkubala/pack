package config_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/pack/internal/style"

	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	cmdConfig "github.com/buildpacks/pack/internal/commands/config"
	"github.com/buildpacks/pack/internal/config"
	ilogging "github.com/buildpacks/pack/internal/logging"
	"github.com/buildpacks/pack/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestTrustedBuilderCommand(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "TrustedBuilderCommands", testTrustedBuilderCommand, spec.Random(), spec.Report(report.Terminal{}))
}

func testTrustedBuilderCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		command      *cobra.Command
		logger       logging.Logger
		outBuf       bytes.Buffer
		tempPackHome string
		configPath   string
	)

	it.Before(func() {
		var err error

		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		tempPackHome, err = ioutil.TempDir("", "pack-home")
		h.AssertNil(t, err)
		h.AssertNil(t, os.Setenv("PACK_HOME", tempPackHome))
		configPath = filepath.Join(tempPackHome, "config.toml")

		command = cmdConfig.Config(logger, config.Config{}, configPath)
		command.SetOut(logging.GetWriterForLevel(logger, logging.InfoLevel))
		//command.SetErr(logging.GetWriterForLevel(logger, logging.ErrorLevel))
	})

	it.After(func() {
		h.AssertNil(t, os.Unsetenv("PACK_HOME"))
		h.AssertNil(t, os.RemoveAll(tempPackHome))
	})

	when("list", func() {
		var args = []string{"trusted-builder", "list"}

		it("shows suggested builders and locally trusted builder in alphabetical order", func() {
			builderName := "great-builder-" + h.RandString(8)

			command.SetArgs(args)
			h.AssertNil(t, command.Execute())
			h.AssertNotContains(t, outBuf.String(), builderName)
			h.AssertContainsAllInOrder(t,
				outBuf,
				"gcr.io/buildpacks/builder:v1",
				"heroku/buildpacks:18",
				"paketobuildpacks/builder:base",
				"paketobuildpacks/builder:full",
				"paketobuildpacks/builder:tiny",
			)
			outBuf.Reset()

			configManager := newConfigManager(t, configPath)
			command = cmdConfig.Config(logger, configManager.configWithTrustedBuilders(builderName), configPath)
			command.SetArgs(args)
			h.AssertNil(t, command.Execute())

			h.AssertContainsAllInOrder(t,
				outBuf,
				"gcr.io/buildpacks/builder:v1",
				builderName,
				"heroku/buildpacks:18",
				"paketobuildpacks/builder:base",
				"paketobuildpacks/builder:full",
				"paketobuildpacks/builder:tiny",
			)
		})
	})

	when("trusted-builder add", func() {
		var args = []string{"trusted-builder", "add"}
		when("no builder is provided", func() {
			it("prints usage", func() {
				command.SetArgs(args)
				h.AssertError(t, command.Execute(), "accepts 1 arg(s)")
			})
		})

		when("can't write to config path", func() {
			it("fails", func() {
				tempPath := filepath.Join(tempPackHome, "non-existent-file.toml")
				h.AssertNil(t, ioutil.WriteFile(tempPath, []byte("something"), 0111))
				command = cmdConfig.Config(logger, config.Config{}, tempPath)
				command.SetOut(logging.GetWriterForLevel(logger, logging.InfoLevel))
				command.SetArgs(append(args, "some-builder"))
				h.AssertError(t, command.Execute(), "writing config")
			})
		})

		when("builder is provided", func() {
			when("builder is not already trusted", func() {
				it("updates the config", func() {
					command.SetArgs(append(args, "some-builder"))
					h.AssertNil(t, command.Execute())

					b, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)
					h.AssertContains(t, string(b), `[[trusted-builders]]
  name = "some-builder"`)
				})
			})

			when("builder is already trusted", func() {
				it("does nothing", func() {
					command.SetArgs(append(args, "some-already-trusted-builder"))
					h.AssertNil(t, command.Execute())
					oldContents, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)

					command.SetArgs(append(args, "some-already-trusted-builder"))
					h.AssertNil(t, command.Execute())

					newContents, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)
					h.AssertEq(t, newContents, oldContents)
				})
			})

			when("builder is a suggested builder", func() {
				it("does nothing", func() {
					h.AssertNil(t, ioutil.WriteFile(configPath, []byte(""), os.ModePerm))

					command.SetArgs(append(args, "paketobuildpacks/builder:base"))
					h.AssertNil(t, command.Execute())
					oldContents, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)
					h.AssertEq(t, string(oldContents), "")
				})
			})
		})
	})

	when("trusted-builder remove", func() {
		var (
			args          = []string{"trusted-builder", "remove"}
			configManager configManager
		)

		it.Before(func() {
			configManager = newConfigManager(t, configPath)
		})

		when("no builder is provided", func() {
			it("prints usage", func() {
				cfg := configManager.configWithTrustedBuilders()
				command := cmdConfig.Config(logger, cfg, configPath)
				command.SetArgs(args)
				command.SetOut(&outBuf)

				err := command.Execute()
				h.AssertError(t, err, "accepts 1 arg(s), received 0")
				h.AssertContains(t, outBuf.String(), "Usage:")
			})
		})

		when("builder is already trusted", func() {
			it("removes builder from the config", func() {
				builderName := "some-builder"

				cfg := configManager.configWithTrustedBuilders(builderName)
				command := cmdConfig.Config(logger, cfg, configPath)
				command.SetArgs(append(args, builderName))

				h.AssertNil(t, command.Execute())

				b, err := ioutil.ReadFile(configPath)
				h.AssertNil(t, err)
				h.AssertNotContains(t, string(b), builderName)

				h.AssertContains(t,
					outBuf.String(),
					fmt.Sprintf("Builder %s is no longer trusted", style.Symbol(builderName)),
				)
			})

			it("removes only the named builder when multiple builders are trusted", func() {
				untrustBuilder := "stop/trusting:me"
				stillTrustedBuilder := "very/safe/builder"

				cfg := configManager.configWithTrustedBuilders(untrustBuilder, stillTrustedBuilder)
				command := cmdConfig.Config(logger, cfg, configPath)
				command.SetArgs(append(args, untrustBuilder))

				h.AssertNil(t, command.Execute())

				b, err := ioutil.ReadFile(configPath)
				h.AssertNil(t, err)
				h.AssertContains(t, string(b), stillTrustedBuilder)
				h.AssertNotContains(t, string(b), untrustBuilder)
			})
		})

		when("builder wasn't already trusted", func() {
			it("does nothing and reports builder wasn't trusted", func() {
				neverTrustedBuilder := "never/trusted-builder"
				stillTrustedBuilder := "very/safe/builder"

				cfg := configManager.configWithTrustedBuilders(stillTrustedBuilder)
				command := cmdConfig.Config(logger, cfg, configPath)
				command.SetArgs(append(args, neverTrustedBuilder))

				h.AssertNil(t, command.Execute())

				b, err := ioutil.ReadFile(configPath)
				h.AssertNil(t, err)
				h.AssertContains(t, string(b), stillTrustedBuilder)
				h.AssertNotContains(t, string(b), neverTrustedBuilder)

				h.AssertContains(t,
					outBuf.String(),
					fmt.Sprintf("Builder %s wasn't trusted", style.Symbol(neverTrustedBuilder)),
				)
			})
		})

		when("builder is a suggested builder", func() {
			it("does nothing and reports that ", func() {
				builder := "paketobuildpacks/builder:base"
				command := cmdConfig.Config(logger, config.Config{}, configPath)
				command.SetArgs(append(args, builder))

				err := command.Execute()
				h.AssertError(t, err, fmt.Sprintf("Builder %s is a suggested builder, and is trusted by default", style.Symbol(builder)))
			})
		})
	})
}
