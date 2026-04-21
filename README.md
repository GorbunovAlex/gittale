# gittale

gittale is a command-line tool that wraps common git commands and can generate commit messages from staged changes using an LLM.

## Features

- Run any git command via gittale.
- For commit, automatically generate a commit message from staged diff.
- Provider-based LLM integration (Gemini and Ollama).
- Diff batching to avoid overflowing LLM context windows.

## Prerequisites

- Go
- git installed and available in PATH
- One LLM provider:
  - Gemini API key, or
  - local/remote Ollama endpoint with a model

## Configuration

Copy the example config and fill in your provider details before building:

```sh
cp config/example.yaml config/config.yaml
```

```yaml
env: "local"            # local, dev, prod
model_provider: "ollama" # ollama | gemini
gemini_api_key: ""
gemini_model: "gemini-2.0-flash"
ollama_model: "llama3.1"
ollama_url: "http://localhost:11434"
```

Provider notes:

- `model_provider: gemini` requires `gemini_api_key`.
- `model_provider: ollama` requires `ollama_model`.
- `gemini_model` and `ollama_url` have sensible defaults.

All values are read from the YAML file **at build time** by a Go tool (`cmd/configgen`) and baked into the binary via `-ldflags`. No config file is needed at runtime.

## Installation

```sh
git clone git@github.com:GorbunovAlex/gittale.git
cd gittale
cp config/example.yaml config/config.yaml
# edit config/config.yaml with your provider and keys
make install
```

`make install` reads `config/config.yaml`, compiles all values into the binary, and installs it to `/usr/local/bin`. After that, `gittale` works from anywhere with no environment variables.

| Target | Description |
|--------|-------------|
| `make build` | Build binary locally using `config/config.yaml` |
| `make install` | Build and install to `/usr/local/bin` |
| `make uninstall` | Remove the installed binary |

Use a different config file:

```sh
make install CONFIG_FILE=/path/to/other.yaml
```

## Usage

Run any git command through gittale:

```sh
./gittale status
./gittale branch
./gittale commit
```

Commit flow:

1. Stage your changes with git add.
2. Run ./gittale commit.
3. gittale reads staged diff, splits it into batches, summarizes each batch, and generates final commit message.
4. Commit is executed as git commit -m "<generated message>".

Pass-through flow:

- Commands other than commit are forwarded to git as-is.

## License

Apache License
