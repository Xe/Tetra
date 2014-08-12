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

	if command == "help" || command == "_index" {
		service = "tetra"
	}

	text := template.New(service + "/" + command)
	helptext, _ := ioutil.ReadFile("help/" + service + "/" + command)
	t, err := text.Parse(string(helptext))

	if err != nil {
		c.ServicesLog(err.Error())
		c.tetra.Log.Print(err)
		return
	}

	data := struct {
		Me     *Client
		Tetra  *Tetra
		Target *Client
	}{
		c,
		c.tetra,
		target,
	}

	t.Execute(buffer, data)

	output := buffer.String()

	for _, line := range strings.Split(output, "\n") {
		c.Notice(target, line)
	}
}
