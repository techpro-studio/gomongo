package gomongo

import (
	"fmt"
	"github.com/techpro-studio/gohttplib"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"reflect"
	"strings"
)

type Mapper[T any, U any] func(input T) U

func SliceMap[T any, U any](values []T, mapper Mapper[T, U]) []U {
	result := make([]U, len(values))
	for i, value := range values {
		result[i] = mapper(value)
	}
	return result
}

func Int64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}
	v64 := int64(*v)
	return &v64
}

func IsZeroValue(v reflect.Value) bool {
	return v.IsZero()
}

func BuildUpdateDoc(v interface{}) bson.M {
	update := bson.M{}
	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		bsonTag := field.Tag.Get("bson")
		if commaIndex := strings.Index(bsonTag, ","); commaIndex != -1 {
			bsonTag = bsonTag[:commaIndex]
		}

		// Check for zero values (e.g., empty strings, zero ints)
		if bsonTag != "" && !IsZeroValue(value) {
			update[bsonTag] = value.Interface()
		}
	}

	return update
}

func ExtractCorrectObjectIdFromRequest(r *http.Request, key string, allowNull bool) (*string, error) {
	id := gohttplib.GetParameterFromURLInRequest(r, key)
	if id == nil {
		if allowNull {
			return nil, nil
		} else {
			return nil, gohttplib.NewServerError(400, "REQUIRED", "REQUIRED", key, nil)
		}
	}
	return GetValidObjectId(*id, key)
}

func GetValidObjectIdFromMap(body map[string]any, key string) (*string, error) {
	id, ok := body[key].(string)
	if !ok {
		return nil, gohttplib.NewServerError(400, "INVALID_OBJECT_ID", fmt.Sprintf("%s is not event string", id), key, nil)
	}
	return GetValidObjectId(id, key)
}

func GetValidObjectId(id string, key string) (*string, error) {
	_, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, gohttplib.NewServerError(400, "INVALID_OBJECT_ID", fmt.Sprintf("%s is not a valid ID", id), key, nil)
	}
	return &id, nil
}
