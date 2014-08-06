package server

import (
	"encoding/json"

	"net/http"

	"go-rest/server/context"
)

// Resource represents a domain model.
type Resource interface{}

// Payload is the unmarshalled request body.
type Payload map[string]interface{}

// ResourceHandler specifies the endpoint handlers for working with a resource. This
// consists of the business logic for performing CRUD operations.
type ResourceHandler interface {
	ResourceName() string
	CreateResource(context.RequestContext, Payload, string) (Resource, error)
	ReadResourceList(context.RequestContext, int, string) ([]Resource, string, error)
	ReadResource(context.RequestContext, string, string) (Resource, error)
	UpdateResource(context.RequestContext, string, Payload, string) (Resource, error)
	DeleteResource(context.RequestContext, string, string) (Resource, error)
	IsAuthorized(http.Request) bool
}

type BaseResourceHandler struct{}

func (b BaseResourceHandler) ResourceName() string {
	panic("ResourceName not implemented")
}

func (b BaseResourceHandler) CreateResource(ctx context.RequestContext, data Payload,
	version string) (Resource, error) {
	panic("CreateResource not implemented")
}

func (b BaseResourceHandler) ReadResourceList(ctx context.RequestContext, limit int,
	version string) ([]Resource, string, error) {
	panic("ReadResourceList not implemented")
}

func (b BaseResourceHandler) ReadResource(ctx context.RequestContext, id string,
	version string) (Resource, error) {
	panic("ReadResource not implemented")
}

func (b BaseResourceHandler) UpdateResource(ctx context.RequestContext, id string,
	data Payload, version string) (Resource, error) {
	panic("UpdateResource not implemented")
}

func (b BaseResourceHandler) DeleteResource(ctx context.RequestContext, id string,
	version string) (Resource, error) {
	panic("DeleteResource not implemented")
}

func (b BaseResourceHandler) IsAuthorized(r http.Request) bool {
	return true
}

// requestHandler constructs http.HandlerFuncs responsible for handling HTTP requests.
type requestHandler struct {
	RestApi
}

// handleCreate returns a HandlerFunc which will deserialize the request payload, pass
// it to the provided create function, and then serialize and dispatch the response.
// The serialization mechanism used is specified by the "format" query parameter.
func (h requestHandler) handleCreate(createFunc func(context.RequestContext, Payload,
	string) (Resource, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(nil, r)

		decoder := json.NewDecoder(r.Body)
		var data map[string]interface{}
		if err := decoder.Decode(&data); err != nil {
			ctx = ctx.SetError(err)
			ctx = ctx.SetStatus(http.StatusInternalServerError)
		} else {
			resource, err := createFunc(ctx, data, ctx.Version())
			ctx = ctx.SetResult(resource)
			ctx = ctx.SetStatus(http.StatusCreated)
			if err != nil {
				ctx = ctx.SetError(err)
			}
		}

		h.sendResponse(w, ctx)
	}
}

// handleReadList returns a HandlerFunc which will pass the request context to the
// provided read function and then serialize and dispatch the response. The
// serialization mechanism used is specified by the "format" query parameter.
func (h requestHandler) handleReadList(readFunc func(context.RequestContext, int,
	string) ([]Resource, string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(nil, r)

		resources, cursor, err := readFunc(ctx, ctx.Limit(), ctx.Version())
		ctx = ctx.SetResult(resources)
		ctx = ctx.SetCursor(cursor)
		ctx = ctx.SetError(err)
		ctx = ctx.SetStatus(http.StatusOK)

		h.sendResponse(w, ctx)
	}
}

// handleRead returns a HandlerFunc which will pass the resource id to the provided
// read function and then serialize and dispatch the response. The serialization
// mechanism used is specified by the "format" query parameter.
func (h requestHandler) handleRead(readFunc func(context.RequestContext, string,
	string) (Resource, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(nil, r)

		resource, err := readFunc(ctx, ctx.ResourceId(), ctx.Version())
		ctx = ctx.SetResult(resource)
		ctx = ctx.SetError(err)
		ctx = ctx.SetStatus(http.StatusOK)

		h.sendResponse(w, ctx)
	}
}

// handleUpdate returns a HandlerFunc which will deserialize the request payload,
// pass it to the provided update function, and then serialize and dispatch the
// response. The serialization mechanism used is specified by the "format" query
// parameter.
func (h requestHandler) handleUpdate(updateFunc func(context.RequestContext,
	string, Payload, string) (Resource, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(nil, r)

		decoder := json.NewDecoder(r.Body)
		var data map[string]interface{}
		if err := decoder.Decode(&data); err != nil {
			ctx = ctx.SetError(err)
			ctx = ctx.SetStatus(http.StatusInternalServerError)
		} else {
			resource, err := updateFunc(ctx, ctx.ResourceId(), data, ctx.Version())
			ctx = ctx.SetResult(resource)
			ctx = ctx.SetError(err)
			ctx = ctx.SetStatus(http.StatusOK)
		}

		h.sendResponse(w, ctx)
	}
}

// handleDelete returns a HandlerFunc which will pass the resource id to the provided
// delete function and then serialize and dispatch the response. The serialization
// mechanism used is specified by the "format" query parameter.
func (h requestHandler) handleDelete(deleteFunc func(context.RequestContext, string,
	string) (Resource, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(nil, r)

		resource, err := deleteFunc(ctx, ctx.ResourceId(), ctx.Version())
		ctx = ctx.SetResult(resource)
		ctx = ctx.SetError(err)
		ctx = ctx.SetStatus(http.StatusOK)

		h.sendResponse(w, ctx)
	}
}

// sendResponse writes a success or error response to the provided http.ResponseWriter
// based on the contents of the context.RequestContext.
func (h requestHandler) sendResponse(w http.ResponseWriter, ctx context.RequestContext) {
	status := ctx.Status()
	requestError := ctx.Error()
	result := ctx.Result()

	serializer, err := h.responseSerializer(ctx.ResponseFormat())
	if err != nil {
		// Fall back to json serialization.
		serializer = jsonSerializer{}
		status = http.StatusNotImplemented
		requestError = err
	}

	if requestError != nil {
		if status < 400 {
			status = http.StatusInternalServerError
		}
		serializer.sendErrorResponse(w, requestError, status)
		return
	}

	nextURL, _ := ctx.NextURL()
	serializer.sendSuccessResponse(w, newSuccessResponse(result, nextURL), status)
}
