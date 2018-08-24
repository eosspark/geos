package common

const BlockIntervalMs int64 = 500
const BlockIntervalUs int64 = 1000 * BlockIntervalMs
const BlockTimestampEpochMs uint64 = 946684800000
const BlockTimestamoEpochNanos int64 = 946684800000000000

/**
 *  The number of sequential blocks produced by a single producer
 */
const ProducerRepetitions int = 12
const MaxProducers int = 125
