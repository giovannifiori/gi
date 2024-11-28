# gi

A CLI tool built with [Go](https://github.com/golang/go) and [Cobra](https://github.com/spf13/cobra) to create `.gitignore` files.

## Description

`gi` is a command-line tool that helps you create `.gitignore` files for your projects. It provides a simple interface to add ignore patterns, making it easier to keep your repository clean and organized.
All the contents are from the [Toptal API](https://docs.gitignore.io/use/api).

## Installation

To install `gi`, you need to have Go installed on your machine. Then, you can use the following command to install the tool:

```sh
go install github.com/giovannifiori/gi@latest
```

## Usage

Here are some examples of how to use `gi`:

### Generate a `.gitignore` file for a Go project

```sh
gi go
```

### Generate a `.gitignore` file for a Node project in a MacOS environment

```sh
gi node macos
```

### Generate a `.gitignore` file choosing from a list of options

```sh
gi
```

## Contributing

We welcome contributions to `gi`. If you would like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'feat: add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
