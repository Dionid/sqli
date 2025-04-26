package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	xo "github.com/xo/xo/cmd"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
)

type GenerateCmdConfig struct {
	Src    string
	Out    string
	Schema string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "sqli",
		Short: "SQLification CLI tool",
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
		Run: func(cmd *cobra.Command, args []string) {
			xoArgs := make([]string, 0)

			xoArgs = append(xoArgs, "--src", generateCmdConfig.Src)
			xoArgs = append(xoArgs, "--schema", generateCmdConfig.Schema)
			xoArgs = append(xoArgs, "--out", generateCmdConfig.Out)

			xoArgs = append(xoArgs, "schema", args[0])

			println("Generating files...")

			err := xo.Run(
				context.Background(),
				"xo",
				"0.0.0-dev",
				xoArgs...,
			)
			if err != nil {
				fmt.Printf("Error generating: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Files generated successfully into " + generateCmdConfig.Out)
		},
	}

	generateCmd.Flags().StringVar(
		&generateCmdConfig.Src,
		"src",
		"../../xo",
		"directory with templates",
	)

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
