package locationwsdb

import "github.com/TSMC-Uber/server/foundation/logger"

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger) *Store {
	return &Store{
		log: log,
	}
}
