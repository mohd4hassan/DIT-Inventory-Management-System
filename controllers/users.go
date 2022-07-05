package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"IMS/context"
	"IMS/models"
	"IMS/rand"
	"IMS/views"

	"github.com/gorilla/mux"
)

const (
	ShowUsers = "show_user"
	EditUsers = "edit_user"
)

type Users struct {
	NewView       *views.View
	IndexView     *views.View
	EditView      *views.View
	LoginView     *views.View
	ForgotPwdView *views.View
	ResetPwdView  *views.View
	us            models.UserService
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:       views.NewView("masterLayout", "users/new"),
		IndexView:     views.NewView("masterLayout", "users/users"),
		EditView:      views.NewView("masterLayout", "users/edit"),
		LoginView:     views.NewView("loginLayout", "users/login"),
		ForgotPwdView: views.NewView("pwdForgot", "users/pwd_forgot"),
		ResetPwdView:  views.NewView("pwdReset", "users/pwd_reset"),
		us:            us,
	}
}

type SignupForm struct {
	//using struct tags because the name on form and struct arent matching
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Role     string `schema:"role"`
	Password string `schema:"password"`
}

//New is used to render the form where a user can create
//a new user account.
//
//GET /register
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	user := context.User(r.Context())

	parseURLParams(r, &form)
	u.NewView.Render(w, r, user.Role)
}

//Create is used to process the signup form when a user tries to
//create a new user account.
//
//POST /register
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Username: form.Username,
		Email:    form.Email,
		Role:     form.Role,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	/* err := u.signIn(w, &user)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} */

	views.RedirectAlert(w, r, "/register", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "User created successfully",
	})

	// http.Redirect(w, r, "/register", http.StatusFound)
}

//GET /users
func (usr *Users) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	users, err := usr.us.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		views.RedirectAlert(w, r, "/", http.StatusFound, views.Alert{
			Level:   views.AlertLvlError,
			Message: "Only admins can access this page",
		})
		// http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: users,
	}

	usr.IndexView.Render(w, r, vd)
}

//Get /users/:id/edit
func (usr *Users) Edit(w http.ResponseWriter, r *http.Request) {
	users, err := usr.userByID(w, r)
	if err != nil {
		log.Println(err)
		views.RedirectAlert(w, r, "/", http.StatusFound, views.Alert{
			Level:   views.AlertLvlError,
			Message: "Only admins can access this page",
		})
		// http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var vd views.Data
	vd.Yield = users
	usr.EditView.Render(w, r, vd)
}

//POST /users/:id/update
func (usr *Users) Update(w http.ResponseWriter, r *http.Request) {
	users, err := usr.userByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Yield = users
	var form SignupForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		usr.EditView.Render(w, r, vd)
		return
	}

	users.Username = form.Username
	users.Email = form.Email
	users.Password = form.Password
	users.Role = form.Role

	err = usr.us.Update(users)
	if err != nil {
		vd.SetAlert(err)
		usr.EditView.Render(w, r, vd)
		return
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "User updated successfully",
	}

	usr.EditView.Render(w, r, vd)

	http.Redirect(w, r, "/users", http.StatusFound)
}

//POST /users/:id/delete
func (usr *Users) Delete(w http.ResponseWriter, r *http.Request) {
	users, err := usr.userByID(w, r)
	if err != nil {
		log.Println(err)
		views.RedirectAlert(w, r, "/", http.StatusFound, views.Alert{
			Level:   views.AlertLvlError,
			Message: "Only admins can delete users",
		})
		// http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())

	var vd views.Data
	err = usr.us.Delete(user.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = users
		usr.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/users", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete a user's session cookie
// and invalidate their current remember token, which will
// sign the current user out.
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// First expire the user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	// Then we update the user with a new remember token
	user := context.User(r.Context())

	// We are ignoring errors for now because they are
	// unlikely, and even if they do occur we can't recover
	// now that the user doesn't have a valid cookie
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)

	// Send the user to the login page
	http.Redirect(w, r, "/login", http.StatusFound)
}

// ResetPwForm is used to process the forgot password form
// and the reset password form.
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwdView.Render(w, r, vd)
		return
	}
	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwdView.Render(w, r, vd)
		return
	}
	// TODO: Email the user their password reset token. In the
	// meantime we need to "use" the token variable
	_ = token
	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for password reset has been sent to your email.",
	})
}

// ResetPw displays the reset password form and has a method
// so that we can prefill the form data with a token provided
// via the URL query params.
//
// GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwdView.Render(w, r, vd)
}

// CompleteReset processed the reset password form
//
// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	var form ResetPwForm

	vd.Yield = &form

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwdView.Render(w, r, vd)
		return
	}

	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwdView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)

	views.RedirectAlert(w, r, "/", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
}

func (usr *Users) userByID(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid User ID", http.StatusNotFound)
		return nil, err
	}
	user, err := usr.us.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return user, nil
}
