package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"os"

	"github.com/faststackco/faststack/auth/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	stripe "github.com/stripe/stripe-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	testDb *gorm.DB
)

func TestMain(m *testing.M) {
	stripe.Key = "sk_test_wXzOwaHdvv9UivuOF8b0pvTG"

	var err error

	testDb, err = model.Db("sqlite3", "/tmp/test.db")
	if err != nil {
		panic("failed to connect to database")
	}

	retCode := m.Run()

	testDb.Close()
	os.Remove("/tmp/test.db")

	os.Exit(retCode)
}

func TestCreateUser(t *testing.T) {
	json := `{"name":"username", "email":"user@email.com", "password": "qqqqwwww"}`

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/user", strings.NewReader(json))
	if assert.NoError(t, err) {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := NewHandler(testDb)

		if assert.NoError(t, h.SignUp(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "created", rec.Body.String())
		}
	}

}

func TestLogin(t *testing.T) {
	pass, err := bcrypt.GenerateFromPassword([]byte("password3021"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	testDb.Create(&model.User{
		Name:           "username3021",
		PasswordBcrypt: pass,
	})
	assert.NoError(t, testDb.Error)

	json := `{"name":"username3021", "password": "password3021"}`

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/user/login", strings.NewReader(json))
	if assert.NoError(t, err) {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := &Handler{testDb}

		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// TODO validate token
		}
	}

}