package ptttype

import "github.com/Ptt-official-app/go-pttbbs/types"

type MQType int

const (
	MQ_TEXT MQType = iota
	MQ_UUENCODE
	MQ_JUSTIFY
)

var MQTypeString = []string{
	"text",
	"uuencode",
	"mq-justify",
}

func (m MQType) String() string {
	if int(m) < len(MQTypeString) {
		return MQTypeString[m]
	}

	return "unknown"
}

type MailQueue struct {
	FilePath Filename_t
	Subject  Subject_t
	MailTime types.Time4
	Sender   UserID_t
	Username Nickname_t
	RCPT     RCPT_t
	Method   int32
	Niamod   []byte
}
