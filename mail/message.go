package mail

import (
	"fmt"
	"net/mail"
)

const (
	toKey        = "To"
	fromKey      = "From"
	subjectKey   = "Subject"
	messageIDKey = "Message-ID"

	dateFormatPrefix = "Mon, 02 Jan 2006 15:04:05 "
)

// Msg represents an email message.
type Msg struct {
	To        []*mail.Address
	From      []*mail.Address
	Date      string
	Subject   string
	MessageID string
}

// MsgFromMIMEMessage creates a Msg by extracting the relevant fields from a MIME mail.Message.
func MsgFromMIMEMessage(mimeMsg *mail.Message) (*Msg, error) {

	toList, err := mimeMsg.Header.AddressList(toKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse To value from MIME header: %s", err)
	}

	fromList, err := mimeMsg.Header.AddressList(fromKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse From value from MIME header: %s", err)
	}

	date, err := mimeMsg.Header.Date()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Date from MIME header: %s", err)
	}
	zone, _ := date.Zone()
	dateStr := date.Format(dateFormatPrefix + zone)

	subject := mimeMsg.Header.Get(subjectKey)

	messageID := mimeMsg.Header.Get(messageIDKey)
	if len(messageID) == 0 {
		return nil, fmt.Errorf("could not find Message-ID in message with subject '%s'", subject)
	}

	return &Msg{
		To:        toList,
		From:      fromList,
		Date:      dateStr,
		Subject:   subject,
		MessageID: messageID,
	}, nil
}
