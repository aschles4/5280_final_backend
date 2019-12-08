package users

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (s Store) CreateUserAccount(ctx context.Context, account UserAccount) error {
	return s.writeToTable(account, "Accounts")
}

func (s Store) FindUserAccountByEmail(ctx context.Context, email string) (*UserAccount, error) {
	return s.readFromAccountsByType("email", email)
}

func (s Store) FindUserAccountByToken(ctx context.Context, token string) (*UserAccount, error) {
	return s.readFromAccountsByType("token", token)
}

func (s Store) RemoveUserAccountByID(ctx context.Context, ID string) error {
	return s.deleteFromTableByID(ID, "Accounts")
}

func (s Store) UpdateUserAccountToken(ctx context.Context, ID, token string) (*UserAccount, error) {
	// create the api params
	params := &dynamodb.UpdateItemInput{
		TableName: aws.String("Accounts"),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(ID),
			},
		},
		UpdateExpression: aws.String("set token=:t"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {N: aws.String(token)},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueAllNew),
	}

	// update the item
	resp, err := s.db.UpdateItem(params)
	if err != nil {
		return nil, err
	}

	// unmarshal the dynamodb attribute values into a custom struct
	var updatedAct UserAccount
	err = dynamodbattribute.UnmarshalMap(resp.Attributes, &updatedAct)
	if err != nil {
		return nil, err
	}

	return &updatedAct, nil
}

//DynamoDB Helpers
func (s Store) readFromAccountsByType(typ, val string) (*UserAccount, error) {
	// create the api params
	params := &dynamodb.GetItemInput{
		TableName: aws.String("Accounts"),
		Key: map[string]*dynamodb.AttributeValue{
			typ: {
				S: aws.String(val),
			},
		},
	}

	// read the item
	resp, err := s.db.GetItem(params)
	if err != nil {
		return nil, err
	}

	var act UserAccount
	err = dynamodbattribute.UnmarshalMap(resp.Item, &act)
	if err != nil {
		return nil, err
	}

	return &act, nil
}
