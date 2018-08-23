package common

const BlockIntervalMs int = 500
const BlockIntervalUs int = 1000 * BlockIntervalMs
const BlockTimestampEpochMs uint64 = 946684800000

/**
 *  The number of sequential blocks produced by a single producer
 */
const ProducerRepetitions int = 12
const MaxProducers int = 125
