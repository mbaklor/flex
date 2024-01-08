# Flex - A dev tool for Barix Flexa

## About

This command line tool is meant to help you initialize flexa packages by scaffolding a simple hello world app, as well as package and send the app and send config files.

## Installation

Place the binary somewhere in your path, and call the app from your command line.

### Release

You can download the binary from the [release page](https://github.com/mbaklor/flex/releases)

### Build from source

To build a binary from source you need git and go installed on your system.
Simply `git clone`, `cd` into the folder, and `go install`

## Usage

### Initialize a package

```
flex init [-n project name] [-l app log file name] [-w show app log in the webui] [-y confirm default log name and webui settings] [-g git init in project]
```

The init command can also be called with no arguments and will prompt you for them.

### Zip and send package to a device

```
flex package [-d directory name] [-b bundle in zip without sending to device] [-a -u -p address username and password of device] [-f json file containing device info, can be multiple]
```

If no directory name is given, cwd is assumed.
If no device information is supplied, you get prompted if you want to bundle without sending.

### Send config file

```
flex config [-a -u -p address username and password of device] [-f json file containing device info, can be multiple] (config file name eg: config.json)
```
