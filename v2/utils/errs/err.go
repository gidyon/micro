package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromJSONMarshal wraps error returned from json.Marshal to a status error
func FromJSONMarshal(err error, obj string) error {
	return status.Errorf(codes.Internal, "failed to json marshal %s: %v", obj, err)
}

// FromJSONUnMarshal wraps error returned from json.Unmarshal to a status error
func FromJSONUnMarshal(err error, obj string) error {
	return status.Errorf(codes.Internal, "failed to json unmarshal %s: %v", obj, err)
}

// FromProtoMarshal wraps error returned from proto.Marshal to a status error
func FromProtoMarshal(err error, obj string) error {
	return status.Errorf(codes.Internal, "failed to proto marshal %s: %v", obj, err)
}

// FromProtoUnMarshal wraps error returned from proto.Unmarshal to a status error
func FromProtoUnMarshal(err error, obj string) error {
	return status.Errorf(codes.Internal, "failed to proto unmarshal %s: %v", obj, err)
}

// MissingField returns a status error caused by a missing message field
func MissingField(field string) error {
	return status.Errorf(codes.InvalidArgument, "missing message field: %v", field)
}

// DuplicateField returns a status error for a duplicate field
func DuplicateField(fieldName, fieldValue string) error {
	return status.Errorf(codes.AlreadyExists, "%s with value %s exists", fieldName, fieldValue)
}

// NilObject is error resulting from using nil references to objects
func NilObject(object string) error {
	return status.Errorf(codes.InvalidArgument, "nil object not allowed: %s", object)
}

// ConvertingType wraps error that occured during type assertion to grpc status error
func ConvertingType(err error, from, to string) error {
	return status.Errorf(codes.Internal, "couldn't convert from %s to %s: %v", from, to, err)
}

// IncorrectVal returns a status error indicating val was incorrect
func IncorrectVal(val string) error {
	return status.Errorf(codes.InvalidArgument, "incorrect value for %q", val)
}

// WriteFailed returns a status error for a write operation error
func WriteFailed(err error) error {
	return status.Errorf(codes.Internal, "write operation failed: %v", err)
}

// ReadFailed returns a status error for a read operation error
func ReadFailed(err error) error {
	return status.Errorf(codes.Internal, "read operation failed: %v", err)
}

// SQLQueryFailed wraps a sql query error to a status error
func SQLQueryFailed(err error, queryType string) error {
	return status.Errorf(codes.Internal, "failed to execute %s query: %v", queryType, err)
}

// DoesNotExist returns status error indicating that the resource does not exist
func DoesNotExist(resource, id string) error {
	return status.Errorf(codes.NotFound, "%s with id %s does not exist", resource, id)
}

// DoesExist returns status error indicating the resource does exist
func DoesExist(resource, id string) error {
	return status.Errorf(codes.AlreadyExists, "%s with id %s already exists", resource, id)
}

// FailedToSave is status error returned from failed save operation
func FailedToSave(resource string, err error) error {
	return status.Errorf(codes.Internal, "failed to save %s: %v", resource, err)
}

// FailedToUpdate wraps error returned from failed update operation to a status error
func FailedToUpdate(resource string, err error) error {
	return status.Errorf(codes.Internal, "failed to update %s: %v", resource, err)
}

// FailedToDelete wraps error returned from failed delete operation to a status error
func FailedToDelete(resource string, err error) error {
	return status.Errorf(codes.Internal, "failed to delete %s: %v", resource, err)
}

// FailedToFind wraps error returned from failed find operation to a status error
func FailedToFind(resource string, err error) error {
	return status.Errorf(codes.Internal, "failed to find %s: %v", resource, err)
}

// FailedToBeginTx wraps error returned from failed transaction to a status error
func FailedToBeginTx(err error) error {
	return status.Errorf(codes.Internal, "failed to begin db transaction: %v", err)
}

// FailedToCommitTx wraps error returned from failed commit of transaction to a status error
func FailedToCommitTx(err error) error {
	return status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
}

// FailedToRollbackTx wraps error returned from failed rollback of transaction to a status error
func FailedToRollbackTx(err error) error {
	return status.Errorf(codes.Internal, "failed to rollback transaction: %v", err)
}

// TokenCredentialNotMatching creates a status error caused by mismatch in token credential
func TokenCredentialNotMatching(cred string) error {
	return status.Errorf(codes.PermissionDenied, "token credential %v do not match", cred)
}

// FailedToEncrypt is status error from failed encryption operation
func FailedToEncrypt(err error) error {
	return status.Errorf(codes.Internal, "failed to encrypt data: %v", err)
}

// FailedToDecrypt is status error from failed decryption operation
func FailedToDecrypt(err error) error {
	return status.Errorf(codes.Internal, "failed to decrypt data: %v", err)
}

// RedisCmdFailed wraps error returned from a failed redis command to a status error
func RedisCmdFailed(err error, command string) error {
	return status.Errorf(codes.Internal, "failed to execute %s command: %v", command, err)
}

// FailedToExecuteTemplate returns a status error for a failed template execution
func FailedToExecuteTemplate(err error) error {
	return status.Errorf(codes.Internal, "failed to execute template: %v", err)
}

// FailedToGetToken wraps the error returned while getting token to a status error
func FailedToGetToken(err error) error {
	return status.Errorf(codes.Internal, "failed to get token: %v", err)
}

// WrapErrorWithCode is a wraps generic error to a status error with provided code
func WrapErrorWithCode(code codes.Code, err error) error {
	return status.Error(code, err.Error())
}

// WrapError is a wraps generic error to a status error
func WrapError(err error) error {
	return status.Error(status.Code(err), err.Error())
}

// WrapErrorWithCodeAndMsg wraps generic error to a status error with provided code and msg
func WrapErrorWithCodeAndMsg(code codes.Code, err error, msg string) error {
	return status.Errorf(code, "%s: %v", msg, err.Error())
}

// WrapErrorWithCodeAndMsgFunc is a common message wrapper for WrapErrorWithCodeAndMsg
func WrapErrorWithCodeAndMsgFunc(msg string) func(codes.Code, error) error {
	return func(code codes.Code, err error) error {
		if err != nil {
			return WrapErrorWithCodeAndMsg(code, err, msg)
		}
		return nil
	}
}

// WrapErrorWithMsg is a wraps generic error to a status error with code and msg formt
func WrapErrorWithMsg(err error, msg string) error {
	return status.Errorf(status.Code(err), "%s: %v", msg, err.Error())
}

// WrapErrorWithMsgFunc is a common message wrapper for WrapErrorWithMsg
func WrapErrorWithMsgFunc(msg string) func(error) error {
	return func(err error) error {
		if err != nil {
			return WrapErrorWithMsg(err, msg)
		}
		return nil
	}
}

// WrapMessage is a wraps message provided to a status error
func WrapMessage(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

// WrapMessagef is a wraps message provided to a status error
func WrapMessagef(code codes.Code, format string, args ...interface{}) error {
	return status.Error(code, fmt.Sprintf(format, args...))
}
