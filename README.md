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

Ensure you set your Gemini API key in a `.env` file as `API_KEY=YOUR_API_KEY`.

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

- For `commit`, the app will generate a commit message using the LLM based on your staged changes.

## License

Apache License
