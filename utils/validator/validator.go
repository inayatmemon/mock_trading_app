package validator

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// initializing validator and defining custom validators
func Init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.CustomTypeTagMap.Set("objectid", func(field, object interface{}) bool {
		switch v := field.(type) {
		case string:
			_, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return false
			} else {
				return true
			}
		default:
			return false
		}
	})

	govalidator.CustomTypeTagMap.Set("aadhar", func(field, object interface{}) bool {
		switch v := field.(type) {
		case string:
			if aadharValidator(v) {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	})

	govalidator.CustomTypeTagMap.Set("pan", func(field, object interface{}) bool {
		switch v := field.(type) {
		case string:
			if panValidator(v) {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	})

	govalidator.CustomTypeTagMap.Set("mobile", func(field, object interface{}) bool {
		switch v := field.(type) {
		case string:
			if mobileValidator(v) {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	})

	govalidator.CustomTypeTagMap.Set("ordertype", func(field, object interface{}) bool {
		switch v := field.(type) {
		case string:
			if strings.ToLower(v) == "buy" || strings.ToLower(v) == "sell" {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	})

	govalidator.CustomTypeTagMap.Set("ordercategory", func(field, object interface{}) bool {
		switch v := field.(type) {
		case string:
			if strings.ToLower(v) == "market" || strings.ToLower(v) == "stop" || strings.ToLower(v) == "limit" {
				return true
			} else {
				return false
			}
		default:
			return false
		}
	})
}

// PAN number regex pattern
const panPattern = `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`

// Mobile number regex pattern
const mobilePattern = `^[6-9]\d{9}$`

// Aadhar number regex pattern
const aadharPattern = `^[2-9][0-9]{11}$`

// Custom PAN validator function
func panValidator(panNumber string) bool {
	re := regexp.MustCompile(panPattern)
	return re.MatchString(panNumber)
}

// Custom Mobile validator function
func mobileValidator(mobile string) bool {
	re := regexp.MustCompile(mobilePattern)
	return re.MatchString(mobile)
}

// Custom Aadhar validator function
func aadharValidator(aadharNumber string) bool {
	re := regexp.MustCompile(aadharPattern)
	return re.MatchString(aadharNumber)
}
