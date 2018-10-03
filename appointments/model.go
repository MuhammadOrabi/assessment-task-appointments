package appointments

// Appointment represents an appointment
type Appointment struct {
	ID     	   uint64  `bson:"_id"`
	DoctorId   string  `bson:"doctor_id"`
	PatientId  int  `bson:"patient_id"`
	Day   	   int  `bson:"day"`
	Hour       int  `bson:"hour"`
}

// Appointments is an array of Appointment objects
type Appointments []Appointment