# Discord_Shell

## What is it?
The most ridiculous, unpractical and insecure way to get command line access to your computer.
Or in different words, a Discord bot that redirects in and output of a program running on your computer to Discord.
## How can I get it working?
Get yourself a Discord Bot and copy its Token to the second line of the main function where it says Token.
Get yourself the golang compiler, change to the directory you cloned this repository to and type
```bash
go run main.go
```
## What can I do with it?
You can either run simple commands that do not require input in the form:
```bash
command arg arg arg
```
or when your command is interactive and requires input you use the form:
```bash
[i] command arg arg arg
```
Using the interactive form you can even have a working bash session in Discord.

## Issues
Python interpreter is not working.