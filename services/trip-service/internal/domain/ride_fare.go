package domain

import (
	"slices"

	tripTypes "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string // ex. van, luxury, sedan
	TotalPriceInCents float64
	Route             *tripTypes.OsrmAPIResponse
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func RideFareModelsToProtos(fares []*RideFareModel) []*pb.RideFare {
	return slices.Collect(func(yield func(*pb.RideFare) bool) {
		for _, f := range fares {
			if !yield(f.ToProto()) {
				return
			}
		}
	})
}
