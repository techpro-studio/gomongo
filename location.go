package gomongo

import (
	"encoding/json"
	"github.com/techpro-studio/gohttplib"
	"github.com/techpro-studio/gohttplib/validator"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

type LocationParameters struct {
	Latitude    float64
	Longitude   float64
	MinDistance int64
	MaxDistance int64
}

func NewObjectLocation(latitude, longitude float64) *ObjectLocation {
	return &ObjectLocation{"Point", []float64{longitude, latitude}}
}

func LocationParametersToMongoQuery(parameters *LocationParameters) bson.M {
	if parameters == nil {
		return nil
	}
	return bson.M{"$near": bson.M{
		"$geometry":    bson.M{"type": "Point", "coordinates": []float64{parameters.Longitude, parameters.Latitude}},
		"$maxDistance": parameters.MaxDistance,
		"$minDistance": parameters.MinDistance,
	}}
}

type ObjectLocation struct {
	Type        string
	Coordinates []float64
}

func (self ObjectLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{
		Latitude:  self.Coordinates[1],
		Longitude: self.Coordinates[0],
	})
}

func LongitudeValidators(key string) []validator.Validator {
	upper := 180.0
	bottom := -180.0
	return validator.RequiredFloatValidators(key, validator.FloatInRangeValidator(key, validator.FloatRange{&upper, &bottom}))
}

func LatitudeValidators(key string) []validator.Validator {
	upper := 90.0
	bottom := -90.0
	return validator.RequiredFloatValidators("latitude", validator.FloatInRangeValidator("longitude", validator.FloatRange{&upper, &bottom}))
}

func DistanceValidator() validator.Validator {
	var bottom int64 = 0
	var upper int64 = 32000000000
	return validator.Int64InRangeValidator("max_distance", validator.Int64Range{Bottom: &bottom, Upper: &upper})
}

func ParseGeoLocation(longitudeRaw string, latitudeRaw string) (float64, float64, error) {
	longitude, err := strconv.ParseFloat(longitudeRaw, 64)
	if err != nil {
		return 0, 0, gohttplib.NewServerError(400, "longitude", "INVALID_FLOAT", "INVALID_FLOAT", nil)
	}
	latitude, err := strconv.ParseFloat(latitudeRaw, 64)
	if err != nil {
		return 0, 0, gohttplib.NewServerError(400, "longitude", "INVALID_FLOAT", "INVALID_FLOAT", nil)
	}

	errs := validator.ValidateValue(longitude, LongitudeValidators("longitude"))
	if len(errs) > 0 {
		return 0, 0, gohttplib.ServerError{StatusCode: 400, Errors: gohttplib.Errors{Errors: errs}}
	}
	errs = validator.ValidateValue(latitude, LatitudeValidators("latitude"))
	if len(errs) > 0 {
		return 0, 0, gohttplib.ServerError{StatusCode: 400, Errors: gohttplib.Errors{Errors: errs}}
	}
	return longitude, latitude, nil
}

func LocationParametersFromRequest(req *http.Request, defaultMaxDistance int64) (*LocationParameters, error) {
	longitudeRaw := gohttplib.GetParameterFromURLInRequest(req, "longitude")
	latitudeRaw := gohttplib.GetParameterFromURLInRequest(req, "latitude")
	if longitudeRaw == nil || latitudeRaw == nil {
		return nil, nil
	}

	longitude, latitude, err := ParseGeoLocation(*longitudeRaw, *latitudeRaw)

	var minDistance int64 = 0
	maxDistance := defaultMaxDistance

	maxDistanceRaw := gohttplib.GetParameterFromURLInRequest(req, "max_distance")
	minDistanceRaw := gohttplib.GetParameterFromURLInRequest(req, "min_distance")

	if maxDistanceRaw != nil {

		maxDistance, err = strconv.ParseInt(*maxDistanceRaw, 10, 64)
		if err != nil {
			return nil, gohttplib.NewServerError(400, "max_distance", "INVALID_INT", "INVALID_INT", nil)
		}

		errs := validator.ValidateValue(maxDistance, validator.RequiredIntValidators("max_distance", DistanceValidator()))
		if len(errs) > 0 {
			return nil, gohttplib.ServerError{StatusCode: 400, Errors: gohttplib.Errors{Errors: errs}}
		}
	}

	if minDistanceRaw != nil {

		minDistance, err = strconv.ParseInt(*minDistanceRaw, 10, 64)
		if err != nil {
			return nil, gohttplib.NewServerError(400, "min_distance", "INVALID_INT", "INVALID_INT", nil)
		}

		errs := validator.ValidateValue(minDistance, validator.RequiredIntValidators("min_distance", DistanceValidator()))
		if len(errs) > 0 {
			return nil, gohttplib.ServerError{StatusCode: 400, Errors: gohttplib.Errors{Errors: errs}}
		}
	}

	if minDistance > maxDistance {
		return nil, gohttplib.NewServerError(400, "min_distance", "Min distance more tham max", "MIN_MAX_ERROR", nil)
	}

	return &LocationParameters{Latitude: latitude, Longitude: longitude, MinDistance: minDistance, MaxDistance: maxDistance}, nil

}
