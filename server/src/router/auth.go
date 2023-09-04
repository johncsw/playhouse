package router

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/model"
	"playhouse-server/repo"
	"playhouse-server/request"
	"playhouse-server/response"
)

func newAuthRouter() *chi.Mux {

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", registrationHandler())
	})
	return r
}

func registrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := &request.AuthRegistrationBody{}
		request.ToRequestBody(b, r)
		if request.IsNotValid(b) {
			panic(response.Error{
				Code:  http.StatusBadRequest,
				Cause: errors.New("email is not valid"),
			})
		}
		email := b.Email

		var usr *model.User
		checkUsr, checkErr := repo.UserRepo().GetUserByEmail(email, nil)
		if checkUsr != nil && checkErr == nil {
			usr = checkUsr
		} else {
			if checkErr != nil {
				panic(response.Error{
					Code:  http.StatusInternalServerError,
					Cause: checkErr,
				})
			}

			createUsr, createErr := repo.UserRepo().NewUser(b, nil)
			if createErr != nil {
				panic(response.Error{
					Code:  http.StatusInternalServerError,
					Cause: createErr,
				})
			}

			usr = createUsr
		}

		session, sessionErr := repo.SessionRepo().NewSession(usr.ID)
		if sessionErr != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: sessionErr,
			})
		}

		token := auth.CreateSessionToken(session)
		builder := response.Builder{Writer: w}
		builder.Header(map[string]string{
			"Authorization": token,
		}).Status(http.StatusCreated).BuildWithBytes(nil)
	}
}
