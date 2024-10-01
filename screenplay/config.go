package screenplay

import "time"

// DefaultPolling is the default polling interval that is used when polling data.
const DefaultPolling = 500 * time.Millisecond

// DefaultTimeout is the default timeout for actions that are waiting on something.
const DefaultTimeout = 2 * time.Second
