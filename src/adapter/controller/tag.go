package controller

import dynamodbTable "scraper-backend/src/driver/database/dynamodb/table"

type controllerTag struct {
	TablePicture *dynamodbTable.TableTag
}
