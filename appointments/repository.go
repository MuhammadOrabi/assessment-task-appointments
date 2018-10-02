package appointments

import (
	"fmt"
	"log"
	"github.com/night-codes/mgo-ai"
	"gopkg.in/mgo.v2"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"os"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
var SERVER = os.Getenv("DATABASE_URL")

// DBNAME the name of the DB instance
const DBNAME = "appointments"

// COLLECTION is the name of the collection in DB
const COLLECTION = "appointments"


// GetAppointmet returns the list of Appointments
func (r Repository) GetAppointments() Appointments {
	session, err := mgo.Dial(SERVER)

	if err != nil {
	 	fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Appointments{}

	if err := c.Find(nil).All(&results); err != nil {
	  	fmt.Println("Failed to write results:", err)
	}

	return results
}

// GetAppointmentById returns a unique Appointment
func (r Repository) GetAppointmentById(id int) Appointment {
	session, err := mgo.Dial(SERVER)

	if err != nil {
	 	fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result Appointment

	fmt.Println("ID in GetAppointmentById", id);

	if err := c.FindId(id).One(&result); err != nil {
	  	fmt.Println("Failed to write result:", err)
	}

	return result
}

// GetAppointmentsByString takes a search string as input and returns Appointments
func (r Repository) GetAppointmentsByString(query string) Appointments {
	session, err := mgo.Dial(SERVER)

	if err != nil {
	 	fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	result := Appointments{}

	// Logic to create filter
	qs := strings.Split(query, " ")
	and := make([]bson.M, len(qs))
	for i, q := range qs {
    	and[i] = bson.M{"DoctorId": bson.M{
        	"$regex": bson.RegEx{Pattern: ".*" + q + ".*", Options: "i"},
    	}}
	}
	filter := bson.M{"$and": and}

	if err := c.Find(&filter).Limit(5).All(&result); err != nil {
	  	fmt.Println("Failed to write result:", err)
	}

	return result
}

// AddAppointment adds a Appointment in the DB
func (r Repository) AddAppointment(appointment Appointment) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	ai.Connect(session.DB(DBNAME).C(COLLECTION))

	appointment.ID = ai.Next(COLLECTION)
	session.DB(DBNAME).C(COLLECTION).Insert(appointment)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Appointment ID- ", appointment.ID)

	return true
}

// UpdateAppointment updates a Appointment in the DB
func (r Repository) UpdateAppointment(appointment Appointment) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	err = session.DB(DBNAME).C(COLLECTION).UpdateId(appointment.ID, appointment)
	
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Updated Appointment ID - ", appointment.ID)

	return true
}

// DeleteAppointment deletes an Appointment
func (r Repository) DeleteAppointment(id int) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	// Remove appointment
	if err = session.DB(DBNAME).C(COLLECTION).RemoveId(id); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}

	fmt.Println("Deleted Appointment ID - ", id)
	// Write status
	return "OK"
}
