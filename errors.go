package hashdb

import (
	"errors"
)

var (
	ErrCreatingDatabase       = errors.New("Error while creating Database.")
	ErrCreatingTable          = errors.New("Error while creating Table.")
	ErrCreatingIndex          = errors.New("Error while creating Index.")
	ErrInitExactHandlerFailed = errors.New("Exact-handler initialization failed.")
	ErrInitLikeHandlerFailed  = errors.New("Like-handler initialization failed.")
	ErrInitPutHandlerFailed   = errors.New("Put-handler initialization failed.")
	ErrHashTooShort           = errors.New("The length of the hash is too short.")
	ErrDatastoreNotFound      = errors.New("The requested datastore was not found.")
	ErrDirectoryNotFound      = errors.New("The directory doesn't exist.")
)
