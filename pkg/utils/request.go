package utils

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

// SnapshotRequestBodyWith reads the request body and handles both form data and other types.
// It returns the raw body as a byte slice.
func SnapshotRequestBodyWith(ctx *gin.Context) ([]byte, error) {
	// Check if the request has a form content type
	if IsPostForm(ctx) {
		// err := ctx.Request.ParseForm()
		// if err != nil {
		// 	return nil, err
		// }
		ctx.MultipartForm()
		// Form data is already parsed, so you can access it using ctx.Request.PostForm
		// Example: formData := ctx.Request.PostForm
		// Capture the raw form data as a byte slice
		rawFormData := []byte(ctx.Request.PostForm.Encode())
		// Replace the request's body with a new NopCloser
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawFormData))
		return rawFormData, nil
	}
	// If the content type is not a form, read the request body as is
	raw, err := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
	return raw, err
}

func IsPostForm(ctx *gin.Context) bool {
	value := ctx.Request.Header.Get("Content-Type")
	return value == "application/x-www-form-urlencoded" || strings.Contains(value, "multipart/form-data")
}

func GetPostFromValues(ctx *gin.Context) map[string]string {
	ctx.MultipartForm()
	formValues := make(map[string]string)
	for key, values := range ctx.Request.PostForm {
		if len(values) > 0 {
			formValues[key] = values[0]
		}
	}
	return formValues
}
