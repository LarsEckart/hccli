package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func CreateQueryAnnotationCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-query-annotation",
		Category: "Query Annotations",
		Usage:    "Create a query annotation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "query-id",
				Usage:    "Query ID to annotate",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Annotation name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Annotation description",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			annotation := &api.QueryAnnotation{
				Name:    cmd.String("name"),
				QueryID: cmd.String("query-id"),
			}
			if v := cmd.String("description"); v != "" {
				annotation.Description = v
			}

			created, err := client.CreateQueryAnnotation(ctx, cmd.String("dataset"), annotation)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func ListQueryAnnotationsCmd() *cli.Command {
	return &cli.Command{
		Name:     "query-annotations",
		Category: "Query Annotations",
		Usage:    "List all query annotations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			annotations, err := client.ListQueryAnnotations(ctx, cmd.String("dataset"))
			if err != nil {
				return err
			}

			return printJSON(annotations)
		},
	}
}

func GetQueryAnnotationCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-query-annotation",
		Category: "Query Annotations",
		Usage:    "Get a query annotation by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Query annotation ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			annotation, err := client.GetQueryAnnotation(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}

			return printJSON(annotation)
		},
	}
}

func UpdateQueryAnnotationCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-query-annotation",
		Category: "Query Annotations",
		Usage:    "Update a query annotation by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Query annotation ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Annotation name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "query-id",
				Usage:    "Query ID to annotate",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Annotation description",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			annotation := &api.QueryAnnotation{
				Name:    cmd.String("name"),
				QueryID: cmd.String("query-id"),
			}
			if v := cmd.String("description"); v != "" {
				annotation.Description = v
			}

			updated, err := client.UpdateQueryAnnotation(ctx, cmd.String("dataset"), cmd.String("id"), annotation)
			if err != nil {
				return err
			}

			return printJSON(updated)
		},
	}
}

func DeleteQueryAnnotationCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-query-annotation",
		Category: "Query Annotations",
		Usage:    "Delete a query annotation by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Query annotation ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteQueryAnnotation(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}
