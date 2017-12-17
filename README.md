# mailcave

For when you just don't feel safe stashing your mail anywhere that's not a cave.

mailcave is a REST application that provides an interface for the storing of MIME email messages.

## Currently Supported Endpoints
### Store a Message
* Method - POST
* URL - `/store`
* Parameters
  * `mime_message=[mime message content]`
