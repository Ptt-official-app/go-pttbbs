package api

type EmailTokenContext string

const (
	CONTEXT_CHANGE_EMAIL EmailTokenContext = "email"
	CONTEXT_SET_ID_EMAIL EmailTokenContext = "id_email"
)
