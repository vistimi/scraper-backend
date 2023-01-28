package controller

import (
	interfaceAdapter "scraper-backend/src/adapter/interface"
	"scraper-backend/src/util"

	driverHost "scraper-backend/src/driver/host"
	driverDynamodb "scraper-backend/src/driver/database/dynamodb"
	driverBucket "scraper-backend/src/driver/storage/bucket"
)

func ConstructorPicture(cfg util.Config) interfaceAdapter.ControllerPicture {
	return &ControllerPicture{
		S3:         driverBucket.Constructor(cfg.AwsS3Client),
		BucketName: cfg.S3BucketNamePictures,
		DynamodbProcess: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureProcess.TableName,
			cfg.AwsDynamodbTablePictureProcess.PrimaryKey,
			*cfg.AwsDynamodbTablePictureProcess.SortKey,
		),
		DynamodbValidation: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureValidation.TableName,
			cfg.AwsDynamodbTablePictureValidation.PrimaryKey,
			*cfg.AwsDynamodbTablePictureValidation.SortKey,
		),
		DynamodbProduction: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureProduction.TableName,
			cfg.AwsDynamodbTablePictureProduction.PrimaryKey,
			*cfg.AwsDynamodbTablePictureProduction.SortKey,
		),
		DynamodbBlocked: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureBlocked.TableName,
			cfg.AwsDynamodbTablePictureBlocked.PrimaryKey,
			*cfg.AwsDynamodbTablePictureBlocked.SortKey,
		),
	}
}

func ConstructorTag(cfg util.Config, controllerPicture interfaceAdapter.ControllerPicture) interfaceAdapter.ControllerTag {
	return &ControllerTag{
		Dynamodb: driverDynamodb.ConstructorTag(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTableTag.TableName,
			cfg.AwsDynamodbTableTag.PrimaryKey,
			*cfg.AwsDynamodbTableTag.SortKey,
		),
		ControllerPicture: controllerPicture,
	}
}

func ConstructorUser(cfg util.Config) interfaceAdapter.ControllerUser {
	return &ControllerUser{
		Dynamodb: driverDynamodb.ConstructorUser(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTableUser.TableName,
			cfg.AwsDynamodbTableUser.PrimaryKey,
			*cfg.AwsDynamodbTableUser.SortKey,
		),
	}
}

func ConstructorFlickr(cfg util.Config, controllerPicture interfaceAdapter.ControllerPicture, controllerTag interfaceAdapter.ControllerTag, controllerUser interfaceAdapter.ControllerUser) interfaceAdapter.ControllerFlickr {
	return &ControllerFlickr{
		Api:               driverHost.ConstructorApiFlickr(),
		ControllerPicture: controllerPicture,
		ControllerTag:     controllerTag,
		ControllerUser:    controllerUser,
	}
}

func ConstructorPexels(cfg util.Config, controllerPicture interfaceAdapter.ControllerPicture, controllerTag interfaceAdapter.ControllerTag, controllerUser interfaceAdapter.ControllerUser) interfaceAdapter.ControllerPexels {
	return &ControllerPexels{
		Api:               driverHost.ConstructorApiPexels(),
		ControllerPicture: controllerPicture,
		ControllerTag:     controllerTag,
		ControllerUser:    controllerUser,
	}
}

func ConstructorUnsplash(cfg util.Config, controllerPicture interfaceAdapter.ControllerPicture, controllerTag interfaceAdapter.ControllerTag, controllerUser interfaceAdapter.ControllerUser) interfaceAdapter.ControllerUnsplash {
	return &ControllerUnsplash{
		Api:               driverHost.ConstructorApiUnsplash(),
		ControllerPicture: controllerPicture,
		ControllerTag:     controllerTag,
		ControllerUser:    controllerUser,
	}
}
