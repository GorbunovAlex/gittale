# gittale

`gittale` – from Tale (German: story) and hints at tala (Norse: speak) – is a command-line tool that wraps common `git` commands and leverages a Gemini API to generate concise commit messages based on your staged changes.

## Features

- Run any `git` command via `gittale`.
- For `git commit`, automatically generates a commit message using an LLM based on your staged diff.
- Simple integration with local Ollama model.

## Prerequisites

- Go (for building the app)
- `git` installed and available in your `PATH`
- Access to a local Ollama instance running your preferred model
- (Optional) `.env` file to configure your Ollama model (e.g., `OLLAMA_MODEL=your-model`)

## Installation

Ensure you set your Ollama model in a `.env` file as `OLLAMA_MODEL=your-model`.

```sh
git clone git@github.com:GorbunovAlex/gittale.git
cd gittale
go build -o gittale
```

## Usage

Run git commands as follows:

```sh
./gittale status
./gittale branch
./gittale commit
# Optionally, move the binary to /usr/local/bin for global usage:
sudo mv gittale /usr/local/bin/
```

- For `commit`, the app will generate a commit message using the LLM via your local Ollama instance, based on your staged changes.

## License

Apache License
