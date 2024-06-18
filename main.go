package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/axseem/pixv/convert"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pixv",
		Usage: "Enter image path to vectorize it",
		Authors: []*cli.Author{
			{
				Name:  "@axseem",
				Email: "https://github.com/axseem",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "method",
				Aliases:     []string{"m"},
				Usage:       "Choose vectorization method",
				Value:       "path",
				DefaultText: "path",
			},
			&cli.StringFlag{
				Name:        "scale",
				Aliases:     []string{"s"},
				Usage:       "Change the scale of the pixels",
				Value:       "1",
				DefaultText: "1",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				return cli.ShowAppHelp(cCtx)
			}

			filePath := cCtx.Args().First()
			file := path.Base(filePath)
			fileName := strings.TrimSuffix(file, filepath.Ext(file))

			img, err := convert.OpenImage(filePath)
			if err != nil {
				return err
			}

			scale, err := strconv.Atoi(cCtx.String("scale"))
			if err != nil {
				return fmt.Errorf("invalid scale value")
			}
			if scale < 1 {
				return fmt.Errorf("only positive integer can be used as scale value")
			}

			var svgString string
			if cCtx.String("method") == "path" {
				svgString, err = convert.Path(img, uint(scale))
				if err != nil {
					return err
				}
			} else if cCtx.String("method") == "rectangle" {
				svgString, err = convert.Rectangle(img, uint(scale))
				if err != nil {
					return err
				}
			} else if cCtx.String("method") == "square" {
				svgString, err = convert.Square(img, uint(scale))
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("invalid method")
			}

			f, err := os.Create(fileName + ".svg")
			if err != nil {
				return err
			}

			f.WriteString(svgString)
			return f.Close()
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
