package main

import (
	"chat/clargs"
	"chat/config"
	"chat/config/config_sources"
	"chat/llm"
	"flag"
	"fmt"
	"os"
)

var configurationSources = []config.ConfigurationSource{
	config_sources.LocalConfigurationSource{},
	config_sources.GlobalConfigurationSource{},
	config_sources.EnvConfigurationSource{},
}

func main() {
	// Load command line arguments.
	commandLineArguments, err := clargs.ParseCommandLineArgs(os.Args)
	if err != nil {
		fmt.Println("Error parsing command line arguments:", err)
		flag.Usage()
		os.Exit(1)
	}

	// Print help text if requested.
	if commandLineArguments.Help {
		flag.Usage()
		os.Exit(0)
	}

	// Load all configurations.
	var configurations []config.Configuration
	for _, configurationSource := range configurationSources {
		configuration, err := configurationSource.Load()
		if err != nil {
			fmt.Println("Error while loading configurations:", err)
			os.Exit(1)
		}

		configurations = append(configurations, configuration)
	}

	// Merge all configurations.
	mergedConfiguration := config.MergeConfigurations(configurations)

	// Construct a chat request by combining the command line arguments and
	// the configuration files.
	chatRequest, err := llm.ConstructChatRequest(commandLineArguments, mergedConfiguration)
	if err != nil {
		fmt.Println("Error constructing chat request:", err)
		os.Exit(1)
	}

	chatResponse, err := llm.Chat(chatRequest)
	if err != nil {
		fmt.Println("Error during chat:", err)
		os.Exit(1)
	}

	fmt.Println(chatResponse.Text)
}
