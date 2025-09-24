package config

// Action strategy constants.
const (
	// ActionStrategyDryRun represents the dry run action strategy.
	ActionStrategyDryRun = "dryRun"
	// ActionStrategyCopy represents the copy action strategy for sorting.
	ActionStrategyCopy = "copy"
	// ActionStrategyMoveToTrash represents the move to trash action strategy for deduplication.
	ActionStrategyMoveToTrash = "moveToTrash"
)

// Keep strategy constants.
const (
	// KeepStrategyOldest represents keeping the oldest file.
	KeepStrategyOldest = "keepOldest"
	// KeepStrategyShortestPath represents keeping the file with the shortest path.
	KeepStrategyShortestPath = "keepShortestPath"
)
