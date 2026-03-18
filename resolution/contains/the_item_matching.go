package contains

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheItemMatching creates a matcher to tell if an iterable contains an item matching a regular expression.
func TheItemMatching(regex *regexp.Regexp) *TheItemMatchingResolution {
	return &TheItemMatchingResolution{
		regex: regex,
	}
}

// TheItemMatchingRegexString is a matcher to tell if an iterable contains an item matching a regular expression.
func TheItemMatchingRegexString(regexStr string) *TheItemMatchingResolution {
	return &TheItemMatchingResolution{
		regex: regexp.MustCompile(regexStr),
	}
}

// TheItemMatchingResolution is a matcher to tell if an iterable contains an item matching a regular expression.
type TheItemMatchingResolution struct {
	regex *regexp.Regexp
}

// Resolve creates a matcher to make an assertion.
func (r *TheItemMatchingResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		return hasMatchingItem(r.regex, obj)
	}
}

func hasMatchingItem(regex *regexp.Regexp, obj any) (bool, error) {
	if obj == nil {
		return false, nil
	}

	objValue := reflect.ValueOf(obj)

	switch objValue.Kind() { //nolint:exhaustive // we handle only the cases we need
	case reflect.Chan:
		return itemInChanMatchesRegex(regex, objValue)
	case reflect.Slice, reflect.Array:
		return itemInSliceMatchesRegex(regex, objValue)
	case reflect.Map:
		return false, errors.New("the object should be an iterable on strings")
	case reflect.Ptr:
		if objValue.IsNil() {
			return false, nil
		}

		ref := objValue.Elem().Interface()

		return hasMatchingItem(regex, ref)
	default:
		return valueMatchesRegex(regex, objValue)
	}
}

func itemInChanMatchesRegex(regex *regexp.Regexp, objValue reflect.Value) (bool, error) {
	for {
		chanValue, ok := objValue.Recv()
		if !ok {
			return false, nil
		}

		matched, err := valueMatchesRegex(regex, chanValue)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}
}

func itemInSliceMatchesRegex(regex *regexp.Regexp, objValue reflect.Value) (bool, error) {
	for i := range objValue.Len() {
		sliceValue := objValue.Index(i)

		matched, err := valueMatchesRegex(regex, sliceValue)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func valueMatchesRegex(regex *regexp.Regexp, value reflect.Value) (bool, error) {
	switch value.Kind() { //nolint:exhaustive // we handle only the cases we need
	case reflect.String:
		return regex.MatchString(value.String()), nil
	case reflect.Ptr, reflect.UnsafePointer:
		if value.IsNil() {
			return false, errors.New("the object should not be nil")
		}

		ref := value.Elem()

		return valueMatchesRegex(regex, ref)
	default:
		return false, errors.New("the object should be an iterable on strings")
	}
}

// String describe the resolution's expectation.
func (r *TheItemMatchingResolution) String() string {
	return fmt.Sprintf("containing an item which matches pattern %s", r.regex)
}

// TheItemMatchingResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheItemMatchingResolution)(nil)
