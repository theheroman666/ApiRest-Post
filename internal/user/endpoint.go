package user

import (
	"encoding/json"
	"goweb/pkg/meta"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	EndPoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	// ErrorRespose struct {
	// 	Error string `json:"error"`
	// }
	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndPoints(serv Service) EndPoints {
	return EndPoints{
		Create: makeCreateEndPoints(serv),
		Get:    makeGetEndPoints(serv),
		GetAll: makeGetAllEndPoints(serv),
		Update: makeUpdateEndPoints(serv),
		Delete: makeDeleteEndPoints(serv),
	}
}

func makeCreateEndPoints(serv Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req CreateReq
		govalidator.IsEmail(req.Email)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}
		if strings.TrimSpace(req.FirstName) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "The First Name is required"})
			return
		}
		if strings.TrimSpace(req.LastName) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "The Last Name is required"})
			return

		}
		if strings.TrimSpace(req.Email) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "The Email is required"})

		}
		if strings.TrimSpace(req.Phone) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "T&he Phone is required"})
			return

		}

		// fmt.Println(req)
		user, err := serv.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		// json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		// json.NewEncoder(w).Encode(user)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
	}
}

func makeGetAllEndPoints(serv Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		paramfilter := r.URL.Query()

		filters := Filters{
			FirstnameF: paramfilter.Get("first_name"),
			LastnameF:  paramfilter.Get("last_name"),
		}
		limit, _ := strconv.Atoi(paramfilter.Get("limit"))
		page, _ := strconv.Atoi(paramfilter.Get("page"))

		count, err := serv.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		meta, err := meta.New(page, limit, count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		users, err := serv.GetAll(filters, meta.Offser(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: users, Meta: meta})
	}
}

func makeGetEndPoints(serv Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := mux.Vars(r)
		id := path["id"]
		user, err := serv.Get(id)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "user doesn't exist"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
	}
}

func makeUpdateEndPoints(serv Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req UpdateReq
		// firstname := strings.TrimSpace(*req.FirstName)
		// lastname := strings.TrimSpace(*req.LastName)
		// email := strings.TrimSpace(*req.Email)
		// phone := strings.TrimSpace(*req.Phone)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		if req.FirstName != nil && strings.TrimSpace(*req.FirstName) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "firstname is required"})
			return
		}
		if req.LastName != nil && strings.TrimSpace(*req.LastName) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "firstname is required"})
			return
		}
		if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "firstname is required"})
			return
		}
		if req.Phone != nil && strings.TrimSpace(*req.Phone) == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "firstname is required"})
			return
		}
		path := mux.Vars(r)
		id := path["id"]

		if err := serv.Update(id, req.FirstName, req.LastName, req.Email, req.Phone); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "User not exist"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "Success"})
	}
}

func makeDeleteEndPoints(serv Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := mux.Vars(r)
		id := path["id"]
		if err := serv.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "the user doesn't exist"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}
