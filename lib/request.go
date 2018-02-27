package custody



//CreationRequest: contains the information necessary to request the creation of a new user.
type CreationRequest struct {
	Name 		string
	PublicKey   []byte
}

//RecordRequest: contains the information necessary to request the validation of a new message.
type RecordRequest struct {
	Name string
	Data []byte
	Hash []byte
}

//ListRequest: the filter on the ledger you want to apply.
type ListRequest struct {
	Name string
}
