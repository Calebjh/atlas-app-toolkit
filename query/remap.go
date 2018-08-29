package query

import "strings"

// RemapCollectionArgs takes a Filtering and Sorting object and alters the internal
// paths, replacing any occurrences of the vals map keys with their corresponding
// values. This could be alterred to a list of string pairs to potentially improve
// performance, but would be slightly less simple to create
func RemapCollectionArgs(f *Filtering, s *Sorting, vals map[string]string) {
	if f != nil {
		if f.GetOperator() != nil {
			remapFilterOperator(f.GetOperator(), vals)
		}
		var path *[]string
		if cond, ok := f.GetRoot().(*Filtering_StringCondition); ok {
			path = &cond.StringCondition.FieldPath
		} else if cond, ok := f.GetRoot().(*Filtering_NumberCondition); ok {
			path = &cond.NumberCondition.FieldPath
		} else if cond, ok := f.GetRoot().(*Filtering_NullCondition); ok {
			path = &cond.NullCondition.FieldPath
		}
		if path != nil {
			fpath := strings.Join(*path, ".")
			for k, v := range vals {
				fpath = strings.Replace(fpath, k, v, -1)
				*path = strings.Split(fpath, ".")
			}
		}
	}
	if s != nil {
		for _, c := range s.Criterias {
			for k, v := range vals {
				tag := strings.Replace(c.Tag, k, v, -1)
				c.Tag = tag
			}
		}
	}
}

func remapFilterOperator(op *LogicalOperator, vals map[string]string) {
	if ppath := getFieldPathLeft(op); ppath != nil {
		fpath := strings.Join(*ppath, ".")
		for k, v := range vals {
			fpath = strings.Replace(fpath, k, v, -1)
			*ppath = strings.Split(fpath, ".")
		}
	} else if op := op.GetLeftOperator(); op != nil {
		remapFilterOperator(op, vals)
	}
	if ppath := getFieldPathRight(op); ppath != nil {
		fpath := strings.Join(*ppath, ".")
		for k, v := range vals {
			fpath = strings.Replace(fpath, k, v, -1)
			*ppath = strings.Split(fpath, ".")
		}
	} else if op := op.GetRightOperator(); op != nil {
		remapFilterOperator(op, vals)
	}
}

func getFieldPathLeft(op *LogicalOperator) *[]string {
	switch op.GetLeft().(type) {
	case *LogicalOperator_LeftStringCondition:
		return &op.GetLeft().(*LogicalOperator_LeftStringCondition).LeftStringCondition.FieldPath
	case *LogicalOperator_LeftNumberCondition:
		return &op.GetLeft().(*LogicalOperator_LeftNumberCondition).LeftNumberCondition.FieldPath
	case *LogicalOperator_LeftNullCondition:
		return &op.GetLeft().(*LogicalOperator_LeftNullCondition).LeftNullCondition.FieldPath
	}
	return nil
}
func getFieldPathRight(op *LogicalOperator) *[]string {
	switch op.GetRight().(type) {
	case *LogicalOperator_RightStringCondition:
		return &op.GetRight().(*LogicalOperator_RightStringCondition).RightStringCondition.FieldPath
	case *LogicalOperator_RightNumberCondition:
		return &op.GetRight().(*LogicalOperator_RightNumberCondition).RightNumberCondition.FieldPath
	case *LogicalOperator_RightNullCondition:
		return &op.GetRight().(*LogicalOperator_RightNullCondition).RightNullCondition.FieldPath
	}
	return nil
}
