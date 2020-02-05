package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	obsws "github.com/christopher-dG/go-obs-websocket"
	"github.com/spf13/cobra"
)

var (
	configFile string
	host       string
	port       uint32

	config Config

	rootCmd = &cobra.Command{
		Use:   "obs-scene-switcher",
		Short: "obs-scene-switcher tracks your active window and switches scenes accordingly",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute()
		},
	}

	client        *obsws.Client
	recentWindows []Window
)

type Scene struct {
	SceneName   string `toml:"scene_name"`
	WindowClass string `toml:"window_class"`
	WindowName  string `toml:"window_name"`
}

type Scenes []Scene

type Config struct {
	Scenes     Scenes `toml:"scenes"`
	AwayScenes Scenes `toml:"away_scenes"`
}

func LoadConfig(filename string) (Config, error) {
	config := Config{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	_, err = toml.Decode(string(b), &config)
	return config, err
}

func main() {
	var err error
	config, err = LoadConfig(configFile)
	if err != nil {
		fmt.Println("could not load config file:", configFile)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if client != nil {
		client.Disconnect()
	}
}

func init() {
	cobra.OnInitialize(connectOBS)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "scenes.toml", "path to config file")
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "host to connect to")
	rootCmd.PersistentFlags().Uint32VarP(&port, "port", "p", 4444, "port to connect to")
}

func handleActiveWindowChanged(event ActiveWindowChangedEvent) {
	fmt.Println(fmt.Sprintf("Active window changed to %s (%d, %s)",
		event.Window.Class, event.Window.ID, event.Window.Name))

	// remove dupes
	i := 0
	for _, rw := range recentWindows {
		if rw.ID == event.Window.ID {
			continue
		}

		recentWindows[i] = rw
		i++
	}
	recentWindows = recentWindows[:i]

	recentWindows = append([]Window{event.Window}, recentWindows...)
	if len(recentWindows) > 15 {
		recentWindows = recentWindows[0:15]
	}

	req := obsws.NewGetCurrentSceneRequest()
	resp, err := req.SendReceive(*client)
	for _, v := range config.AwayScenes {
		if err == nil && resp.Name == v.SceneName {
			fmt.Println("Skipping switch, in Away mode!")
			return
		}
	}

	for _, v := range config.Scenes {
		if event.Window.Class == v.WindowClass ||
			event.Window.Name == v.WindowName {
			req := obsws.NewSetCurrentSceneRequest(v.SceneName)
			req.Send(*client)
		}
	}
}

func handleWindowClosed(event WindowClosedEvent) {
	i := 0
	for _, rw := range recentWindows {
		if rw.ID == event.Window.ID {
			continue
		}

		recentWindows[i] = rw
		i++
	}
	recentWindows = recentWindows[:i]
}

func execute() error {
	x := Connect(os.Getenv("DISPLAY"))
	defer x.Close()

	tch := make(chan interface{})
	x.TrackWindows(tch, time.Second)

	for e := range tch {
		switch event := e.(type) {
		case ActiveWindowChangedEvent:
			handleActiveWindowChanged(event)

		case WindowClosedEvent:
			handleWindowClosed(event)
		}
	}

	return nil
}

func connectOBS() {
	// disable obsws logging
	obsws.Logger = log.New(ioutil.Discard, "", log.LstdFlags)

	client = &obsws.Client{Host: host, Port: int(port)}
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set the amount of time we can wait for a response.
	obsws.SetReceiveTimeout(time.Second * 2)
}
