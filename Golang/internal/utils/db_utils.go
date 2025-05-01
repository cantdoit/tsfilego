package utils

import (
	"fmt"
)

// NodeID represents a node identifier.
type NodeID uint16

// TsID represents a Time Series Identifier.
type TsID struct {
	DbNID          NodeID
	DeviceNID      NodeID
	MeasurementNID NodeID
}

// NewTsID creates a new TsID.
func NewTsID(dbNID NodeID, deviceNID NodeID, measurementNID NodeID) TsID {
	return TsID{
		DbNID:          dbNID,
		DeviceNID:      deviceNID,
		MeasurementNID: measurementNID,
	}
}

// Reset resets the TsID to its default state.
func (ts TsID) Reset() {
	ts.DbNID = 0
	ts.DeviceNID = 0
	ts.MeasurementNID = 0
}

// IsValid checks whether the TsID is valid.
func (ts TsID) IsValid() bool {
	return ts.DbNID != 0 || ts.DeviceNID != 0 || ts.MeasurementNID != 0
}

// ToInt64 converts the TsID to a 64-bit integer representation.
func (ts TsID) ToInt64() int64 {
	res := int64(ts.DbNID)
	res = (res << 16) | int64(ts.DeviceNID)
	res = (res << 16) | int64(ts.MeasurementNID)
	return res
}

// String returns a string representation of the TsID.
func (ts TsID) String() string {
	return fmt.Sprintf("<%d,%d,%d>", ts.DbNID, ts.DeviceNID, ts.MeasurementNID)
}

// Equals compares two TsID objects for equality.
func (ts TsID) Equals(other TsID) bool {
	return ts.DbNID == other.DbNID &&
		ts.DeviceNID == other.DeviceNID &&
		ts.MeasurementNID == other.MeasurementNID
}

// LessThan checks if the current TsID is less than another TsID.
func (ts TsID) LessThan(other TsID) bool {
	return ts.ToInt64() < other.ToInt64()
}

// GreaterThan checks if the current TsID is greater than another TsID.
func (ts TsID) GreaterThan(other TsID) bool {
	return ts.ToInt64() > other.ToInt64()
}

// DeviceID represents a Device Identifier.
type DeviceID struct {
	DbNID     NodeID
	DeviceNID NodeID
}

// NewDeviceID creates a new DeviceID.
func NewDeviceID(dbNID NodeID, deviceNID NodeID) DeviceID {
	return DeviceID{
		DbNID:     dbNID,
		DeviceNID: deviceNID,
	}
}

// FromTsID initializes a DeviceID from a TsID.
func (d DeviceID) FromTsID(ts TsID) {
	d.DbNID = ts.DbNID
	d.DeviceNID = ts.DeviceNID
}

// Equals compares two DeviceID objects for equality.
func (d DeviceID) Equals(other DeviceID) bool {
	return d.DbNID == other.DbNID &&
		d.DeviceNID == other.DeviceNID
}

// NotEquals checks if the current DeviceID is not equal to another DeviceID.
func (d DeviceID) NotEquals(other DeviceID) bool {
	return !d.Equals(other)
}

// LessThan checks if the current DeviceID is less than another DeviceID.
func (d DeviceID) LessThan(other DeviceID) bool {
	thisI32 := (int32(d.DbNID) << 16) | int32(d.DeviceNID)
	thatI32 := (int32(other.DbNID) << 16) | int32(other.DeviceNID)
	return thisI32 < thatI32
}

// String returns a string representation of the DeviceID.
func (d DeviceID) String() string {
	return fmt.Sprintf("<%d,%d>", d.DbNID, d.DeviceNID)
}
