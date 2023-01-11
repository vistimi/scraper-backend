package controller

import dynamodbTable "scraper-backend/src/driver/database/dynamodb/table"

type controllerUser struct {
	TablePicture *dynamodbTable.TableUser
}
