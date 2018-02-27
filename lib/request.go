package custody

//RecordRequest: contains the information necessary to request a Clerk operation.
// All operations use the same RecordRequest format, you only need to provide values for the necessary arguments.
type RecordRequest struct {
	Name      string
	PublicKey []byte
	Data      []byte
	Hash      []byte
}
