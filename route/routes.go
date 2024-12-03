package route

import (
	"net/http"

	"github.com/IsraelTeo/api-paw-go/auth"
	"github.com/IsraelTeo/api-paw-go/handler"
	"github.com/IsraelTeo/api-paw-go/middelware"
	"github.com/gorilla/mux"
)

const (
	registerPath = "/sign-up"
	loginPath    = "/login"

	userBasicPath = "/user"
	userIDPath    = "/user/{id}"
	usersPath     = "/users"

	employeTypeBasicPath = "/type"
	employeTypeIDPath    = "/type/{id}"
	employeTypesPath     = "/types"

	employeeBasicPath = "/employee"
	employeeIDPath    = "/employee/{id}"
	employeesPath     = "/employees"

	customerBasicPath = "/customer"
	customerIDPath    = "/customer/{id}"
	customersPath     = "/customers"

	petBasicPath = "/pet"
	petIDPath    = "/pet/{id}"
	petsPath     = "/pets"
)

func Init() *mux.Router {
	routes := mux.NewRouter()

	routes.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	apiAuth := routes.PathPrefix("/auth").Subrouter()

	apiAuth.HandleFunc(registerPath, middelware.Log(handler.RegisterUser)).Methods("POST")
	apiAuth.HandleFunc(loginPath, middelware.Log(auth.Login)).Methods("POST")

	api := routes.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc(userIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.GetUserById))).Methods("GET")
	api.HandleFunc(usersPath, middelware.ValidateJWTAdmin(middelware.Log(handler.GetAllUsers))).Methods("GET")
	api.HandleFunc(userIDPath, middelware.Log(handler.UpdateUser)).Methods("PUT")
	api.HandleFunc(userIDPath, middelware.Log(handler.DeleteUser)).Methods("DELETE")

	api.HandleFunc(employeTypeBasicPath, middelware.ValidateJWTAdmin(middelware.Log(handler.SaveEmployeeType))).Methods("POST")
	api.HandleFunc(employeTypeIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.GetEmployeeTypeById))).Methods("GET")
	api.HandleFunc(employeTypesPath, middelware.ValidateJWTAdmin(middelware.Log(handler.GetAllEmployeeTypes))).Methods("GET")
	api.HandleFunc(employeTypeIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.UpdateEmployeeType))).Methods("PUT")
	api.HandleFunc(employeTypeIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.DeleteEmployeeType))).Methods("DELETE")

	api.HandleFunc(employeeBasicPath, middelware.ValidateJWTAdmin(middelware.Log(handler.SaveEmployee))).Methods("POST")
	api.HandleFunc(employeeIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.GetEmployeeById))).Methods("GET")
	api.HandleFunc(employeesPath, middelware.ValidateJWTAdmin(middelware.Log(handler.GetAllEmployees))).Methods("GET")
	api.HandleFunc(employeeIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.UpdateEmployee))).Methods("PUT")
	api.HandleFunc(employeeIDPath, middelware.ValidateJWTAdmin(middelware.Log(handler.DeleteEmployee))).Methods("DELETE")

	api.HandleFunc(customerBasicPath, middelware.ValidateJWT(middelware.Log(handler.SaveCustomer))).Methods("POST")
	api.HandleFunc(customerIDPath, middelware.ValidateJWT(middelware.Log(handler.GetCustomerById))).Methods("GET")
	api.HandleFunc(customersPath, middelware.ValidateJWT(middelware.Log(handler.GetAllCustomers))).Methods("GET")
	api.HandleFunc(customerIDPath, middelware.ValidateJWT(middelware.Log(handler.UpdateCustomer))).Methods("PUT")
	api.HandleFunc(customerIDPath, middelware.ValidateJWT(middelware.Log(handler.DeleteCustomer))).Methods("DELETE")

	api.HandleFunc(petBasicPath, middelware.ValidateJWT(middelware.Log(handler.SavePet))).Methods("POST")
	api.HandleFunc(petIDPath, middelware.ValidateJWT(middelware.Log(handler.GetPetById))).Methods("GET")
	api.HandleFunc(petsPath, middelware.ValidateJWT(middelware.Log(handler.GetAllPets))).Methods("GET")
	api.HandleFunc(petIDPath, middelware.ValidateJWT(middelware.Log(handler.UpdatePet))).Methods("PUT")
	api.HandleFunc(petIDPath, middelware.ValidateJWT(middelware.Log(handler.DeletePet))).Methods("DELETE")

	return routes
}
