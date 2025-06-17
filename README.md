# gittale

`gittale` – from Tale (German: story) and hints at tala (Norse: speak) – is a command-line tool that wraps common `git` commands and leverages a local Ollama LLM (such as `gemma3:4b-it-qat`) to generate concise commit messages based on your staged changes.

`gittale` is a command-line tool that wraps common `git` commands and leverages a local Ollama LLM (such as `gemma3:4b-it-qat`) to generate concise commit messages based on your staged changes.

## Features

- Run any `git` command via `gittale`.
- For `git commit`, automatically generates a commit message using an LLM based on your staged diff.
- Simple integration with local Ollama API.

## Prerequisites

- Go (for building the app)
- `git` installed and available in your `PATH`
- [Ollama](https://ollama.com/) running locally with a supported model pulled (e.g., `gemma3:4b-it-qat`)
- (Optional) `.env` file to configure Ollama endpoint and model.

## Installation

```sh
git clone git@github.com:GorbunovAlex/gittale.git
cd gittale
go build -o gittale
```

## Usage

1. Ensure your Ollama server is running locally and you have pulled a model (e.g., `ollama pull gemma3:4b-it-qat`).
2. Configure the Ollama API endpoint and model using environment variables (`OLLAMA_API_URL`, `OLLAMA_MODEL`) or a `.env` file. See the Configuration section below.
3. Run git commands as follows:

```sh
./gittale status
./gittale branch
./gittale commit
# Optionally, move the binary to /usr/local/bin for global usage:
sudo mv gittale /usr/local/bin/
```

- For `commit`, the app will generate a commit message using the LLM based on your staged changes.

## Configuration

You can configure `gittale` using the following environment variables:

- `OLLAMA_API_URL`: The URL of the Ollama API endpoint (default: `http://localhost:11434/api/chat`).
- `OLLAMA_MODEL`: The name of the Ollama model to use (default: `gemma3:4b-it-qat`).

Alternatively, you can create a `.env` file in the same directory as the `gittale` executable with the following format:

```dotenv
OLLAMA_API_URL=http://localhost:11434/api/chat
OLLAMA_MODEL=gemma3:4b-it-qat
```

Environment variables take precedence over values in the `.env` file.

## License

Apache License
