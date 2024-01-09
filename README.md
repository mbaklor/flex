# Flex - A dev tool for Barix Flexa

## About

This command line tool is meant to help you initialize flexa packages by scaffolding a simple hello world app, as well as package and send the app and send config files.

## Installation

Place the binary somewhere in your path, and call the app from your command line.

### Release

You can download the binary from the [release page](https://github.com/mbaklor/flex/releases)

### Go

You can install this tool with go using the following command:

```
go install github.com/mbaklor/flex
```

note: you may need to add `~/go/bin` to your path to access the tool after install

### Build from source

To build a binary from source you need git and go installed on your system.
Simply `git clone`, `cd` into the folder, and `go install`

note: you may need to add `~/go/bin` to your path to access the tool after install

## Usage

### Initialize a package

```
flex init [-n] [-l] [-w] [-y] [-g]
```

- `-n` - project name
- `-l` - name of app log file, default is `app_log.log`
- `-w` - use this flag to add the app log as a tab in the Flexa web UI
- `-y` - use default options for both `-l` and `-w`
- `-g` - initialize a git repo in the package folder

The init command can also be called with no arguments and will prompt you for them.

### Zip and send package to a device

```
flex package [-d] [-b] [-a -u -p] [-f]
```

- `-d` - directory to zip, if none provided `cwd` is assumed
- `-b` - use this flag to bundle (zip) the package without sending to a device, if provided the next flags are ignored
- `-a` - IP address of Flexa device to send package to
- `-u` - username of Flexa device
- `-p` - password of Flexa device
- `-f` - json file or files containing address, username, and password of Flexa device

If no device information is supplied, you get prompted if you want to bundle without sending.

### Send config file

```
flex config [-a -u -p] [-f] (arg)
```

- `-a` - IP address of Flexa device to send package to
- `-u` - username of Flexa device
- `-p` - password of Flexa device
- `-f` - json file or files containing address, username, and password of Flexa device
- `(arg)` - config file to send, needs to be in json format

### device json example

```json
{
  "address": "1.2.3.4",
  "username": "admin",
  "password": "1234"
}
```
