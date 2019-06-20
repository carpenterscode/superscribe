package superscriber

type Env string

const (
	Sandbox Env = "Sandbox"
	Prod    Env = "PROD"
)

type NoteType string

const (
	Cancel               NoteType = "CANCEL"
	DidChangeRenewalPref NoteType = "DID_CHANGE_RENEWAL_PREF"
	InitialBuy           NoteType = "INITIAL_BUY"
	InteractiveRenewal   NoteType = "INTERACTIVE_RENEWAL"
	Renewal              NoteType = "RENEWAL"

	// Introduced in June 2019 at WWDC
	DidChangeRenewalStatus NoteType = "DID_CHANGE_RENEWAL_STATUS"
)
