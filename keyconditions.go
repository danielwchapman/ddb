package ddb

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

type KeyCondition func(pkColName, skColName string) expression.KeyConditionBuilder

func KeyPkOnly(pk string) KeyCondition {
	return func(pkColName, skColName string) expression.KeyConditionBuilder {
		return expression.Key(pkColName).Equal(expression.Value(pk))
	}
}

func KeySkBeginsWith(pk, skPrefix string) KeyCondition {
	return func(pkColName, skColName string) expression.KeyConditionBuilder {
		return expression.Key(pkColName).Equal(expression.Value(pk)).
			And(expression.Key(skColName).BeginsWith(skPrefix))
	}
}

func KeySkBetween(pk, skStart, skEnd string) KeyCondition {
	return func(pkColName, skColName string) expression.KeyConditionBuilder {
		return expression.Key(pkColName).Equal(expression.Value(pk)).
			And(expression.Key(skColName).Between(expression.Value(skStart), expression.Value(skEnd)))
	}
}

func KeySkGreaterThan(pk, skStart string) KeyCondition {
	return func(pkColumnName, skColumnName string) expression.KeyConditionBuilder {
		return expression.Key(pkColumnName).Equal(expression.Value(pk)).
			And(expression.Key(skColumnName).GreaterThan(expression.Value(skStart)))
	}
}

func KeySkGreaterThanEqual(pk, skStart string) KeyCondition {
	return func(pkColumnName, skColumnName string) expression.KeyConditionBuilder {
		return expression.Key(pkColumnName).Equal(expression.Value(pk)).
			And(expression.Key(skColumnName).GreaterThanEqual(expression.Value(skStart)))
	}
}

func KeySkLessThan(pk, skEnd string) KeyCondition {
	return func(pkColumnName, skColumnName string) expression.KeyConditionBuilder {
		return expression.Key(pkColumnName).Equal(expression.Value(pk)).
			And(expression.Key(skColumnName).LessThan(expression.Value(skEnd)))
	}
}

func KeySkLessThanEqual(pk, skEnd string) KeyCondition {
	return func(pkColumnName, skColumnName string) expression.KeyConditionBuilder {
		return expression.Key(pkColumnName).Equal(expression.Value(pk)).
			And(expression.Key(skColumnName).LessThanEqual(expression.Value(skEnd)))
	}
}
