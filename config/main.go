package config

var Host string = "0.0.0.0"
var Port int = 7379
var MaxKeys = 100

var EvictionStrategy = "approx-lru"

var EvictionFactor = 0.2

var LRU_BITS = 24
var POOL_SIZE = 16
