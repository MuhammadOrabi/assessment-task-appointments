package appointments

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
    // "fmt"
    "net/http"
    "strings"
    "strconv"

    "github.com/gorilla/mux"
    // "github.com/gorilla/context"
    // "github.com/dgrijalva/jwt-go"

)

//Controller ...
type Controller struct {
    Repository Repository
}


// Index GET /
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
    appointments := c.Repository.GetAppointments() // list of all appointments
    // log.Println(appointments)
    data, _ := json.Marshal(appointments)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}

// AddAppointment POST /
func (c *Controller) AddAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment Appointment
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
    
    log.Println(body)

    if err != nil {
        log.Fatalln("Error AddAppointment", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := r.Body.Close(); err != nil {
        log.Fatalln("Error AddAppointment", err)
    }

    if err := json.Unmarshal(body, &appointment); err != nil { // unmarshall body contents as a type Candidate
        w.WriteHeader(422) // unprocessable entity
        log.Println(err)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatalln("Error AddAppointment unmarshalling data", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }

    log.Println(appointment)
    success := c.Repository.AddAppointment(appointment) // adds the appointment to the DB
    if !success {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusCreated)
    return
}

// SearchAppointment GET /
func (c *Controller) SearchAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println(vars)

    query := vars["query"] // param query
    log.Println("Search Query - " + query);

    appointments := c.Repository.GetAppointmentsByString(query, "doctor_id")
    data, _ := json.Marshal(appointments)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}


// SearchAppointmentByPatient GET /
func (c *Controller) SearchAppointmentByPatient(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println(vars)

    query := vars["query"] // param query
    log.Println("Search Query - " + query);

    appointments := c.Repository.GetAppointmentsByString(query, "patient_id")
    data, _ := json.Marshal(appointments)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}

// UpdateAppointment PUT /
func (c *Controller) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment Appointment
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
    if err != nil {
        log.Fatalln("Error UpdateAppointment", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if err := r.Body.Close(); err != nil {
        log.Fatalln("Error UpdateAppointment", err)
    }

    if err := json.Unmarshal(body, &appointment); err != nil { // unmarshall body contents as a type Candidate
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatalln("Error UpdateAppointment unmarshalling data", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }

    log.Println(appointment.ID)
    success := c.Repository.UpdateAppointment(appointment) // updates the product in the DB
    
    if !success {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    return
}

// GetAppointment GET - Gets a single appointment by ID /
func (c *Controller) GetAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println(vars)

    id := vars["id"] // param id
    log.Println(id);

    appointmentid, err := strconv.Atoi(id);

    if err != nil {
        log.Fatalln("Error GetAppointment", err)
    }

    appointment := c.Repository.GetAppointmentById(appointmentid)
    data, _ := json.Marshal(appointment)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}

// DeleteAppointment DELETE /
func (c *Controller) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println(vars)
    id := vars["id"] // param id
    log.Println(id);

    appointmentid, err := strconv.Atoi(id);

    if err != nil {
        log.Fatalln("Error GetAppointment", err)
    }

    if err := c.Repository.DeleteAppointment(appointmentid); err != "" { // delete a appointment by id
        log.Println(err);
        if strings.Contains(err, "404") {
            w.WriteHeader(http.StatusNotFound)
        } else if strings.Contains(err, "500") {
            w.WriteHeader(http.StatusInternalServerError)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    return
}