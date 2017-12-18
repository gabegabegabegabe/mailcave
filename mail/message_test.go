package mail

import (
	"fmt"
	"net/mail"
	"reflect"
	"strings"
	"testing"
)

func TestMsgFromMIMEMessage(t *testing.T) {

	mimeMissingFrom, err := mail.ReadMessage(strings.NewReader(mimeStrMissingFrom))
	if err != nil {
		panic(err)
	}

	validMIME, err := mail.ReadMessage(strings.NewReader(mimeStrValid))
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		mime        *mail.Message
		expectedMsg *Msg
		expectedErr error
	}{
		{
			mime:        mimeMissingFrom,
			expectedMsg: nil,
			expectedErr: fmt.Errorf("failed to parse From value from MIME header with error 'mail: header not in message'"),
		},
		{
			mime: validMIME,
			expectedMsg: &Msg{
				To: []*mail.Address{
					&mail.Address{
						Name:    "",
						Address: "specialrecipient@mailcave.com",
					},
				},
				From: []*mail.Address{
					&mail.Address{
						Name:    "",
						Address: "mail.place@company.com",
					},
				},
				Date:      "Fri, 01 Apr 2011 02:57:21 GMT",
				Subject:   "Some Information",
				MessageID: "<1479419471.1301626641534.CaveMail.pjfxbg1@eijfp99>",
			},
			expectedErr: nil,
		},
	}

	for _, c := range testCases {
		msg, err := MsgFromMIMEMessage(c.mime)

		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expcted err to be %q but was found to be %q", c.expectedErr, err)
		}

		if !reflect.DeepEqual(msg, c.expectedMsg) {
			t.Errorf("Expcted msg to be %q but was found to be %q", c.expectedMsg, msg)
		}
	}
}

const (
	mimeStrMissingFrom = `Return-Path: <mail.place@company.com>
Delivered-To: insurance@localhost.mailcave.com
Date: Fri, 1 Apr 2011 02:57:21 GMT
Message-ID: <1479419471.1301626641534.CaveMail.pjfxbg1@eijfp99>
To: specialrecipient@mailcave.com
Subject: Some Information
Mime-Version: 1.0
Content-Type: multipart/alternative; 
	boundary="----=_Part_26607_1527703119.1301626641533"
urn:schemas:mailheader:content-type: multipart/mixed

------=_Part_26607_1527703119.1301626641533
Content-Type: text/plain; charset=ISO-8859-1
Content-Transfer-Encoding: binary

Dear Fakename,

We would like to update you on your account balance:

Date: 10/24/29
Account no.: XXXXXXXXXX  
Account balance: USD 0.00 

Sorry about the poorness.  Thank you for using the internet.

Company Online
Fancier Company Name

Please do not reply to this email. We are not fond of you.
------=_Part_26607_1527703119.1301626641533--

`

	mimeStrValid = `Return-Path: <mail.place@company.com>
Delivered-To: insurance@localhost.mailcave.com
Date: Fri, 1 Apr 2011 02:57:21 GMT
Message-ID: <1479419471.1301626641534.CaveMail.pjfxbg1@eijfp99>
From: mail.place@company.com
To: specialrecipient@mailcave.com
Subject: Some Information
Mime-Version: 1.0
Content-Type: multipart/alternative; 
	boundary="----=_Part_26607_1527703119.1301626641533"
urn:schemas:mailheader:content-type: multipart/mixed

------=_Part_26607_1527703119.1301626641533
Content-Type: text/plain; charset=ISO-8859-1
Content-Transfer-Encoding: binary

Dear Fakename,

We would like to update you on your account balance:

Date: 10/24/29
Account no.: XXXXXXXXXX  
Account balance: USD 0.00 

Sorry about the poorness.  Thank you for using the internet.

Company Online
Fancier Company Name

Please do not reply to this email. We are not fond of you.
------=_Part_26607_1527703119.1301626641533--

`
)
