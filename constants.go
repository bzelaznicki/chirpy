package main

import "time"

const maxChirpLength = 140
const defaultTokenExpiration = time.Hour
const refreshTokenExpirationDays = 60
const refreshTokenExpiration = refreshTokenExpirationDays * 24 * time.Hour
