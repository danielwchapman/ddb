package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type options struct {
	updates         expression.UpdateBuilder
	updatesCount    int
	conditions      expression.ConditionBuilder
	conditionsCount int
	returnValues    types.ReturnValue
	returnValuesOut any
}

type Option func(options *options) error

func WithFieldUpdates(updates map[string]any) Option {
	return func(options *options) error {
		item, err := attributevalue.MarshalMap(updates)
		if err != nil {
			return fmt.Errorf("WithFieldUpdates: MarshalMap: %w", err)
		}

		for k, v := range item {
			options.updates = options.updates.Set(expression.Name(k), expression.Value(v))
		}

		options.updatesCount += len(item)

		return nil
	}
}

func WithItemExists() Option {
	return func(options *options) error {
		cond := expression.AttributeExists(expression.Name("PK"))

		if options.conditionsCount == 0 {
			options.conditions = cond
		} else {
			options.conditions.And(cond)
		}
		options.conditionsCount++

		return nil
	}
}

func WithItemNotExist() Option {
	return func(options *options) error {
		cond := expression.AttributeNotExists(expression.Name("PK"))

		if options.conditionsCount == 0 {
			options.conditions = cond
		} else {
			options.conditions.And(cond)
		}
		options.conditionsCount++

		return nil
	}
}

func WithReturnValues(returnValues types.ReturnValue, out any) Option {
	return func(options *options) error {
		options.returnValues = returnValues
		options.returnValuesOut = out
		return nil
	}
}
