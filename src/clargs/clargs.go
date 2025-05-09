package clargs

import (
	"flag"
	"fmt"
)

type CommandLineArgs struct {
	ApiKey       string
	Prompt       string
	SystemPrompt string
	Model        string
	File         string
	WebSearch    bool
	Help         bool
}

func ParseCommandLineArgs(args []string) (CommandLineArgs, error) {
	flag.Usage = func() {
		writer := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		preamble := `chat

A simple command line tool to interact with OpenAI's large language models.

EXAMPLE:
Performs simple prompt.

  chat "What is the capital of France?"

Performs a prompt with a file provided as context.

  chat -f chapter.txt "Summarize the chapter."

USAGE:
chat [<options>] <prompt>
`

		fmt.Fprintln(writer, preamble)
		flag.PrintDefaults()

	}

	clArgs := CommandLineArgs{}

	flag.StringVar(&clArgs.ApiKey, "k", "", "OpenAI API key.")
	flag.StringVar(&clArgs.Model, "m", "", "OpenAI model to use.")
	flag.StringVar(&clArgs.SystemPrompt, "s", "", "System prompt to provide to the model.")
	flag.StringVar(&clArgs.File, "f", "", "A file provided to the model as context. Can also include line specifiers such as <file path>:300 or <file path>:300-400. Can include multiple file paths separated by commas (,).")
	flag.BoolVar(&clArgs.WebSearch, "w", false, "Allows the model to perform a web search as part of its query.")
	flag.BoolVar(&clArgs.Help, "h", false, "Prints help text.")

	// Parse the flags from the provided args slice
	err := flag.CommandLine.Parse(args[1:])
	if err != nil {
		return clArgs, err
	}

	// Collect remaining non-flag arguments as the prompt
	remainingArgs := flag.CommandLine.Args()
	if len(remainingArgs) == 0 {
		return clArgs, fmt.Errorf("prompt is required")
	}
	clArgs.Prompt = remainingArgs[0]

	return clArgs, nil
}
