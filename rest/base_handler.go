package rest

import (
	"fmt"
	"net/http"
	"reflect"
)

// BaseResourceHandler is a base implementation of ResourceHandler with stubs for the
// CRUD operations. This allows ResourceHandler implementations to only implement
// what they need.
type BaseResourceHandler struct{}

// ResourceName is a stub. It must be implemented.
func (b BaseResourceHandler) ResourceName() string {
	return ""
}

// CreateURI is a stub. Implement if necessary. The default create URI is
// /api/v{version:[^/]+}/resourceName.
func (b BaseResourceHandler) CreateURI() string {
	return ""
}

// ReadURI is a stub. Implement if necessary. The default read URI is
// /api/v{version:[^/]+}/resourceName/{resource_id}.
func (b BaseResourceHandler) ReadURI() string {
	return ""
}

// ReadListURI is a stub. Implement if necessary. The default read list URI is
// /api/v{version:[^/]+}/resourceName.
func (b BaseResourceHandler) ReadListURI() string {
	return ""
}

// UpdateURI is a stub. Implement if necessary. The default update URI is
// /api/v{version:[^/]+}/resourceName/{resource_id}.
func (b BaseResourceHandler) UpdateURI() string {
	return ""
}

// DeleteURI is a stub. Implement if necessary. The default delete URI is
// /api/v{version:[^/]+}/resourceName/{resource_id}.
func (b BaseResourceHandler) DeleteURI() string {
	return ""
}

// EmptyResource is a stub. Implement if Rules are defined.
func (b BaseResourceHandler) EmptyResource() interface{} {
	return nil
}

// CreateResource is a stub. Implement if necessary.
func (b BaseResourceHandler) CreateResource(ctx RequestContext, data Payload,
	version string) (Resource, error) {
	return nil, NotImplemented("CreateResource is not implemented")
}

// ReadResourceList is a stub. Implement if necessary.
func (b BaseResourceHandler) ReadResourceList(ctx RequestContext, limit int,
	cursor string, version string) ([]Resource, string, error) {
	return nil, "", NotImplemented("ReadResourceList not implemented")
}

// ReadResource is a stub. Implement if necessary.
func (b BaseResourceHandler) ReadResource(ctx RequestContext, id string,
	version string) (Resource, error) {
	return nil, NotImplemented("ReadResource not implemented")
}

// UpdateResource is a stub. Implement if necessary.
func (b BaseResourceHandler) UpdateResource(ctx RequestContext, id string,
	data Payload, version string) (Resource, error) {
	return nil, NotImplemented("UpdateResource not implemented")
}

// DeleteResource is a stub. Implement if necessary.
func (b BaseResourceHandler) DeleteResource(ctx RequestContext, id string,
	version string) (Resource, error) {
	return nil, NotImplemented("DeleteResource not implemented")
}

// Authenticate is the default authentication logic. All requests are authorized.
// Implement custom authentication logic if necessary.
func (b BaseResourceHandler) Authenticate(r http.Request) error {
	return nil
}

// Rules returns the resource rules to apply to incoming requests and outgoing
// responses. No rules are applied by default. Implement if necessary.
func (b BaseResourceHandler) Rules() Rules {
	return Rules{}
}

// resourceHandlerProxy wraps a ResourceHandler and injects the resource type into its
// Rules. This allows the framework to provide additional logic around the proxied
// ResourceHandler. It also provides default logic such as REST URIs.
type resourceHandlerProxy struct {
	ResourceHandler
}

// ResourceName returns the wrapped ResourceHandler's resource name. If the proxied
// handler doesn't have ResourceName implemented, it panics.
func (r resourceHandlerProxy) ResourceName() string {
	name := r.ResourceHandler.ResourceName()
	if name == "" {
		panic("ResourceHandler must implement ResourceName()")
	}
	return name
}

// Rules returns the wrapped ResourceHandler's Rules after injecting them with the
// resource type.
func (r resourceHandlerProxy) Rules() Rules {
	rules := r.ResourceHandler.Rules()
	for _, rule := range rules {
		// Associate Rules with their respective type.
		rule.resourceType = reflect.TypeOf(r.EmptyResource())
	}

	return rules
}

// CreateURI returns the URI for creating a resource using the handler-specified
// URI while falling back to a sensible default if not provided.
func (r resourceHandlerProxy) CreateURI() string {
	url := r.ResourceHandler.CreateURI()
	if url == "" {
		url = fmt.Sprintf("/api/v{%s:[^/]+}/%s", versionKey, r.ResourceName())
	}
	return url
}

// ReadURI returns the URI for reading a specific resource using the handler-specified
// URI while falling back to a sensible default if not provided.
func (r resourceHandlerProxy) ReadURI() string {
	url := r.ResourceHandler.ReadURI()
	if url == "" {
		url = fmt.Sprintf("/api/v{%s:[^/]+}/%s/{%s}", versionKey, r.ResourceName(),
			resourceIDKey)
	}
	return url
}

// ReadListURI returns the URI for reading a list of resources using the handler-
// specified URI while falling back to a sensible default if not provided.
func (r resourceHandlerProxy) ReadListURI() string {
	url := r.ResourceHandler.ReadListURI()
	if url == "" {
		url = fmt.Sprintf("/api/v{%s:[^/]+}/%s", versionKey, r.ResourceName())
	}
	return url
}

// UpdateURI returns the URI for updating a specific resource using the handler-
// specified URI while falling back to a sensible default if not provided.
func (r resourceHandlerProxy) UpdateURI() string {
	url := r.ResourceHandler.UpdateURI()
	if url == "" {
		url = fmt.Sprintf("/api/v{%s:[^/]+}/%s/{%s}", versionKey, r.ResourceName(),
			resourceIDKey)
	}
	return url
}

// DeleteURI returns the URI for deleting a specific resource using the handler-
// specified URI while falling back to a sensible default if not provided.
func (r resourceHandlerProxy) DeleteURI() string {
	url := r.ResourceHandler.DeleteURI()
	if url == "" {
		url = fmt.Sprintf("/api/v{%s:[^/]+}/%s/{%s}", versionKey,
			r.ResourceHandler.ResourceName(), resourceIDKey)
	}
	return url
}
