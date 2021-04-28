package db

import (
	"fmt"
	"strconv"
)

type Stack []bool

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(b bool) {
	*s = append(*s, b) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (bool, bool) {
	if s.IsEmpty() {
		return false, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// postfixNotation converts the Rule to postfix notation inorder to perform stack methodb
func postfixNotation(rules Rule, parent, child rowInfo) (pfn []interface{}) {
	var opr QueryOp
	oprFlag := false
	for _, rule := range rules {
		switch t := rule.(type) {
		case Rule:
			pfn = append(pfn, postfixNotation(t, parent, child)...)
		case RuleFeild:
			left, right := getOperands(t, parent, child)
			pfn = append(pfn, ruleOperations(left, right, t.Operator))
			fmt.Println("pfn-->>", pfn)
			if oprFlag {
				pfn = append(pfn, opr)
				oprFlag = false
			}
		case QueryOp:
			opr = t
			oprFlag = true
			continue
		}
	}
	return
}
func whereClouse(rule Rule, parent, child rowInfo) bool {
	var stack Stack

	fmt.Println(rule...)
	res := postfixNotation(rule, parent, child)
	for _, rl := range res {
		if v, ok := rl.(bool); ok {
			stack.Push(v)
		} else if v, ok := rl.(QueryOp); ok {
			rval, empt := stack.Pop()
			if empt != false {
				return true
			}
			lval, empt := stack.Pop()
			if empt != false {
				return true
			}
			fval := boolOperaiton(lval, rval, v)
			stack.Push(fval)
		}
		fmt.Println(stack)
	}
	if result, ok := stack.Pop(); ok {
		return result
	}

	return false
}
func getOperands(rf RuleFeild, parent, child rowInfo) (left, right interface{}) {
	if rf.LeftRule != nil {
		if rf.LeftRule.TableName == parent.name {
			left = parent.row[rf.LeftRule.FieldName]

		} else if rf.LeftRule.TableName == child.name {
			left = child.row[rf.LeftRule.FieldName]
		}
	} else {
		left = rf.LeftOperand
	}
	if rf.RightRule != nil {
		if rf.RightRule.TableName == parent.name {
			right = parent.row[rf.RightRule.FieldName]
		} else if rf.RightRule.TableName == child.name {
			right = child.row[rf.RightRule.FieldName]
		}
	} else {
		right = rf.RightOperand
	}
	return
}
func ruleOperations(left, right QueryOprnd, operator QueryOp) bool {
	switch operator {
	case opEQ:
		if lVal, rVal, ok := getInt(left, right); ok {
			println(lVal, rVal)
			return lVal == rVal
		}
		if lVal, rVal, ok := getStr(left, right); ok {
			return lVal == rVal
		}
	case OpNTEQ:
		if lVal, rVal, ok := getInt(left, right); ok {
			return lVal != rVal
		}
		if lVal, rVal, ok := getStr(left, right); ok {
			return lVal != rVal
		}
	case OpGTEQ:
		if lVal, rVal, ok := getInt(left, right); ok {
			return lVal >= rVal
		}
		if lVal, rVal, ok := getStr(left, right); ok {
			return lVal >= rVal
		}
	case OpGT:
		if lVal, rVal, ok := getInt(left, right); ok {
			return lVal > rVal
		}
		if lVal, rVal, ok := getStr(left, right); ok {
			return lVal > rVal
		}
	case OpLT:
		if lVal, rVal, ok := getInt(left, right); ok {
			return lVal < rVal
		}
		if lVal, rVal, ok := getStr(left, right); ok {
			return lVal < rVal
		}
	case OpLTEQ:
		if lVal, rVal, ok := getInt(left, right); ok {
			return lVal <= rVal
		}
		if lVal, rVal, ok := getStr(left, right); ok {
			return lVal <= rVal
		}

	}

	return true
}
func boolOperaiton(left, right bool, operator QueryOp) bool {
	switch operator {
	case opEQ:
		return left == right
	case OpNTEQ:
		return left != right
	case OpAND:
		return left && right
	case OpOR:
		return left || right
	}

	return true
}

func getInt(left, right interface{}) (lval, rval int32, ok bool) {
	lval, lOk := left.(int32)
	rval, rOk := right.(int32)
	if lOk && rOk {
		return lval, rval, true
	} else if lOk {
		if rl, ok := right.(string); ok {
			rval, err := strconv.Atoi(rl)
			if err == nil {
				return lval, int32(rval), true
			}
		}
	} else if rOk {
		if ll, ok := left.(string); ok {
			lval, err := strconv.Atoi(ll)
			if err == nil {
				return int32(lval), rval, true
			}
		}
	}
	return lval, rval, false
}
func getStr(left, right interface{}) (lval, rval string, ok bool) {
	lval, lOk := left.(string)
	rval, rOk := right.(string)
	if lOk && rOk {
		return lval, rval, true
	}
	return lval, rval, false
}

// returns if $JOIN rule found
func getJoinInfo(rules []Rules) Rule {
	if rules == nil {
		return nil
	}
	for _, rule := range rules {
		if rule.Label == "$JOIN" {
			return rule.Rule
		}
	}
	return nil
}

// returns if $WHERE rule found
func getWhereInfo(rules []Rules) Rule {
	if rules == nil {
		return nil
	}
	for _, rule := range rules {
		if rule.Label == "$WHERE" {
			return rule.Rule
		}
	}
	return nil
}

// sets the rlation on join rule
func setJoinRelation(rule Rule, par, chld string) (rel relation) {
	if rule == nil {
		return
	}
	for _, rul := range rule {
		if rf, ok := rul.(RuleFeild); ok {
			if par == rf.LeftRule.TableName && chld == rf.RightRule.FieldName || chld == rf.LeftRule.TableName && par == rf.RightRule.FieldName {
				rel.relToPar = &relInfo{rf.RightRule, rf.LeftRule}
				rel.relToChld = nil
				return
			}

			// set parent relation
			if par == rf.LeftRule.TableName {
				rel.relToPar = &relInfo{rf.RightRule, rf.LeftRule}
			} else if par == rf.RightRule.TableName {
				rel.relToPar = &relInfo{rf.LeftRule, rf.RightRule}
			}
			// set child relation
			if chld == rf.LeftRule.TableName {
				rel.relToChld = &relInfo{rf.RightRule, rf.LeftRule}
			} else if chld == rf.RightRule.TableName {
				rel.relToChld = &relInfo{rf.LeftRule, rf.RightRule}
			}
		}
	}
	return
}
