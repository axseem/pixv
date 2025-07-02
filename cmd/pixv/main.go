package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/axseem/pixv"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "pixv",
		Version: "v0.3.0",
		Usage:   "A CLI tool to vectorize pixel-art images",
		Authors: []any{"axseem: max@axseem.me"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "method",
				Aliases: []string{"m"},
				Value:   string(pixv.StrategyPath),
				Usage:   "Vectorization method: path, rectangle, or square",
			},
			&cli.IntFlag{
				Name:    "scale",
				Aliases: []string{"s"},
				Value:   1,
				Usage:   "Scale multiplier for the output SVG",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path. Defaults to input with .svg extension",
			},
		},
		Action: runVectorize,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runVectorize(ctx context.Context, cmd *cli.Command) error {
	if !cmd.Args().Present() {
		_ = cli.ShowAppHelp(cmd)
		return errors.New("image path argument is required")
	}
	inputPath := cmd.Args().First()

	scale := cmd.Int("scale")
	if scale < 1 {
		return errors.New("scale must be a positive integer")
	}

	method := pixv.Strategy(cmd.String("method"))
	if !method.IsValid() {
		return fmt.Errorf("invalid method %q, must be one of: path, rectangle, square", method)
	}

	img, err := openImage(inputPath)
	if err != nil {
		return fmt.Errorf("could not open image %q: %w", inputPath, err)
	}

	opts := pixv.Options{
		Method: method,
		Scale:  scale,
	}
	svgContent, err := pixv.Vectorize(ctx, img, opts)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("operation cancelled")
		}
		return fmt.Errorf("failed to vectorize image: %w", err)
	}

	outputPath := cmd.String("output")
	if outputPath == "" {
		outputPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".svg"
	}

	if err := os.WriteFile(outputPath, []byte(svgContent), 0644); err != nil {
		return fmt.Errorf("failed to write svg to %q: %w", outputPath, err)
	}

	fmt.Printf("Successfully vectorized %q to %q\n", inputPath, outputPath)
	return nil
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
