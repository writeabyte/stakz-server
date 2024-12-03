# Server Repository for stakz.dev

What is it? Stakz-server is the backend for [stakz.dev](https://stakz.dev)

It is a basic automation and file server that allows users to run commands from stakz.dev on their local machine or a remote server.

## Installation



## How does it work?

```
Usage of /Users/michael/Development/web/stakz/stakz-server/tmp/main:
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
