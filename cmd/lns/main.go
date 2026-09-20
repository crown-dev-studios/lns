package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/crown-dev-studios/lns/internal/caddy"
	"github.com/crown-dev-studios/lns/internal/config"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:               "lns",
	Short:             "Run local services at stable local names",
	Args:              cobra.NoArgs,
	SilenceUsage:      true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Long: `lns discovers the default local app graph without repository setup,
leases dynamic ports, starts safe Compose dependencies, and routes stable
http://*.localhost names through Caddy.`,
	Example: `  lns
  lns plan
  lns run web`,
}

func init() {
	rootCmd.AddCommand(versionCmd, startCmd, stopCmd, configCmd, doctorCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), versionString())
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the shared local proxy",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		installation, err := requireCaddy(cmd.Context())
		if err != nil {
			return err
		}
		if err := config.EnsureConfigDirs(); err != nil {
			return fmt.Errorf("create LNS state directories: %w", err)
		}
		if _, err := caddy.RegenerateAllCaddyfiles(); err != nil {
			return fmt.Errorf("generate proxy configuration: %w", err)
		}
		if isTCPListening(config.CaddyAdminAddr) {
			printSuccess("Proxy is already running")
			return nil
		}
		if err := startCaddy(cmd.Context(), installation.Path, config.GetGlobalCaddyfilePath(), config.DefaultHTTPPort); err != nil {
			return err
		}
		printSuccess("Proxy started at http://*.localhost")
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the shared local proxy",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isTCPListening(config.CaddyAdminAddr) {
			printSuccess("Proxy is already stopped")
			return nil
		}
		installation, err := requireCaddy(cmd.Context())
		if err != nil {
			return fmt.Errorf("%w; stop the process listening at %s manually", err, config.CaddyAdminAddr)
		}
		command := exec.Command(installation.Path, "stop", "--address", config.CaddyAdminAddr)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return fmt.Errorf("stop Caddy: %w", err)
		}
		printSuccess("Proxy stopped")
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show machine-local LNS paths",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("State directory: %s\n", config.GetConfigDir())
		fmt.Printf("Proxy config:    %s\n", config.GetGlobalCaddyfilePath())
		fmt.Printf("Runtime routes:  %s\n", config.GetRuntimePath())
		fmt.Printf("Dependencies:    %s\n", filepath.Join(config.GetConfigDir(), "dependencies"))
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the local runtime requirements",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var issues []string
		if installation, err := requireCaddy(cmd.Context()); err != nil {
			issues = append(issues, err.Error())
		} else {
			printSuccess("Caddy %s: %s", installation.Version, installation.Path)
		}
		if path, err := exec.LookPath("docker"); err != nil {
			printWarning("Docker is missing; projects without Compose dependencies can still run")
		} else if output, err := exec.Command(path, "compose", "version", "--short").CombinedOutput(); err != nil {
			issues = append(issues, "Docker Compose is unavailable: "+strings.TrimSpace(string(output)))
		} else {
			printSuccess("Docker Compose: %s", strings.TrimSpace(string(output)))
		}
		if isTCPListening(config.CaddyAdminAddr) {
			printSuccess("Proxy is running")
		} else {
			printWarning("Proxy will start on the first `lns` run")
		}
		if len(issues) > 0 {
			return fmt.Errorf("doctor found %d issue(s): %s", len(issues), strings.Join(issues, "; "))
		}
		return nil
	},
}

func printSuccess(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", color.GreenString("✓"), fmt.Sprintf(format, args...))
}

func printWarning(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", color.YellowString("!"), fmt.Sprintf(format, args...))
}
