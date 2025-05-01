# Chat

Chat is a simple command line tool for interacting with OpenAI's ChatGPT. Chat
is designed to make it simple to send a quick prompt to ChatGPT without
requiring opening a browser.

## Features

- Simple interface.
- Global configuration or per-directory configuration.
- Provide a file or part of a file as context.
- Override configuration through command line arguments.

## Usage and Examples
Usage:

`chat [<options>] <prompt>`

Simple prompt.

`chat "What is the capital of France?"`

Prompt with file context. Can be limited to single line or range.

`chat -f code.c "Summarize this code."`  
`chat -f code.c:25 "Explain this line of code."`  
`chat -f code.c:38-45 "Explain this function."`

## Options
Command line options can be provided to modify the behavior of Chat. Command
line options ALWAYS override configuration. See the "Configuration" section
below to understand how configuration works.

| Option | Description|
| :-- | :-- |
| `-h` | Shows help text. |
| `-k` | OpenAI API key. Required if not provided in configuration. |
| `-f <path>[:<start>[-<end>]]` | Add a file or part of a file as context. |
| `-m <model>` | Model to use. |
| `-s <prompt>` | System prompt to use. |

## Configuration
A configuration can be provided to set defaults for the model, system prompt,
and API key. Configuration can be created as a combination of global and local
configurations.

### Global Configuration
The global configuration should be placed in `<HOME>/.chatconfig` where `<HOME>`
is the users home directory.

- **Windows**: `%USERPROFILE%`
- **Mac OS**: `$HOME` or `~`
- **Linux**: `$HOME` or `~`

### Local Configuration
A `.chatconfig` file can be placed in any directory. Chat will automatically
merge all configurations found in the current working directory and all parent
directories. The closer the configuration is to the working directory, the more
priority it will have (overwriting keys from lower priority configurations.)

The global configuration is always loaded last with the least priority.

### Configuration File

The `.chatconfig` file should be a simple text file with a `KEY=VALUE` pair per
line. The following configuration options are currently available:

| Option | Description|
| :-- | :-- |
| `CHAT_API_KEY` | OpenAI API key. |
| `CHAT_MODEL` | Model to use. |
| `CHAT_SYSTEM_PROMPT` | System prompt to use. |

#### Example Configuration File
```
CHAT_API_KEY=<API Key>
CHAT_MODEL=4o-mini
CHAT_SYSTEM_PROMPT=You are a command line assistant. You answer questions and response in single line responses if possible.
```

### Environment Variable Configuration

All configuration keys can also be set as environment variables. Configuration
from environment variables will have the least priority.

## Installation

### Portable
Chat is a single file executable and is fully self-contained. You can place the
binary on a USB-stick or wherever you want to run it from and call the
executable.

### Windows
Download the `chat.exe` binary and place in a location of your choice. For
example `C:\Program Files\Chat\chat.exe`. Then add the directory to your path.

### Mac OS or Linux
Download the `chat` binary and place it in `/usr/local/bin`. Then make the
binary exectuable by calling: 

`chmod +x /usr/local/bin/chat`