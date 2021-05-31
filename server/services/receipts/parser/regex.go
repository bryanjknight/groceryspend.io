package parser

import (
	"regexp"
)

// collection of regexes
var phoneNumberRegex = regexp.MustCompile(`(\([0-9]{3}\) |[0-9]{3}-)[0-9]{3}-[0-9]{4}`)
var addressRegex = regexp.MustCompile(`(?i)\d+ [A-Z 0-9]+ [A-Z]{2,5}`)
var townCityZipRegex = regexp.MustCompile(`(?i)[A-Z ]+, [A-Z]{2} \d{5}(-\d{4})?`)
var cashierRegex = regexp.MustCompile(`(?i)cashier:? [a-z0-9 ]+`)
var storeRegex = regexp.MustCompile(`(?i)store:? \d{1,5}`)

// TODO: use levenshtein distance to judge whether stop word hit because of misread characters
var departmentNamesRegex = regexp.MustCompile(`(?i)(dairy|produce|froduce|meat|grocery)`)

var taxRegex = regexp.MustCompile(`(?i)(total )?tax`)
var subtotalRegex = regexp.MustCompile(`(?i)subtotal`)
var totalRegex = regexp.MustCompile(`(?i)total|balance`)

var dateRegex = regexp.MustCompile(`(\d{2}/\d{2}/\d{2,4}|\d{1,2}[A-Z]{3}\d{4})`)
var timeRegex = regexp.MustCompile(`\d{1,2}:\d{2}:\d{2}`)

var priceRegexStr = `(-?\d{0,5}\.\d{2})`
var priceRegex = regexp.MustCompile(priceRegexStr)

var weightRegexStr = `(?i)(\d{1,3}\.\d{2}) (lb|1b|oz)`
var weightRegex = regexp.MustCompile(weightRegexStr)
