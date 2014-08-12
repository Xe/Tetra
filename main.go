/*
Command Tetra is an extended services package for TS6 IRC daemons with Lua and
Moonscript support.
*/
package main

import (
	"fmt"
	"os"

	"github.com/Xe/Tetra/bot"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "cqbot"
	app.Usage = "Kicks off the cqbot irc bot"
	app.Version = "0.1-dev"
	app.Author = "Sam Dodrill <xena@yolo-swag.com>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "conf, c",
			Value:  "etc/config.yaml",
			Usage:  "Configuration file to use. This can be overridden by TETRA_CONFIG_PATH.",
			EnvVar: "TETRA_CONFIG_PATH",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "run",
			Usage:       "runs the bot",
			Action:      startBot,
			Description: "Kick off the bot and all of its associated goroutines.",
		},
		{
			Name:        "checkconfig",
			Usage:       "checks config validity",
			Description: "Checks a configuration file for validity or not",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "conf, c",
					Value:  "etc/config.yaml",
					Usage:  "Configuration file to use. This can be overridden by TETRA_CONFIG_PATH.",
					EnvVar: "TETRA_CONFIG_PATH",
				},
			},
			Action: func(c *cli.Context) {
				_, err := tetra.NewConfig(c.String("conf"))
				if err != nil {
					fmt.Printf("File %s not readable: %s\n", c.String("conf"), err.Error())
					os.Exit(1)
				}

				fmt.Printf("Config file %s is valid.\n", c.String("conf"))
			},
		},
		{
			Name:        "genconfig",
			Usage:       "generates a config",
			Description: "Generates a config file based on arguments or environment variables",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "conf, c",
					Value:  "etc/config.yaml",
					Usage:  "Configuration file to use. This can be overridden by TETRA_CONFIG_PATH.",
					EnvVar: "TETRA_CONFIG_PATH",
				},
				cli.StringFlag{
					Name:   "server, s",
					Value:  "127.0.0.1",
					Usage:  "Remote server to connect to. This can be overriden by TETRA_HOST.",
					EnvVar: "TETRA_HOST",
				},
				cli.StringFlag{
					Name:  "port, p",
					Value: "6667",
					Usage: "The port the ircd is listening on.",
				},
				cli.StringFlag{
					Name:  "password, p",
					Value: "shameless",
					Usage: "Server password.",
				},
				cli.BoolFlag{
					Name:  "ssl",
					Usage: "Toggles SSL on the uplink Connection.",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "tetra.int",
					Usage: "Server name to send over TS6.",
				},
				cli.StringFlag{
					Name:  "sid",
					Value: "326",
					Usage: "Server ID to send over TS6.",
				},
				cli.StringFlag{
					Name:  "prefix, p",
					Value: "`",
					Usage: "Command prefix to use in channel.",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "Enables the debug flag in the config.",
				},
			},
			Action: func(c *cli.Context) {
				fmt.Println("Writing config file to " + c.String("conf"))
				if _, err := os.Stat(c.String("conf")); err != nil {
					file, err := os.Create(c.String("conf"))
					if err != nil {
						panic(err)
					}

					file.Close()
				}

				config, err := tetra.NewConfig(c.String("conf"))

				config.Uplink.Host = c.String("server")
				config.Uplink.Port = c.String("port")
				config.Uplink.Password = c.String("password")
				config.Uplink.Ssl = c.Bool("ssl")
				config.General.Debug = c.Bool("debug")
				config.General.Prefix = c.String("prefix")
				config.Server.Name = c.String("name")
				config.Server.Sid = c.String("sid")

				config.Autoload = []string{
					"tetra/dispatch",
					"chatbot/fantasy",
					"tetra/load",
					"tetra/scripts",
					"tetra/unload",
					"tetra/die",
					"tetra/version",
				}

				file, err := os.Create(c.String("conf"))
				if err != nil {
					panic(err)
				}

				output, err := yaml.Marshal(config)
				if err != nil {
					panic(err)
				}

				fmt.Fprint(file, output)
				file.Close()

				fmt.Println(output)

				fmt.Println("Wrote config to " + c.String("conf"))
			},
		},
	}

	app.Action = startBot

	app.Run(os.Args)
}

func startBot(c *cli.Context) {
	confloc := c.String("conf")

	if _, err := os.Open(confloc); err != nil {
		fmt.Fprintln(os.Stderr, "Please either set TETRA_CONFIG_PATH to the location of the configuration file or add your config at etc/config.yaml")
		os.Exit(1)
	}

	fmt.Printf("Using config file %s\n", confloc)

	bot := tetra.NewTetra(confloc)

	bot.Connect(bot.Config.Uplink.Host, bot.Config.Uplink.Port)
	defer bot.Conn.Conn.Close()

	bot.Auth()
	bot.StickConfig()
	bot.WebApp()

	bot.Main()
}
