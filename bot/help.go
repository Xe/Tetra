package tetra

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func helpHas(service string, command string) bool {
	command = strings.Join(strings.Split(command, " "), "_")

	if command == "help" || command == "_index" {
		service = "tetra"
	}

	_, err := os.Stat("help/" + service + "/" + command)

	return err == nil
}

func (c *Client) showHelp(target *Client, service, command string) {
	buffer := bytes.NewBufferString("")

	command = strings.TrimSpace(command)

	if command == "" {
		command = "_index"
	}

	if command == "help" || command == "_index" {
		service = "tetra"
	}

	text := template.New(service + "/" + command)
	helptext, _ := ioutil.ReadFile("help/" + service + "/" + command)
	t, err := text.Parse(string(helptext))

	if err != nil {
		c.ServicesLog(err.Error())
		Log.Print(err)
		return
	}

	if _, ok := c.Commands[command]; !ok {
		if command != "help" {
			c.Notice(target, "No such command")
		}
	}

	data := struct {
		Me     *Client
		Target *Client
	}{
		c,
		target,
	}

	t.Execute(buffer, data)

	output := buffer.String()

	for _, line := range strings.Split(output, "\n") {
		c.Notice(target, line)
	}
}
