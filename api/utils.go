package api

import (
	"encoding/json"
	"net/http"
	"reflect"

	"gopkg.in/go-playground/validator.v8"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/schema"
	"github.com/leebenson/conform"
)

var (
	validate *validator.Validate
)

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
	// custom validation
	validate.RegisterValidation("bson", bsonValidator)
}

// Is the ID is a valid bson hex ?
func bsonValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	return bson.IsObjectIdHex(field.String())
}

func valid(obj interface{}) error {
	return validate.Struct(obj)
}

func sanitize(obj interface{}) error {
	return conform.Strings(obj)
}

func BindJson(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, obj interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(obj)
}

func IsValid(obj interface{}) error {
	if err := sanitize(obj); err != nil {
		return err
	}
	if err := valid(obj); err != nil {
		return err
	}
	return nil
}

func BindForm(r *http.Request, obj interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	return decoder.Decode(obj, r.Form)
}
