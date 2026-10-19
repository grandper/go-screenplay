package send

import impl "github.com/grandper/go-screenplay/v2/extensions/http/action/internal/send"

// RequestWithMethod sends an HTTP request with the given method:
//
//	send.RequestWithMethod(http.MethodPatch).To(birdURL)
func RequestWithMethod(method string) *impl.RequestAction {
	return impl.RequestWithMethod(method)
}

// DeleteRequest sends a DELETE request.
func DeleteRequest() *impl.RequestAction {
	return impl.DeleteRequest()
}

// GetRequest sends a GET request.
func GetRequest() *impl.RequestAction {
	return impl.GetRequest()
}

// HeadRequest sends a HEAD request.
func HeadRequest() *impl.RequestAction {
	return impl.HeadRequest()
}

// OptionsRequest sends an OPTIONS request.
func OptionsRequest() *impl.RequestAction {
	return impl.OptionsRequest()
}

// PatchRequest sends a PATCH request.
func PatchRequest() *impl.RequestAction {
	return impl.PatchRequest()
}

// PostRequest sends a POST request.
func PostRequest() *impl.RequestAction {
	return impl.PostRequest()
}

// PutRequest sends a PUT request.
func PutRequest() *impl.RequestAction {
	return impl.PutRequest()
}
