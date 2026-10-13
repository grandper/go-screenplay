package screenplay

import "time"

// DefaultPolling is the default polling interval that is used when polling data.
const DefaultPolling = 500 * time.Millisecond

// DefaultTimeout is the default timeout for actions that are waiting on something.
const DefaultTimeout = 2 * time.Second

// DefaultMaxBufferedSize bounds what one command output, one HTTP body or one
// file read may hold in memory: 64 MiB. Reading more than that fails instead of
// exhausting the memory of the test process.
const DefaultMaxBufferedSize = 64 << 20
