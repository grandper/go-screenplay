// Package resolution holds the tools to create your own resolutions. The
// built-in resolutions live in its sub-packages, one per family: be, contain,
// length, match, start and end.
//
// FromFunc creates a resolution from a description and a function with the
// signature of screenplay.Resolution.Inspect, without declaring a dedicated
// type. The type the resolution inspects is inferred from the function:
//
//	func DivisibleBy(divisor int) screenplay.Resolution[int] {
//		return resolution.FromFunc(fmt.Sprintf("be divisible by %d", divisor), func(value int) error {
//			if value%divisor != 0 {
//				return fmt.Errorf("got %d", value)
//			}
//
//			return nil
//		})
//	}
//
//	err := actor.Should(see.The(CashAccount.Balance()).Should(DivisibleBy(100)))
//
// # Writing a resolution
//
// A description is written to follow the verb the resolution is used with:
// "be divisible by 100" after Should, "a valid order id" after ShouldBe,
// "length exactly 2" after ShouldHave.
//
// An error returned by Inspect means that the value does not satisfy the
// resolution, which is what lets a verb such as ShouldNot simply invert the
// outcome.
//
// The one exception is a resolution that was built wrongly: an invalid
// pattern, an impossible bound. Nothing panics when it is built; its Inspect
// returns an error wrapping screenplay.ErrInvalidResolution, and every verb
// then fails, the negative ones included. The built-in resolutions do so, and
// a resolution of your own does the same by wrapping that error.
package resolution
