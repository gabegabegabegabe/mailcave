# mailcave

For when you just don't feel safe stashing your mail anywhere that's not a cave.

mailcave is a REST application that provides an interface for the storing of MIME email messages.  Currently, the only database it supports is MongoDB.

## Getting Started

### Requirements
* Golang installed
* MongoDB installed
* Docker installed (if running in  a Docker container)

## Running
### Command Line Arguments
* `--dbAddr` - the address of the database - default is `mongodb://localhost:27017`
* `--dbName` - the name of the archive database - default is `mailcave`
* `--ipAddr` - the address on which mailcave should listen on - default is `:8080`
* `--logToStdOut` - whether or not mailcave should log to stdout in addition to file - default is `true`

#### Example
`cd mailcave/cmd/mailcave`<br />
`go build -o mailcave main.go`<br />
`./mailcave --dbName=mailcaveDb --ipAddr=:8088 --logToStdOut=false`

#### With Docker
`cd mailcave`<br />
`docker build -t mailcave .`<br />
`docker run --publish 8080:8080 --name mailcaveContainer --rm mailcave`

## Currently Supported Endpoints
### Store a Message
* Method - POST
* URL - `/store`
* Parameters
  * `mime_message=[mime message content]`

## Testing
A command line app called `mailsend` is provided to allow you to send a MIME message contained in a text file to a running instance of mailcave.

### Example
`cd mailcave/cmd/mailsend`<br />
`go build -o mailsend main.go`<br />
`./mailsend --mimeFile=mime.msg --mailcaveAddr=http://localhost:8088`
