package cel

import (
	"fmt"

	celgo "github.com/google/cel-go/cel"
	"k8s.io/apimachinery/pkg/api/resource"
	dracel "k8s.io/dynamic-resource-allocation/cel"
)

// ComparisonOperator defines available operators for CEL expression.
type ComparisonOperator string

const (
	OpEqual          ComparisonOperator = "=="
	OpNotEqual       ComparisonOperator = "!="
	OpGreater        ComparisonOperator = ">"
	OpGreaterOrEqual ComparisonOperator = ">="
	OpLess           ComparisonOperator = "<"
	OpLessOrEqual    ComparisonOperator = "<="
)

var boolOperators = map[ComparisonOperator]string{
	OpEqual:    "==",
	OpNotEqual: "!=",
}

var intOperators = map[ComparisonOperator]string{
	OpEqual:          "==",
	OpNotEqual:       "!=",
	OpGreater:        ">",
	OpGreaterOrEqual: ">=",
	OpLess:           "<",
	OpLessOrEqual:    "<=",
}

var stringOperators = map[ComparisonOperator]string{
	OpEqual:    "==",
	OpNotEqual: "!=",
}

var semverOperators = map[ComparisonOperator]string{
	OpEqual:          "== 0",
	OpNotEqual:       "!= 0",
	OpGreater:        "> 0",
	OpGreaterOrEqual: ">= 0",
	OpLess:           "< 0",
	OpLessOrEqual:    "<= 0",
}

var quantityOperators = map[ComparisonOperator]string{
	OpEqual:          "== 0",
	OpNotEqual:       "!= 0",
	OpGreater:        "> 0",
	OpGreaterOrEqual: ">= 0",
	OpLess:           "< 0",
	OpLessOrEqual:    "<= 0",
}

// ValueType defines available value types for CEL.
type ValueType string

const (
	TypeBool     ValueType = "bool"
	TypeString   ValueType = "string"
	TypeInt      ValueType = "int"
	TypeSemver   ValueType = "semver"
	TypeQuantity ValueType = "quantity"
	TypeUnknown  ValueType = "unknown"
)

func (vt ValueType) CELType() *celgo.Type {
	switch vt {
	case TypeBool:
		return celgo.BoolType
	case TypeInt:
		return celgo.IntType
	case TypeString, TypeSemver, TypeQuantity:
		return celgo.StringType
	}
	return nil
}

func buildBoolExpr(key string, op ComparisonOperator, value interface{}) (string, error) {
	b, ok := value.(bool)
	if !ok {
		return "", nil
	}
	comp, ok := boolOperators[op]
	if !ok {
		return "", fmt.Errorf("invalid operator %q for bool type", op)
	}
	return fmt.Sprintf("%s %s %t", key, comp, b), nil
}

func buildIntExpr(key string, op ComparisonOperator, value interface{}) (string, error) {
	num, ok := value.(int)
	if !ok {
		return "", nil
	}
	comp, ok := intOperators[op]
	if !ok {
		return "", fmt.Errorf("invalid operator %q for int type", op)
	}
	return fmt.Sprintf("%s %s %d", key, comp, num), nil
}

func buildStringExpr(key string, op ComparisonOperator, value interface{}) (string, error) {
	str, ok := value.(string)
	if !ok {
		return "", nil
	}
	comp, ok := stringOperators[op]
	if !ok {
		return "", fmt.Errorf("invalid operator %q for string type", op)
	}
	return fmt.Sprintf("%s %s %q", key, comp, str), nil
}

func buildSemverExpr(key string, op ComparisonOperator, value interface{}) (string, error) {
	str, ok := value.(string)
	if !ok {
		return "", nil
	}
	comp, ok := semverOperators[op]
	if !ok {
		return "", fmt.Errorf("invalid operator %q for semver type", op)
	}
	return fmt.Sprintf("(%s).compareTo(semver(%q)) %s", key, str, comp), nil
}

func buildQuantityExpr(key string, op ComparisonOperator, value interface{}) (string, error) {
	quantity, ok := value.(*resource.Quantity)
	if !ok {
		return "", nil
	}
	comp, ok := quantityOperators[op]
	if !ok {
		return "", fmt.Errorf("invalid operator %q for quantity type", op)
	}
	return fmt.Sprintf("(%s).compareTo(quantity(%q)) %s", key, quantity.String(), comp), nil
}

// BuildExpr returns a CEL expression given key, operator, value, and type.
// Examples:
// BuildExpr("foo", OpEqual, true, TypeBool) => "foo == true"
// BuildExpr("foo", OpEqual, "bar", TypeString) => "foo == \"bar\""
// BuildExpr("count", OpGreater, 5, TypeInt) => "count > 5"
// BuildExpr("ver", OpGreater, "1.2.3", TypeSemver) => "semver(ver).compareTo(semver(\"1.2.3\")) > 0".
func BuildExpr(key string, op ComparisonOperator, value interface{}, vt ValueType) (string, error) {
	if value == nil {
		return "", fmt.Errorf("value is nil")
	}

	var expr string
	var err error
	switch vt {
	case TypeBool:
		expr, err = buildBoolExpr(key, op, value)
	case TypeInt:
		expr, err = buildIntExpr(key, op, value)
	case TypeString:
		expr, err = buildStringExpr(key, op, value)
	case TypeSemver:
		expr, err = buildSemverExpr(key, op, value)
	case TypeQuantity:
		expr, err = buildQuantityExpr(key, op, value)
	}
	if err != nil {
		return "", err
	}
	if expr == "" {
		return "", fmt.Errorf("invalid value type %q", vt)
	}
	return expr, ValidateExpr(expr)
}

func ValidateExpr(expression string) error {
	compiler := dracel.GetCompiler()
	result := compiler.CompileCELExpression(expression, dracel.Options{DisableCostEstimation: true})

	if result.Error != nil {
		return result.Error
	}
	return nil
}
