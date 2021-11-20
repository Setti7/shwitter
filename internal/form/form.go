package form

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

func BindJSON(c *gin.Context, obj interface{}) []map[string]string {
	if err := c.BindJSON(obj); err != nil {
		errs := listOfErrors(obj, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errs})
		return errs
	}

	return nil
}

// TODO: return list of errors for each field
func listOfErrors(obj interface{}, e error) []map[string]string {
	ve := e.(validator.ValidationErrors)
	InvalidFields := make([]map[string]string, 0)

	for _, e := range ve {
		errors := map[string]string{}

		field, _ := reflect.TypeOf(obj).Elem().FieldByName(e.Field())
		jsonTag := string(field.Tag.Get("json"))
		errors[jsonTag] = ValidationErrorToText(e)
		InvalidFields = append(InvalidFields, errors)
	}

	return InvalidFields
}

func ValidationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required."
	case "max":
		return fmt.Sprintf("This field cannot be longer than %s characters.", e.Param())
	case "min":
		return fmt.Sprintf("This field must be longer than %s characters.", e.Param())
	case "email":
		return "This email is invalid."
	case "len":
		return fmt.Sprintf("This field must be %s characters long.", e.Param())
	}
	return "This field is invalid."
}
