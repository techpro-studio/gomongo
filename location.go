package gomongo

import (
	"github.com/techpro-studio/gohttplib/location"
	"go.mongodb.org/mongo-driver/bson"
)

func LocationParametersToMongoQuery(parameters *location.LocationParameters) bson.M {
	if parameters == nil {
		return nil
	}
	return bson.M{"$near": bson.M{
		"$geometry":    bson.M{"type": "Point", "coordinates": []float64{parameters.Longitude, parameters.Latitude}},
		"$maxDistance": parameters.MaxDistance,
		"$minDistance": parameters.MinDistance,
	}}
}
