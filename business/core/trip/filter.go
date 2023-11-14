package trip

import (
	"fmt"
	"time"

	"github.com/TSMC-Uber/server/business/sys/validate"
	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	ID             *uuid.UUID `validate:"omitempty"`
	DriverID       *uuid.UUID `validate:"omitempty"`
	PassengerLimit *int       `validate:"omitempty"`
	SourceID       *uuid.UUID `validate:"omitempty"`
	DestinationID  *uuid.UUID `validate:"omitempty"`
	StartStartDate *time.Time `validate:"omitempty"`
	EndStartDate   *time.Time `validate:"omitempty"`
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithTripID(id uuid.UUID) {
	qf.ID = &id
}

// WithDriverID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithDriverID(driverID uuid.UUID) {
	qf.DriverID = &driverID
}

// WithPassengerLimit sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithPassengerLimit(passengerLimit int) {
	qf.PassengerLimit = &passengerLimit
}

// WithSourceID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithSourceID(sourceID uuid.UUID) {
	qf.SourceID = &sourceID
}

// WithDestinationID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithDestinationID(destinationID uuid.UUID) {
	qf.DestinationID = &destinationID
}

// WithStartStartDate sets the DateCreated field of the QueryFilter value.
// query the trip which start time is after the startDate
func (qf *QueryFilter) WithStartStartDate(startDate time.Time) {
	d := startDate.UTC()
	qf.StartStartDate = &d
}

// WithEndCreatedDate sets the DateCreated field of the QueryFilter value.
// query the trip which start time is before the endDate
func (qf *QueryFilter) WithEndStartDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndStartDate = &d
}
