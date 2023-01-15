package controller

import (
	interfaceAdapter "scraper-backend/src/adapter/interface"
	"scraper-backend/src/util"

	driverDynamodb "scraper-backend/src/driver/database/dynamodb"
	driverBucket "scraper-backend/src/driver/storage/bucket"
)

func ContructorPicture(cfg util.Config) interfaceAdapter.ControllerPicture {
	return &ControllerPicture{
		S3:         driverBucket.Constructor(cfg.AwsS3Client),
		BucketName: cfg.S3BucketNamePictures,
		DynamodbProcess: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamoDbClient,
			cfg.TablePictureProcessName,
			cfg.TablePicturePrimaryKey,
			cfg.TablePictureSortKey,
		),
		DynamodbValidation: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamoDbClient,
			cfg.TablePictureValidationName,
			cfg.TablePicturePrimaryKey,
			cfg.TablePictureSortKey,
		),
		DynamodbProduction: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamoDbClient,
			cfg.TablePictureProductionName,
			cfg.TablePicturePrimaryKey,
			cfg.TablePictureSortKey,
		),
		DynamodbBlocked: driverDynamodb.ConstructorPicture(
			cfg.AwsDynamoDbClient,
			cfg.TablePictureBlockedName,
			cfg.TablePicturePrimaryKey,
			cfg.TablePictureSortKey,
		),
	}
}

func ContructorTag(cfg util.Config, controllerPicture ControllerPicture) interfaceAdapter.ControllerTag {
	return &ControllerTag{
		Dynamodb: driverDynamodb.ConstructorTag(
			cfg.AwsDynamoDbClient,
			cfg.TableTagName,
			cfg.TableTagPrimaryKey,
			cfg.TableTagSortKey,
		),
		ControllerPicture: controllerPicture,
	}
}

func ContructorUser(cfg util.Config, controllerPicture ControllerPicture) interfaceAdapter.ControllerUser {
	return &ControllerUser{
		Dynamodb: driverDynamodb.ConstructorUser(
			cfg.AwsDynamoDbClient,
			cfg.TableUserName,
			cfg.TableUserPrimaryKey,
			cfg.TableUserSortKey,
		),
	}
}
