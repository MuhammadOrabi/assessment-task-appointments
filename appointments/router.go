package appointments

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes = Routes {
    // Route {
    //     "Authentication",
    //     "POST",
    //     "/get-token",
    //     controller.GetToken,
    // },
    Route {
        "Index",
        "GET",
        "/",
        controller.Index,
    },
    Route {
        "AddAppointment",
        "POST",
        "/add",
        controller.AddAppointment,
    },
    Route {
        "UpdateAppointment",
        "PUT",
        "/update",
        controller.UpdateAppointment,
    },
    // Get Appointment by {id}
    Route {
        "GetAppointment",
        "GET",
        "/{id}",
        controller.GetAppointment,
    },
    // Delete Appointment by {id}
    Route {
        "DeleteAppointment",
        "DELETE",
        "/delete/{id}",
        controller.DeleteAppointment,
    },
    // Search Appointment with string
    Route {
        "SearchAppointment",
        "GET",
        "/Search/{query}",
        controller.SearchAppointment,
    },
}

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        log.Println(route.Name)
        handler = route.HandlerFunc
        
        router.
         Methods(route.Method).
         Path("/api/appointments" + route.Pattern).
         Name(route.Name).
         Handler(handler)
    }
    return router
}