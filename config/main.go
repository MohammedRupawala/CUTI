package config

var Host string = "0.0.0.0"
var Port int = 7379
var MaxKeys = 100

var EvictionStrategy = "all-key-random"

var EvictionFactor = 0.2
