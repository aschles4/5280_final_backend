package users

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Store struct {
	db    *dynamodb.DynamoDB
}

func NewStore(conn, region string) (*Store, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(conn),
	})
	if err != nil {
		return nil, err
	}

	db := dynamodb.New(sess)

	return &Store{
		db:    db,
	}, nil
}