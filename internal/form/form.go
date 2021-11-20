package form

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

type ErrorMap []map[string]string

// Validates the obj and returns user-readable errors
func BindJSONOrAbort(c *gin.Context, obj interface{}) ErrorMap {
	if err := c.BindJSON(obj); err != nil {
		errs := listOfErrors(obj, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errs})
		return errs
	}

	return nil
}

func listOfErrors(obj interface{}, e error) ErrorMap {
	ve := e.(validator.ValidationErrors)
	InvalidFields := make(ErrorMap, 0)

	for _, e := range ve {
		errors := map[string]string{}

		field, _ := reflect.TypeOf(obj).Elem().FieldByName(e.Field())
		jsonTag := string(field.Tag.Get("json"))
		errors[jsonTag] = validationErrorToText(e)
		InvalidFields = append(InvalidFields, errors)
	}

	return InvalidFields
}

func validationErrorToText(e validator.FieldError) string {
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
