package main

import (
	"context"
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

func main() {
	rootCmd := &cobra.Command{
		Use:   "sqli",
		Short: "SQLification CLI tool",
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate files using xo",
		Run: func(cmd *cobra.Command, args []string) {
			defaultSrc := "../../xo"
			defaultOut := "./db"
			defaultSchema := "public"

			if len(args) < 1 {
				fmt.Println("Please provide a database URL")
				return
			}

			err := xo.Run(
				context.Background(),
				"xo",
				"0.0.0-dev",
				"--src",
				defaultSrc,
				"--out",
				defaultOut,
				"--schema",
				defaultSchema,
				"schema",
				args[0],
			)
			if err != nil {
				fmt.Printf("Error running xo: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Files generated successfully into " + defaultOut)

			// xoCmd := exec.Command("xo --src ../../xo --out ./db --schema public schema ${DB_PG}")
			// xoCmd.Stdout = os.Stdout
			// xoCmd.Stderr = os.Stderr

			// if err := xoCmd.Run(); err != nil {
			// 	fmt.Printf("Error executing xo command: %v\n", err)
			// 	os.Exit(1)
			// }
		},
	}

	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
