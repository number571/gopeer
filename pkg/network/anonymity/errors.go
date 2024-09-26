package anonymity

const (
	errPrefix = "pkg/network/anonymity = "
)

type SAnonymityError struct {
	str string
}

func (err *SAnonymityError) Error() string {
	return errPrefix + err.str
}

var (
	ErrSetHashIntoDB          = &SAnonymityError{"set hash into database"}
	ErrNilDB                  = &SAnonymityError{"database is nil"}
	ErrRetryLimit             = &SAnonymityError{"retry limit"}
	ErrEnqueueMessage         = &SAnonymityError{"enqueue message"}
	ErrUnknownType            = &SAnonymityError{"unknown type"}
	ErrLoadMessage            = &SAnonymityError{"load message"}
	ErrStoreHashWithBroadcast = &SAnonymityError{"store hash with broadcast"}
	ErrActionIsNotFound       = &SAnonymityError{"action is not found"}
	ErrActionIsClosed         = &SAnonymityError{"action is closed"}
	ErrActionTimeout          = &SAnonymityError{"action timeout"}
	ErrEnqueuePayload         = &SAnonymityError{"enqueue payload"}
	ErrFetchResponse          = &SAnonymityError{"fetch response"}
	ErrRunning                = &SAnonymityError{"node running"}
	ErrProcessRun             = &SAnonymityError{"process run"}
	ErrHashAlreadyExist       = &SAnonymityError{"hash already exist"}
)
