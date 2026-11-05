
package retry

import "time"

func NextDelay(attempts int) time.Duration {
    const baseDelay = 2 * time.Second
    const maxDelay = 30 * time.Second

    delay := baseDelay

    for i := 1; i < attempts; i++ {
        delay *= 2

        if delay >= maxDelay {
            return maxDelay
        }
    }

    return delay
}

