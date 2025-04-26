package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Dionid/sqli"
	"github.com/spf13/cobra"
	"github.com/xo/xo/templates"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
)

// NewTemplateSet creates a new templates set.
func NewTemplateSet(ctx context.Context) (*templates.Set, error) {
	var err error

	// build template ts
	ts := templates.NewDefaultTemplateSet(ctx)

	// load specified template
	target := "templates"

	sub, err := fs.Sub(sqli.XoTemplates, target)
	if err != nil {
		return nil, err
	}

	// add template
	if target, err = ts.Add(ctx, target, sub, true); err != nil {
		return nil, err
	}

	// use
	ts.Use(target)

	return ts, nil
}

type GenerateCmdConfig struct {
	Out    string
	Schema string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "sqli",
		Short: "SQLi CLI tool",
	}

	generateCmdConfig := &GenerateCmdConfig{}

	generateCmd := &cobra.Command{
		Use:   "generate <database url>",
		Short: "Generate SQLi Query Builder based on the database schema",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Please provide a database URL")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, cmdArgs []string) {
			ctx := context.Background()

			println("Generating files...")

			err := sqli.Generate(
				ctx,
				sqli.GenerateCmdOpts{
					DbSchema: generateCmdConfig.Schema,
					Out:      generateCmdConfig.Out,

					DbUrl: cmdArgs[0],
				},
			)
			if err != nil {
				fmt.Printf("Error generating: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Files generated successfully into " + generateCmdConfig.Out)
		},
	}

	generateCmd.Flags().StringVarP(
		&generateCmdConfig.Out,
		"out",
		"o",
		"db",
		"out path (default 'models')",
	)

	generateCmd.Flags().StringVarP(
		&generateCmdConfig.Schema,
		"schema",
		"c",
		"public",
		"database schema name",
	)

	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
