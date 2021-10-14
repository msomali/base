/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package cli

import (
	"fmt"
	"github.com/techcraftt/base/api/http"
	clix "github.com/urfave/cli/v2"
	"os"
	"strconv"
)

type (
	App struct {
		http *http.Client
		app  *clix.App
	}

	ApiClient struct {
		*http.Client
		//logger *log.Logger
	}
)

func (apiClient *ApiClient) divideCommand() *clix.Command {
	return &clix.Command{
		Name:        "divide",
		Aliases:     []string{"div"},
		Usage:       "divide two numbers",
		UsageText:   "div [number] [number]",
		Description: "perform division",
		ArgsUsage:   "args usage div command",
		Action: func(c *clix.Context) error {
			if c.NArg() > 1 {
				args := c.Args()
				aStr := args.Get(0)
				bStr := args.Get(1)
				a, err := strconv.Atoi(aStr)
				if err != nil {
					return err
				}
				b, err := strconv.Atoi(bStr)
				if err != nil {
					return err
				}
				res, err := apiClient.Divide(int64(a), int64(b))
				if err != nil {
					return err
				}
				fmt.Printf("answer: %v\n", res)
				return nil
			}

			return fmt.Errorf("not enough arguments")

		},
	}
}

func commands(comm ...*clix.Command) []*clix.Command {
	var commands []*clix.Command
	for _, command := range comm {
		commands = append(commands, command)
	}
	return commands
}

func flags(fs ...clix.Flag) []clix.Flag {
	var flgs []clix.Flag
	for _, flg := range fs {
		flgs = append(flgs, flg)
	}
	return flgs
}

func authors(auth ...*clix.Author) []*clix.Author {
	var authors []*clix.Author
	for _, author := range auth {
		authors = append(authors, author)
	}
	return authors
}

func New(base string, port uint64, debug bool) *App {
	client := http.NewClient(base, port, debug)
	author1 := &clix.Author{
		Name:  "Pius Alfred",
		Email: "me.pius1102@gmail.com",
	}
	apiClient := ApiClient{client}

	divideCommand := apiClient.divideCommand()

	app := &clix.App{
		Name:                 "calc",
		Usage:                "perform simple calculations",
		UsageText:            "calc div <first-integer> <second-integer>",
		Version:              "1.0.0",
		Description:          "perform simple addition and division of base64 integers",
		Commands:             commands(divideCommand),
		Flags:                flags(),
		EnableBashCompletion: true,
		Authors:              authors(author1),
		Copyright:            "MIT Licence, Creative Commons",
		ErrWriter:            os.Stderr,
	}

	return &App{
		http: client,
		app:  app,
	}
}

func (app *App) Run(args []string) error {
	return app.app.Run(args)
}
