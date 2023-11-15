package utils

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/sivaosorg/govm/charge"
)

// SnapshotRequestBodyWith reads the request body and handles both form data and other types.
// It returns the raw body as a byte slice.
func SnapshotRequestBodyWith(ctx *gin.Context) ([]byte, error) {
	// Check if the request has a form content type
	if charge.IsPostForms(ctx.Request) {
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

func GetPostForms(ctx *gin.Context) map[string]string {
	ctx.MultipartForm()
	return charge.GetPostForms(ctx.Request)
}

func GetHeaders(ctx *gin.Context) map[string][]string {
	headers := make(map[string][]string)
	for key, values := range ctx.Request.Header {
		headers[key] = values
	}
	return headers
}
