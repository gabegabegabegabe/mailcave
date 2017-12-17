package mail

// Archive provides storage of mail messages.
type Archive interface {
	Open() error
	Close()
	ArchiveMessage(msg *Msg) (string, error)
}
