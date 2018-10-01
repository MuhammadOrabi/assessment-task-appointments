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


/* Middleware handler to handle all requests for authentication */
// func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
//     return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//         authorizationHeader := req.Header.Get("authorization")
//         if authorizationHeader != "" {
//             bearerToken := strings.Split(authorizationHeader, " ")
//             if len(bearerToken) == 2 {
//                 token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
//                     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//                         return nil, fmt.Errorf("There was an error")
//                     }
//                     return []byte("secret"), nil
//                 })
//                 if error != nil {
//                     json.NewEncoder(w).Encode(Exception{Message: error.Error()})
//                     return
//                 }
//                 if token.Valid {
//                     log.Println("TOKEN WAS VALID")
//                     context.Set(req, "decoded", token.Claims)
//                     next(w, req)
//                 } else {
//                     json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
//                 }
//             }
//         } else {
//             json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
//         }
//     })
// }

// Get Authentication token GET /
// func (c *Controller) GetToken(w http.ResponseWriter, req *http.Request) {
//     var user User
//     _ = json.NewDecoder(req.Body).Decode(&user)
//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//         "username": user.Username,
//         "password": user.Password,
//     })

//     log.Println("Username: " + user.Username);
//     log.Println("Password: " + user.Password);

//     tokenString, error := token.SignedString([]byte("secret"))
//     if error != nil {
//         fmt.Println(error)
//     }
//     json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
// }

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
    w.WriteHeader(http.StatusCreated)
    return
}

// SearchAppointment GET /
func (c *Controller) SearchAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println(vars)

    query := vars["query"] // param query
    log.Println("Search Query - " + query);

    appointments := c.Repository.GetAppointmentsByString(query)
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