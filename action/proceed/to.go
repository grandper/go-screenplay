package proceed

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// To starts an action that performs the given performables. Guard it with a
// condition to perform them only when the condition holds:
//
//	proceed.To(RetryTheRequest()).IfThe(statusCode).HappensTo(be.AtLeast(500))
func To(actions ...screenplay.Performable) *impl.ToAction {
	return impl.To(actions...)
}
