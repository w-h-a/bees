package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var knownKeys = map[string]bool{
	"issue-prefix": true,
	"author":       true,
	"json":         true,
}

type config struct {
	data map[string]any
}

func (c *config) Get(key string) string {
	v, ok := c.data[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func (c *config) Set(key, value string) error {
	if key == "json" {
		switch strings.ToLower(value) {
		case "true", "1":
			c.data[key] = true
		case "false", "0":
			c.data[key] = false
		default:
			return fmt.Errorf("json must be true/false or 1/0, got %q", value)
		}
		return nil
	}
	c.data[key] = value
	return nil
}

func (c *config) IssuePrefix() string {
	if v, ok := c.data["issue-prefix"].(string); ok {
		return v
	}
	return ""
}

func (c *config) Keys() []string {
	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func newConfig() *config {
	return &config{data: map[string]any{}}
}

func newConfigCmd(beesDir *string, cfg **config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	cmd.AddCommand(newConfigSetCmd(beesDir, cfg))
	cmd.AddCommand(newConfigGetCmd(cfg))
	cmd.AddCommand(newConfigListCmd(cfg))

	return cmd
}

func newConfigSetCmd(beesDir *string, cfg **config) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]

			if !knownKeys[key] {
				slog.Warn("unknown config key", "key", key)
			}

			if err := (*cfg).Set(key, value); err != nil {
				return err
			}

			if err := saveConfig(*beesDir, *cfg); err != nil {
				return err
			}

			if !jsonOutput {
				fmt.Printf("%s = %s\n", key, value)
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", " ")

			return enc.Encode(map[string]string{
				"key":   key,
				"value": value,
			})
		},
	}
}

func newConfigGetCmd(cfg **config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value, source := resolveConfig(*cfg, key)

			if !jsonOutput {
				fmt.Println(value)
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")

			return enc.Encode(map[string]string{
				"key":    key,
				"value":  value,
				"source": source,
			})
		},
	}
}

func newConfigListCmd(cfg **config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			type entry struct {
				Key    string `json:"key"`
				Value  string `json:"value"`
				Source string `json:"source"`
			}

			var entries []entry
			seen := make(map[string]bool)

			for _, key := range sortedKnownKeys() {
				value, source := resolveConfig(*cfg, key)
				if value != "" {
					entries = append(entries, entry{key, value, source})
				}
				seen[key] = true
			}

			for _, key := range (*cfg).Keys() {
				if seen[key] {
					continue
				}
				value, source := resolveConfig(*cfg, key)
				if value != "" {
					entries = append(entries, entry{key, value, source})
				}
			}

			if !jsonOutput {
				for _, e := range entries {
					fmt.Printf("%s = %s (%s)\n", e.Key, e.Value, e.Source)
				}
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")

			return enc.Encode(entries)
		},
	}
}

func saveConfig(beesDir string, cfg *config) error {
	data, err := yaml.Marshal(cfg.data)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(filepath.Join(beesDir, "config.yaml"), data, 0o644)
}

func loadConfig(beesDir string) (*config, error) {
	data, err := os.ReadFile(filepath.Join(beesDir, "config.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return newConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	if m == nil {
		m = make(map[string]any)
	}

	return &config{data: m}, nil
}

func resolveConfig(cfg *config, key string) (value, source string) {
	envKey := "BEES_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	if v := os.Getenv(envKey); v != "" {
		return v, "env"
	}
	v := cfg.Get(key)
	if v != "" {
		return v, "config"
	}
	return "", ""
}

func sortedKnownKeys() []string {
	keys := make([]string, 0, len(knownKeys))
	for k := range knownKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
