package sqli

import (
	"context"
	"io/fs"

	xo "github.com/xo/xo/cmd"
	"github.com/xo/xo/templates"
)

func NewTemplateSet(ctx context.Context) (*templates.Set, error) {
	var err error

	// # Build template ts
	ts := templates.NewDefaultTemplateSet(ctx)

	// # Load specified template
	target := "templates"

	sub, err := fs.Sub(XoTemplates, target)
	if err != nil {
		return nil, err
	}

	// # Add template
	if target, err = ts.Add(ctx, target, sub, true); err != nil {
		return nil, err
	}

	// # Use
	ts.Use(target)

	return ts, nil
}

type GenerateCmdOpts struct {
	Out      string
	DbSchema string

	DbUrl string
}

func Generate(
	ctx context.Context,
	opts GenerateCmdOpts,
) error {
	xoCmdArgs := make([]string, 0)

	// # Flags
	xoCmdArgs = append(xoCmdArgs, "--schema", opts.DbSchema)
	xoCmdArgs = append(xoCmdArgs, "--out", opts.Out)

	// # Args
	xoCmdArgs = append(xoCmdArgs, "schema", opts.DbUrl)

	// # Create template set
	ts, err := NewTemplateSet(ctx)
	if err != nil {
		return err
	}

	// # Create args
	tmpArgs := xo.NewArgs(ts.Target(), ts.Targets()...)

	// # Create root xo command
	xoCmd, err := xo.RootCommand(ctx, "xo", "0.0.0-dev", ts, tmpArgs, xoCmdArgs...)
	if err != nil {
		return err
	}

	// # Execute
	err = xoCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
