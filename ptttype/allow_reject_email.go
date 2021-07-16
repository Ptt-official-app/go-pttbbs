package ptttype

import "strings"

type AllowRejectEmailOp byte

const (
	ALLOW_REJECT_EMAIL_OP_A AllowRejectEmailOp = 'A'
	ALLOW_REJECT_EMAIL_OP_P AllowRejectEmailOp = 'P'
	ALLOW_REJECT_EMAIL_OP_S AllowRejectEmailOp = 'S'
	ALLOW_REJECT_EMAIL_OP_D AllowRejectEmailOp = 'D'

	ALLOW_REJECT_EMAIL_OP_PERCENT AllowRejectEmailOp = '%'
)

type AllowRejectEmail struct {
	Op           AllowRejectEmailOp
	Pattern      string
	LowerPattern string
}

func NewAllowRejectEmail(line string) *AllowRejectEmail {
	if len(line) < 2 {
		return nil
	}
	op := AllowRejectEmailOp(line[0])
	pattern := strings.TrimSpace(line[1:])
	lowerPattern := strings.ToLower(pattern)

	return &AllowRejectEmail{
		Op:           op,
		Pattern:      pattern,
		LowerPattern: lowerPattern,
	}
}

func (a *AllowRejectEmail) IsValid(email string) (isValid bool, err error) {
	lowerEmail := strings.ToLower(email)
	switch a.Op {
	case ALLOW_REJECT_EMAIL_OP_A: // equal
		return strings.EqualFold(email, a.Pattern), nil
	case ALLOW_REJECT_EMAIL_OP_P: // pattern
		return strings.Contains(lowerEmail, a.LowerPattern), nil
	case ALLOW_REJECT_EMAIL_OP_S: //@domain
		theIdx := strings.Index(email, "@")
		if theIdx == -1 {
			return false, nil
		}
		return strings.EqualFold(email[(theIdx+1):], a.Pattern), nil
	case ALLOW_REJECT_EMAIL_OP_D: // domain
		if len(email) <= len(a.Pattern) {
			return false, nil
		}

		domainIdx := len(email) - len(a.Pattern)
		if !strings.EqualFold(email[domainIdx:], a.Pattern) {
			return false, nil
		}
		atIdx := domainIdx - 1
		atChar := email[atIdx]
		return atChar == '.' || atChar == '@', nil
	case ALLOW_REJECT_EMAIL_OP_PERCENT: // all
		return true, nil
	default:
		return false, ErrInvalidAllowRejectEmailOp
	}
}
