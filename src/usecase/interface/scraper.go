package entity

type UsecaseScraperFlickr interface{
	SearchPhotos(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsSearchPhotoFlickr) ([]primitive.ObjectID, error)
}

type UsecaseScraperPexels interface{
	SearchPhotos(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsSearchPhotoFlickr) ([]primitive.ObjectID, error)
}

type UsecaseScraperUnsplash interface {
	SearchPhotos(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsSearchPhotoUnsplash) ([]primitive.ObjectID, error)
}