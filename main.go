/*
Command Tetra is an extended services package for TS6 IRC daemons with Lua and
Moonscript support.
*/
package main

import (
	"fmt"
	"os"

	"code.google.com/p/go-uuid/uuid"
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
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enables the debug flag in the config.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "run",
			Usage:       "runs the bot",
			Action:      startBot,
			Description: "Kick off the bot and all of its associated goroutines.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "conf, c",
					Value:  "etc/config.yaml",
					Usage:  "Configuration file to use. This can be overridden by TETRA_CONFIG_PATH.",
					EnvVar: "TETRA_CONFIG_PATH",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "Enables the debug flag in the config.",
				},
			},
		},
		{
			Name:        "checkconfig",
			Usage:       "checks config validity",
			Description: "Checks a configuration file for validity or not",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "conf, c",
					Value:  "etc/config.yaml",
					Usage:  "configuration file to use",
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
					Usage:  "Configuration file to use.",
					EnvVar: "TETRA_CONFIG_PATH",
				},
				cli.StringFlag{
					Name:   "server, s",
					Value:  "127.0.0.1",
					Usage:  "Remote server to connect to.",
					EnvVar: "TETRA_HOST",
				},
				cli.StringFlag{
					Name:   "port, p",
					Value:  "6667",
					Usage:  "The port the ircd is listening on.",
					EnvVar: "TETRA_PORT",
				},
				cli.StringFlag{
					Name:   "password",
					Value:  "shameless",
					Usage:  "Server password.",
					EnvVar: "TETRA_PASSWORD",
				},
				cli.BoolFlag{
					Name:  "ssl",
					Usage: "Toggles SSL on the uplink Connection.",
				},
				cli.StringFlag{
					Name:   "name, n",
					Value:  "tetra.int",
					Usage:  "Server name to send over TS6.",
					EnvVar: "TETRA_SERVER_NAME",
				},
				cli.StringFlag{
					Name:   "sid",
					Value:  "326",
					Usage:  "Server ID to send over TS6.",
					EnvVar: "TETRA_SID",
				},
				cli.StringFlag{
					Name:   "prefix",
					Value:  "`",
					Usage:  "Command prefix to use in channel.",
					EnvVar: "TETRA_PREFIX",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "Enables the debug flag in the config.",
				},
				cli.StringFlag{
					Name:   "gecos, g",
					Usage:  "Server GECOS",
					Value:  "Tetra Services",
					EnvVar: "TETRA_SERVER_GECOS",
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

				config := &tetra.Config{
					Uplink: &tetra.UplinkConfig{
						Host:     c.String("server"),
						Port:     c.String("port"),
						Password: c.String("password"),
						Ssl:      c.Bool("ssl"),
					},
					General: &tetra.GeneralConfig{
						StaffChan: "#opers",
						SnoopChan: "#services",
						Debug:     c.Bool("debug"),
						Prefix:    c.String("prefix"),
					},
					Server: &tetra.ServerConfig{
						Gecos: c.String("gecos"),
						Name:  c.String("name"),
						Sid:   c.String("sid"),
					},
					Services: []*tetra.ServiceConfig{
						&tetra.ServiceConfig{
							Name:   "tetra",
							Nick:   "Tetra",
							User:   "tetra",
							Host:   c.String("name"),
							Gecos:  "Tetra admin client",
							Certfp: uuid.New(),
						},
						&tetra.ServiceConfig{
							Name:   "chatbot",
							Nick:   "Chatbot",
							User:   "chatbot",
							Host:   c.String("name"),
							Gecos:  "Chitty chatter bot!",
							Certfp: uuid.New(),
						},
					},
					Stats: &tetra.StatsConfig{
						Host: "NOCOLLECTION",
					},
					Autoload: []string{
						"tetra/upgrade",
						"tetra/load",
						"tetra/scripts",
						"tetra/unload",
						"tetra/die",
						"tetra/version",
						"chatbot/btc",
						"chatbot/sendfile",
						"chatbot/youtube",
						"chatbot/source",
						"chatbot/doge",
						"chatbot/tell",
					},
				}

				file, err := os.Create(c.String("conf"))
				if err != nil {
					panic(err)
				}

				output, err := yaml.Marshal(config)
				if err != nil {
					panic(err)
				}

				plaintext := string(output)

				fmt.Fprint(file, plaintext)
				file.Close()

				fmt.Println(plaintext)

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
