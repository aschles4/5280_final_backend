package users

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (s Store) CreateUserProfile(ctx context.Context, profile UserProfile) error {
	return s.writeToTable(profile, "Profiles")
}

func (s Store) FindUserProfileByID(ctx context.Context, ID string) (*UserProfile, error) {
	return s.readFromProfilesByType("id", ID)
}

func (s Store) RemoveUserProfileByID(ctx context.Context, ID string) error {
	return s.deleteFromTableByID(ID, "Profiles")
}

//DynamoDB Helpers
func (s Store) readFromProfilesByType(typ, val string) (*UserProfile, error) {
	// create the api params
	params := &dynamodb.GetItemInput{
		TableName: aws.String("Profiles"),
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

	var prof UserProfile
	err = dynamodbattribute.UnmarshalMap(resp.Item, &prof)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}
