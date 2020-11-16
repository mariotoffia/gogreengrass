package sdk

/*
#include <stdlib.h>
typedef void (*get_secret) (char *ctx, char* id, char *version, char *stage);

static inline void call_c_get_secret(get_secret ptr, char *ctx, char* id, char *version, char *stage) {
	(ptr)(ctx, id, version, stage);
}
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

// Secret represents a single secret obtained from the
// green grass secrets manager.
type Secret struct {
	// ARN is the ARN of the secret.
	ARN string `json:"arn"`
	// Name is the friendly name of the secret.
	Name string `json:"name"`
	// VersionID is the unique identifier of this version of the secret.
	VersionID string `json:"version"`
	// SecretBinary is the decrypted part of the protected secret information that was originally provided as
	// binary data in the form of a byte array. This parameter is not used if the secret is created by the Secrets Manager console.
	//
	// If you store custom information in this field of the secret, then you must code your Lambda rotation function to parse and
	// interpret whatever you store in the _SecretString_ or _SecretBinary_ fields.
	SecretBinary []byte `json:"bin"`
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
	SecretString string `json:"secret"`
	// VersionStages is a list of all of the staging labels currently attached to this version of the secret.
	Stages []string `json:"stages"`
	// Creates is the date and time that this version of the secret was created.
	Created time.Time `json:"created"`
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

type ggSecretsManager struct {
	fnGetSecret *[0]byte
}

// NewSecretsManagerClient creates a new secrets manager client to communicate with
// the greengrass dataplane API.
func NewSecretsManagerClient(fnGetSecret *[0]byte) SecretsManager {
	return &ggSecretsManager{fnGetSecret: fnGetSecret}
}

func (sm *ggSecretsManager) GetSecret(secretID, versionID, versionStage string) (*Secret, error) {

	if nil == sm.fnGetSecret {

		return nil,
			fmt.Errorf(
				"Could not fetch secret (id: %s, version: %s stage: %s) since no func defined",
				secretID, versionID, versionStage,
			)

	}

	ctx := getContext()

	cs := C.CString(secretID)
	cv := C.CString(versionID)
	ct := C.CString(versionStage)
	cc := C.CString(ctx)

	defer func() {
		C.free(unsafe.Pointer(cs))
		C.free(unsafe.Pointer(cv))
		C.free(unsafe.Pointer(ct))
		C.free(unsafe.Pointer(cc))
	}()

	C.call_c_get_secret(sm.fnGetSecret, cc, cs, cv, ct)

	buf := getProcessBuffer(ctx)
	if "" == buf {
		return nil, nil /*TODO: error*/
	}

	var secret Secret
	err := json.Unmarshal([]byte(buf), &secret)
	fmt.Println(buf)

	return &secret, err

}
