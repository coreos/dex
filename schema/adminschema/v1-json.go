package adminschema
//
// This file is automatically generated by schema/generator
//
// **** DO NOT EDIT ****
//
const DiscoveryJSON = `{
  "kind": "discovery#restDescription",
  "discoveryVersion": "v1",
  "id": "dex:v1",
  "name": "adminschema",
  "version": "v1",
  "title": "Dex Admin API",
  "description": "The Dex Admin API.",
  "documentationLink": "http://github.com/coreos/dex",
  "protocol": "rest",
  "icons": {
    "x16": "",
    "x32": ""
  },
  "labels": [],
  "baseUrl": "$ENDPOINT/api/v1/",
  "basePath": "/api/v1/",
  "rootUrl": "$ENDPOINT/",
  "servicePath": "api/v1/",
  "batchPath": "batch",
  "parameters": {},
  "auth": {},
  "schemas": {
      "Admin": {
          "id": "Admin",
          "type": "object",
          "properties": {
              "id": {
                  "type": "string"
              },
              "email": {
                  "type": "string"
              },
              "password": {
                  "type": "string"
              }
          }
      },
      "State": {
          "id": "State",
          "type": "object",
          "properties": {
              "AdminUserCreated": {
                  "type": "boolean"
              }
          }
      },
      "ClientCreateRequest": {
          "id": "ClientCreateRequest",
          "type": "object",
	      "description": "'client' field is a client registration request as defined by the OpenID Connect dynamic registration spec, and holds fields such as redirect URLs, prefered algorithms, etc. For brevity field names and types of that object have been omitted.",
          "properties": {
              "isAdmin": {
                  "type": "boolean"
              },
              "client": {
                  "type": "object"
              }
          }
      },
	  "ClientRegistrationResponse": {
		  "id": "ClientRegistrationResponse",
		  "type": "object",
		  "description": "This object is a client registration respones as defined by the OpenID Connect dynamic registration spec. For brevity field names and types have been omitted.",
		  "properties": {}
	  }
  },
  "resources": {
      "Admin": {
          "methods": {
              "Get": {
                  "id": "dex.admin.Admin.Get",
                  "description": "Retrieve information about an admin user.",
                  "httpMethod": "GET",
                  "path": "admin/{id}",
                  "parameters": {
                      "id": {
                          "type": "string",
                          "required": true,
                          "location": "path"
                      }
                  },
                  "parameterOrder": [
                      "id"
                  ],
                  "response": {
                      "$ref": "Admin"
                  }
                  
              },
              "Create": {
                  "id": "dex.admin.Admin.Create",
                  "description": "Create a new admin user.",
                  "httpMethod": "POST",
                  "path": "admin",
                  "request": {
                      "$ref": "Admin"
                  },
                  "response": {
                      "$ref": "Admin"
                  }
              }
          }
      },
      "State": {
          "methods": {
              "Get": {
                  "id": "dex.admin.State.Get",
                  "description": "Get the state of the Dex DB",
                  "httpMethod": "GET",
                  "path": "state",
                  "response": {
                      "$ref": "State"
                  }
              }
          }
      },
      "Client": {
          "methods": {
              "Create": {
                  "id": "dex.admin.Client.Create",
                  "description": "Register an OpenID Connect client.",
                  "httpMethod": "POST",
                  "path": "client",
                  "request": {
                      "$ref": "ClientCreateRequest"
                  },
                  "response": {
                      "$ref": "ClientRegistrationResponse"
                  }
              }
          }
      }
  }
}
`