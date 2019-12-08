package users

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type StreamAccounts struct {
	ID int64 `json:"id"`
}

type Content struct {
	ID int64 `json:"id"`
}

type Library struct {
	ContentList []Content `json:"contentList"`
}

type UserAccount struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type UserProfile struct {
	ID             string           `json:"Id"`
	Name           string           `json:"name"`
	Email          string           `json:"email"`
	StreamAccounts []StreamAccounts `json:"StreamAccounts"`
	Library        Library          `json:"library"`
}

//DynamoDB Helpers
func (s Store) deleteFromTableByID(ID, table string) error {
	// create the api params
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(ID),
			},
		},
	}

	// delete the item
	_, err := s.db.DeleteItem(params)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) writeToTable(doc interface{}, table string) error {
	m, err := dynamodbattribute.MarshalMap(doc)
	if err != nil {
		return err
	}

	// create the api params
	params := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      m,
	}

	// put the item
	_, err = s.db.PutItem(params)
	if err != nil {
		return err
	}

	return nil
}
