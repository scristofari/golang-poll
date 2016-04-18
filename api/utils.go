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
	// custom validation for bson
	validate.RegisterValidation("bson", bsonValidator)
}

// Is the ID is a valid bson hex ?
func bsonValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	return bson.IsObjectIdHex(field.String())
}

// Use validate V8 to validate the any struct.
// Validations for structs and individual fields based on tags
func valid(obj interface{}) error {
	return validate.Struct(obj)
}

// Use conform to sanitize the any struct
// Trims, sanitizes & scrubs data based on struct tags
func sanitize(obj interface{}) error {
	return conform.Strings(obj)
}

// Deserialize json to struct
func BindJson(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

// Serialize struct to json
// Write the correct status / content-type
func WriteJson(w http.ResponseWriter, obj interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(obj)
}

// Helper
// Sanitize / Validation
func IsValid(obj interface{}) error {
	if err := sanitize(obj); err != nil {
		return err
	}
	if err := valid(obj); err != nil {
		return err
	}
	return nil
}

// Bind form base to struct
// Gorilla schema
func BindForm(r *http.Request, obj interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	return decoder.Decode(obj, r.Form)
}
