/*exportstruct really rocks
 */
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"of"},
			Value:   "types.go",
			Usage:   "Generated code exports to `FILE`",
		},
		&cli.StringFlag{
			Name:     "user",
			Aliases:  []string{"u"},
			Value:    "postgres",
			Usage:    "Database `USER`",
			EnvVars:  []string{"ES_USER"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Value:    "password",
			Usage:    "Database `PASSWORD`",
			EnvVars:  []string{"ES_PASS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"a"},
			Value:   "::1",
			Usage:   "Address of db `HOST`",
			EnvVars: []string{"ES_HOST"},
		},
		&cli.StringFlag{
			Name:    "port",
			Value:   "5432",
			Usage:   "`PORT` of db",
			EnvVars: []string{"ES_PORT"},
		},
		&cli.StringFlag{
			Name:    "db",
			Value:   "postgres",
			Usage:   "`NAME` of db",
			EnvVars: []string{"ES_DB"},
		},
		&cli.BoolFlag{
			Name:    "ssl-mode",
			Aliases: []string{"sm"},
			Value:   false,
			Usage:   "Set ssl-mode to verify-full or disable (default: `disable`)",
		},
		&cli.BoolFlag{
			Name:  "json",
			Value: false,
			Usage: "Add json tags, camelcase",
		},
		&cli.BoolFlag{
			Name:  "sql",
			Value: false,
			Usage: "Add sql tags, snake_case",
		},
	}

	app.Action = export

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func export(c *cli.Context) error {
	file, user, pass, host, port, db := c.String("file"), c.String("user"), c.String("password"),
		c.String("host"), c.String("port"), c.String("db")
	println(file)
	ssl := c.Bool("ssl-mode")
	verify := "disable"
	if ssl {
		verify = "verify-full"
	}

	psql := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=%s`, host, port, user, pass, db, verify)

	cmd := exec.Command("psql", psql, "-f", "query.sql")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err = cmd.Start(); err != nil {
		panic(err)
	}

	out, errout := ioutil.ReadAll(stdout)
	if err = cmd.Wait(); err != nil {
		fmt.Println(errout)
		panic(err)
	}

	buf := bytes.NewBufferString("package main\n\nimport \"database/sql\"\n\n")
	_, err = buf.Write(out)
	if err != nil {
		panic(err)
	}

	structs := buf.Bytes()
	structs = bytes.ReplaceAll(structs, []byte(`\n`), []byte("\n"))
	structs = bytes.ReplaceAll(structs, []byte(`\t`), []byte("\t"))

	err = ioutil.WriteFile(file, structs, 0600)
	if err != nil {
		panic(err)
	}
	cmd = exec.Command("gofmt", "-s", "-w", file)
	println(cmd.String())
	if err = cmd.Run(); err != nil {
		panic(err)
	}
	cmd = exec.Command("goimports", "-w", file)
	println(cmd.String())
	if err = cmd.Run(); err != nil {
		panic(err)
	}

	j, s := c.Bool("json"), c.Bool("sql")
	if j {
		cmd = exec.Command("gomodifytags", "-file", file, "-w", "-all", "-add-tags", "json", "-transform", "camelcase")
		if err = cmd.Run(); err != nil {
			panic(err)
		}
	}
	if s {
		cmd = exec.Command("gomodifytags", "-file", file, "-w", "-all", "-add-tags", "sql")
		if err = cmd.Run(); err != nil {
			panic(err)
		}
	}

	return nil
}
