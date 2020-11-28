package sdkc

// TODO: https://github.com/aws-samples/aws-greengrass-lambda-functions

/*
#include "lib/greengrasssdk.h"
*/
import "C"

import (
	"encoding/json"
	"time"
	"unsafe"
)

// Secret represents a single secret obtained from the
// green grass secrets manager.
type Secret struct {
	// ARN is the ARN of the secret.
	ARN string
	// Name is the friendly name of the secret.
	Name string
	// VersionID is the unique identifier of this version of the secret.
	VersionID string `json:"VersionId"`
	// SecretBinary is the decrypted part of the protected secret information that was originally provided as
	// binary data in the form of a byte array. This parameter is not used if the secret is created by the Secrets Manager console.
	//
	// If you store custom information in this field of the secret, then you must code your Lambda rotation function to parse and
	// interpret whatever you store in the _SecretString_ or _SecretBinary_ fields.
	SecretBinary []byte
	// SecretString is the decrypted part of the protected secret information that was originally provided as a string.
	//
	// If you create this secret by using the Secrets Manager console then only the ``SecretString``
	// parameter contains data. Secrets Manager stores the information as a JSON structure of
	// key/value pairs that the Lambda rotation function knows how to parse.
	//
	// If you store custom information in the secret by using the CreateSecret , UpdateSecret , or
	// PutSecretValue API operations instead of the Secrets Manager console, or by using the
	// *Other secret type* in the console, then you must code your Lambda rotation function to
	// parse and interpret those values.
	SecretString string
	// VersionStages is a list of all of the staging labels currently attached to this version of the secret.
	VersionStages []string
	// Creates is the date and time that this version of the secret was created.
	CreatedDate *time.Time
}

// SecretsManager is a interface to interact with the
// green grass secrets manager.
type SecretsManager interface {
	// GetSecret calls the secrets manager lambda to obtain the requested secret value.
	//
	// The _secretID_ specifies the secret containing the version that you want to retrieve. You can specify either the
	// Amazon Resource Name (ARN) or the friendly name of the secret.
	//
	// _VersionID_ Specifies the unique identifier of the version of the secret that you want to retrieve. If you
	// specify this parameter then don't specify _VersionStage_. If you don't specify either a _VersionStage_ or
	// _SecretVersionId_ then the default is to perform the operation on the version with the _VersionStage_ value of _AWSCURRENT_.
	// This value is typically a UUID-type value with 32 hexadecimal digits.
	//
	// The _VersionStage_ Specifies the secret version that you want to retrieve by the staging label attached to the version.
	// Staging labels are used to keep track of different versions during the rotation process. If you use this parameter then
	// don't specify _SecretVersionId_. If you don't specify either a _VersionStage_ or _SecretVersionId_ , then the default
	// is to perform the operation on the version with the _VersionStage_ value of _AWSCURRENT_.
	GetSecret(secretID, versionID, versionStage string) (*Secret, error)
}

// SecretAPI encapsulates local secret manager communication.
type SecretAPI struct {
	APIRequest
	response string
}

// NewSecretAPI creates a new instance of the local api communication
// towards the secrets manager.
func NewSecretAPI() *SecretAPI {
	return &SecretAPI{}
}

// GetLastResponse returns the last success response from an API call.
func (s *SecretAPI) GetLastResponse() string {
	return s.response
}

// GetSecret calls the secrets manager lambda to obtain the requested secret value.
//
// The _secretID_ specifies the secret containing the version that you want to retrieve. You can specify either the
// Amazon Resource Name (ARN) or the friendly name of the secret.
//
// _VersionID_ Specifies the unique identifier of the version of the secret that you want to retrieve. If you
// specify this parameter then don't specify _VersionStage_. If you don't specify either a _VersionStage_ or
// _SecretVersionId_ then the default is to perform the operation on the version with the _VersionStage_ value of _AWSCURRENT_.
// This value is typically a UUID-type value with 32 hexadecimal digits.
//
// The _VersionStage_ Specifies the secret version that you want to retrieve by the staging label attached to the version.
// Staging labels are used to keep track of different versions during the rotation process. If you use this parameter then
// don't specify _SecretVersionId_. If you don't specify either a _VersionStage_ or _SecretVersionId_ , then the default
// is to perform the operation on the version with the _VersionStage_ value of _AWSCURRENT_.
func (s *SecretAPI) GetSecret(secretID, versionID, versionStage string) (*Secret, error) {

	if s.err != nil {
		return nil, s.err
	}

	cs := C.CString(secretID)
	cv := C.CString(versionID)
	ct := C.CString(versionStage)

	defer func() {
		s.close()
		C.free(unsafe.Pointer(cs))
		C.free(unsafe.Pointer(cv))
		C.free(unsafe.Pointer(ct))
	}()

	res := C.gg_request_result{}

	s.APIRequest.initialize()

	e := GreenGrassCode(C.gg_get_secret_value(s.request, cs, cv, ct, &res))

	s.response = s.handleRequestResponse(
		"Failed to get secret",
		e, RequestStatus(res.request_status),
	)

	if s.err != nil {
		return nil, s.err
	}

	Log(LogLevelInfo, "Got Secret: %s", s.response)

	var secret Secret
	if err := json.Unmarshal([]byte(s.response), &secret); err != nil {
		return nil, err
	}

	return &secret, nil

}
