// Package main provides a Docker CLI plugin for managing project dependencies.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/cli/cli-plugins/metadata"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

const (
	manifestFile = "docker-deps.json"
)

// Manifest represents the structure of the dependency file.
type Manifest struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

// SetProjectInfo updates the basic project identification metadata.
func (m *Manifest) SetProjectInfo(name, version string) {
	m.Name = name
	m.Version = version
}

// AddDependency adds or updates a single image dependency.
func (m *Manifest) AddDependency(name, version string) {
	if m.Dependencies == nil {
		m.Dependencies = make(map[string]string)
	}
	m.Dependencies[name] = version
}

// RemoveDependency removes an image from the project manifest.
func (m *Manifest) RemoveDependency(name string) error {
	if _, ok := m.Dependencies[name]; !ok {
		return fmt.Errorf("dependency %s not found", name)
	}
	delete(m.Dependencies, name)
	return nil
}

// UpdateDependency changes the version of an existing dependency.
func (m *Manifest) UpdateDependency(name, version string) error {
	if _, ok := m.Dependencies[name]; !ok {
		return fmt.Errorf("dependency %s not found", name)
	}
	m.Dependencies[name] = version
	return nil
}

// GetDependency retrieves the version of a given image dependency.
func (m *Manifest) GetDependency(name string) (string, bool) {
	v, ok := m.Dependencies[name]
	return v, ok
}

// loadManifest reads the manifest from a file and parses its content.
func loadManifest(file string) (*Manifest, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{Dependencies: make(map[string]string)}, nil
		}
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if manifest.Dependencies == nil {
		manifest.Dependencies = make(map[string]string)
	}
	return &manifest, nil
}

// saveManifest writes the current manifest state to a JSON file.
func saveManifest(m *Manifest, file string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func main() {
	meta := metadata.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Philippe Court <philippe.court@gmail.com>",
		Version:       "1.0.0",
	}

	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "deps",
			Short: "Manage project dependencies with a manifest file",
		}

		cmd.AddCommand(newAddCommand(dockerCli))
		cmd.AddCommand(newRemoveCommand(dockerCli))
		cmd.AddCommand(newUpdateCommand(dockerCli))
		cmd.AddCommand(newListCommand(dockerCli))
		cmd.AddCommand(newInfoCommand(dockerCli))
		cmd.AddCommand(newWhyCommand(dockerCli))
		cmd.AddCommand(newVersionCommand(dockerCli, meta))
		cmd.AddCommand(newInitCommand(dockerCli))
		cmd.AddCommand(newInstallCommand(dockerCli))

		return cmd
	}, meta)
}

func newInstallCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Pull all images listed as dependencies in the manifest",
		RunE: func(_ *cobra.Command, _ []string) error {
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}

			if len(m.Dependencies) == 0 {
				_, _ = fmt.Fprintln(dockerCli.Out(), "No dependencies to install.")
				return nil
			}

			for dep, ver := range m.Dependencies {
				imgName := fmt.Sprintf("%s:%s", dep, ver)
				_, _ = fmt.Fprintf(dockerCli.Out(), "Pulling %s...\n", imgName)

				err := runDockerPull(dockerCli, imgName)
				if err != nil {
					_, _ = fmt.Fprintf(dockerCli.Err(), "Error pulling %s: %v\n", imgName, err)
					continue
				}
			}
			return nil
		},
	}
}

func runDockerPull(dockerCli command.Cli, imgName string) error {
	apiClient := dockerCli.Client()

	out, err := apiClient.ImagePull(context.Background(), imgName, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	dec := json.NewDecoder(out)
	for {
		var m map[string]interface{}
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		status, _ := m["status"].(string)
		id, _ := m["id"].(string)
		progress, _ := m["progress"].(string)

		if status != "" {
			if id != "" {
				if progress != "" {
					_, _ = fmt.Fprintf(dockerCli.Out(), "\r[%s] %s: %s", id, status, progress)
				} else {
					_, _ = fmt.Fprintf(dockerCli.Out(), "\r[%s] %s", id, status)
				}
			} else {
				_, _ = fmt.Fprintf(dockerCli.Out(), "\n%s", status)
			}
		}
	}
	_, _ = fmt.Fprintln(dockerCli.Out())
	return nil
}

func newInitCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "init [name] [version]",
		Short: "Initialize or update the manifest with project name and version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			m.SetProjectInfo(args[0], args[1])
			if err := saveManifest(m, manifestFile); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(dockerCli.Out(), "Initialized manifest: %s version %s\n", args[0], args[1])
			return nil
		},
	}
}

func newAddCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "add [dependency] [version]",
		Short: "Add a dependency to the manifest",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			m.AddDependency(args[0], args[1])
			if err := saveManifest(m, manifestFile); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(dockerCli.Out(), "Added dependency %s version %s\n", args[0], args[1])
			return nil
		},
	}
}

func newRemoveCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [dependency]",
		Short: "Remove a dependency from the manifest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			if err := m.RemoveDependency(args[0]); err != nil {
				return err
			}
			if err := saveManifest(m, manifestFile); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(dockerCli.Out(), "Removed dependency %s\n", args[0])
			return nil
		},
	}
}

func newUpdateCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "update [dependency] [version]",
		Short: "Update a dependency version in the manifest",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			if err := m.UpdateDependency(args[0], args[1]); err != nil {
				return err
			}
			if err := saveManifest(m, manifestFile); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(dockerCli.Out(), "Updated dependency %s to version %s\n", args[0], args[1])
			return nil
		},
	}
}

func newListCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all dependencies in the manifest",
		RunE: func(_ *cobra.Command, _ []string) error {
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			if len(m.Dependencies) == 0 {
				_, _ = fmt.Fprintln(dockerCli.Out(), "No dependencies found.")
				return nil
			}
			_, _ = fmt.Fprintln(dockerCli.Out(), "Dependencies:")
			for dep, ver := range m.Dependencies {
				_, _ = fmt.Fprintf(dockerCli.Out(), "  - %s: %s\n", dep, ver)
			}
			return nil
		},
	}
}

func newInfoCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "info [dependency]",
		Short: "Get detailed information about a dependency",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			version, ok := m.GetDependency(args[0])
			if !ok {
				return fmt.Errorf("dependency %s not found", args[0])
			}
			_, _ = fmt.Fprintf(dockerCli.Out(), "Dependency: %s\nVersion: %s\n", args[0], version)
			return nil
		},
	}
}

func newWhyCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "why [dependency]",
		Short: "Explain why a dependency is in the manifest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			m, err := loadManifest(manifestFile)
			if err != nil {
				return err
			}
			if _, ok := m.GetDependency(args[0]); !ok {
				return fmt.Errorf("dependency %s not found in manifest", args[0])
			}
			_, _ = fmt.Fprintf(dockerCli.Out(), "The dependency \"%s\" is present because it was explicitly added as a project requirement.\n", args[0])
			return nil
		},
	}
}

func newVersionCommand(dockerCli command.Cli, meta metadata.Metadata) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the plugin version information",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd
			_, _ = fmt.Fprintf(dockerCli.Out(), "Docker Manifest Plugin Version: %s\nVendor: %s\n", meta.Version, meta.Vendor)
		},
	}
}
