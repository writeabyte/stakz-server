# Server Repository for stakz.dev

What is it? Stakz-server is the backend for [stakz.dev](https://stakz.dev)

It is a basic automation and file server that allows users to run commands from stakz.dev, save files, and browse the file system.

## Installation

#### Using Homebrew

```bash
brew install writeabyte/stakz/stakz ```

#### Using Go
```bash
go install github.com/writeabyte/stakz-server@latest 
```
Create an alias for stakz-server in your shell profile
```bash
alias stakz="stakz-server"
```

#### From Releases
Download the latest release from the [releases page](https://github.com/writeabyte/stakz-server/releases) and move the binary to your $PATH


## How does it work?

run `stakz` to start the server. The server will start on port 3001 by default and will serve files from the current directory.

By default, the server will not execute commands. To enable command execution, run `stakz --execute`.

```
stakz --help
  -dir string
        The directory you want the stakz server to run in. (default ".")
  -execute
        Enable the /execute endpoint allowing the server to run commands.
  -key string
        The server key used to authenticate requests. 
        If not set, a random key will be generated.
  -keyEnabled
        Whether or not to require a server key for requests. 
        If false, the server key will be ignored and requests will not be authenticated. 
        Only do this if you trust the execution context! (e.g. running in a container) (default true)
  -port int
        The port the server will listen on. (default 3001)
```
