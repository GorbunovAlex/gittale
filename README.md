# gittale

`gittale` – from Tale (German: story) and hints at tala (Norse: speak) – is a command-line tool that wraps common `git` commands and leverages a Gemini API to generate concise commit messages based on your staged changes.

## Features

- Run any `git` command via `gittale`.
- For `git commit`, automatically generates a commit message using an LLM based on your staged diff.
- Simple integration with Gemini API.

## Prerequisites

- Go (for building the app)
- `git` installed and available in your `PATH`
- Access to Gemini API and an API key
- (Optional) `.env` file to configure your API key

## Installation

```sh
git clone git@github.com:GorbunovAlex/gittale.git
cd gittale
go build -o gittale
```

## Usage

1. Ensure you have your Gemini API key set in the environment variable `API_KEY` (or in a `.env` file).
2. Run git commands as follows:

```sh
./gitgpt status
./gitgpt branch
./gitgpt commit
# Optionally, move the binary to /usr/local/bin for global usage:
sudo mv gittale /usr/local/bin/
```

- For `commit`, the app will generate a commit message using the LLM based on your staged changes.

## Configuration

You can configure `gittale` using the following environment variable:

- `API_KEY`: Your Gemini API key.

Alternatively, you can create a `.env` file in the same directory as the `gittale` executable with the following format:

```dotenv
API_KEY=your-gemini-api-key
```

Environment variables take precedence over values in the `.env` file.

## License

Apache License
