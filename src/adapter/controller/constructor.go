package controller

import (
	interfaceAdapter "scraper-backend/src/adapter/interface"
	"scraper-backend/src/util"

	driverDynamodb "scraper-backend/src/driver/database/dynamodb"
	driverHost "scraper-backend/src/driver/host"
	driverBucket "scraper-backend/src/driver/storage/bucket"
)

func ConstructorPicture(cfg util.Config) interfaceAdapter.ControllerPicture {
	return &ControllerPicture{
		S3:         driverBucket.Constructor(cfg.AwsS3Client),
		BucketName: cfg.S3BucketNamePictures,
		DynamodbProcess: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureProcess.TableName,
			cfg.AwsDynamodbTablePictureProcess.PrimaryKeyName,
			cfg.AwsDynamodbTablePictureProcess.PrimaryKeyType,
			*cfg.AwsDynamodbTablePictureProcess.SortKeyName,
			*cfg.AwsDynamodbTablePictureProcess.SortKeyType,
		),
		DynamodbValidation: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureValidation.TableName,
			cfg.AwsDynamodbTablePictureValidation.PrimaryKeyName,
			cfg.AwsDynamodbTablePictureValidation.PrimaryKeyType,
			*cfg.AwsDynamodbTablePictureValidation.SortKeyName,
			*cfg.AwsDynamodbTablePictureValidation.SortKeyType,
		),
		DynamodbProduction: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureProduction.TableName,
			cfg.AwsDynamodbTablePictureProduction.PrimaryKeyName,
			cfg.AwsDynamodbTablePictureProduction.PrimaryKeyType,
			*cfg.AwsDynamodbTablePictureProduction.SortKeyName,
			*cfg.AwsDynamodbTablePictureProduction.SortKeyType,
		),
		DynamodbBlocked: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTablePictureBlocked.TableName,
			cfg.AwsDynamodbTablePictureBlocked.PrimaryKeyName,
			cfg.AwsDynamodbTablePictureBlocked.PrimaryKeyType,
			*cfg.AwsDynamodbTablePictureBlocked.SortKeyName,
			*cfg.AwsDynamodbTablePictureBlocked.SortKeyType,
		),
	}
}

func ConstructorTag(cfg util.Config, controllerPicture interfaceAdapter.ControllerPicture) interfaceAdapter.ControllerTag {
	return &ControllerTag{
		Dynamodb: driverDynamodb.ConstructorTag(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTableTag.TableName,
			cfg.AwsDynamodbTableTag.PrimaryKeyName,
			cfg.AwsDynamodbTableTag.PrimaryKeyType,
			*cfg.AwsDynamodbTableTag.SortKeyName,
			*cfg.AwsDynamodbTableTag.SortKeyType,
		),
		ControllerPicture: controllerPicture,
	}
}

func ConstructorUser(cfg util.Config) interfaceAdapter.ControllerUser {
	return &ControllerUser{
		Dynamodb: driverDynamodb.ConstructorUser(
			cfg.AwsDynamodbClient,
			cfg.AwsDynamodbTableUser.TableName,
			cfg.AwsDynamodbTableUser.PrimaryKeyName,
			cfg.AwsDynamodbTableUser.PrimaryKeyType,
			*cfg.AwsDynamodbTableUser.SortKeyName,
			*cfg.AwsDynamodbTableUser.SortKeyType,
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
