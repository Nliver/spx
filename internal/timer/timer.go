package timer

var (
	gameTimer float64

	timestamps     []int64
	nextTimerIndex int
)

const TIME_PERCISION = 1000

func Timer() float64 {
	return float64(int64(gameTimer*TIME_PERCISION)) / TIME_PERCISION
}

func ResetTimer() {
	gameTimer = 0
	nextTimerIndex = 0
}

func OnReload() {
	ResetTimer()
	timestamps = timestamps[:0]
	nextTimerIndex = 0
}

func RegisterTimer(timer float64) {
	timeStamp := int64(timer * TIME_PERCISION)

	low, high := 0, len(timestamps)
	for low < high {
		mid := (low + high) / 2
		if timestamps[mid] == timeStamp {
			return
		} else if timestamps[mid] < timeStamp {
			low = mid + 1
		} else {
			high = mid
		}
	}

	timestamps = append(timestamps, 0)
	copy(timestamps[low+1:], timestamps[low:])
	timestamps[low] = timeStamp
}

func CheckTimerEvent() float64 {
	if len(timestamps) == 0 {
		return -1
	}

	if len(timestamps) <= nextTimerIndex {
		return -1
	}
	targetTimer := timestamps[nextTimerIndex]
	if targetTimer > int64(gameTimer*TIME_PERCISION) {
		return -1
	}
	nextTimerIndex++
	return float64(targetTimer) / TIME_PERCISION
}

func OnUpdate(deltaTime float64) {
	gameTimer += deltaTime
}
