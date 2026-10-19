package choose

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/choose"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// To starts a multi-way choice with the performables of its first branch; close
// the branch with When or WhenThe, and name the next one with To:
//
//	choose.To(ReadTheResponse()).WhenThe(statusCode).HappensTo(be.EqualTo(200)).
//		To(AbortTheRequest()).Otherwise()
func To(actions ...screenplay.Performable) *impl.ToAction {
	return impl.To(actions...)
}
