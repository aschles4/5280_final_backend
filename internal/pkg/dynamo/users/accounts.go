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
	// create the api params
	params := &dynamodb.ScanInput{
		TableName:        aws.String("Accounts"),
		FilterExpression: aws.String("email = :e"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":e": {S: aws.String(email)},
		},
	}
	// read the item
	resp, err := s.db.Scan(params)
	if err != nil {
		return nil, err
	}

	if *resp.Count >= int64(1) {
		var act UserAccount
		err = dynamodbattribute.UnmarshalMap(resp.Items[0], &act)
		if err != nil {
			return nil, err
		}
		return &act, nil
	}

	return nil, nil
}

func (s Store) FindUserAccountByToken(ctx context.Context, token string) (*UserAccount, error) {

	// create the api params
	params := &dynamodb.ScanInput{
		TableName:        aws.String("Accounts"),
		FilterExpression: aws.String("tkn = :t"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {S: aws.String(token)},
		},
	}
	// read the item
	resp, err := s.db.Scan(params)
	if err != nil {
		return nil, err
	}

	if *resp.Count >= int64(1) {
		var act UserAccount
		err = dynamodbattribute.UnmarshalMap(resp.Items[0], &act)
		if err != nil {
			return nil, err
		}
		return &act, nil
	}

	return nil, nil
}

func (s Store) RemoveUserAccountByID(ctx context.Context, ID string) error {
	return s.deleteFromTableByID(ID, "Accounts")
}

func (s Store) UpdateUserAccountToken(ctx context.Context, ID, token string) (*UserAccount, error) {
	// create the api params
	params := &dynamodb.UpdateItemInput{
		TableName: aws.String("Accounts"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(ID),
			},
		},
		UpdateExpression: aws.String("set tkn = :t"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {S: aws.String(token)},
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
