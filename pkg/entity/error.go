package entity

import (
	"errors"
)

/*
 * Add service specific error scenarios here.
 */
var (
	ErrDefault       = errors.New("ErrDefault")       // Default to fallback
	ErrInvalidConfig = errors.New("ErrInvalidConfig") // Invalid configuration in config yml

	// Field errors - Invalid input -- CUSTOM CODE- MODIFY AS REQUIRED
	ErrInvalidInputAttr1 = errors.New("ErrInvalidInputAttr1") // The user input value for InputAttr1 is empty or invalid
	ErrInvalidInputAttr2 = errors.New("ErrInvalidInputAttr2") // The user input value for InputAttr2 is empty or invalid

	// Business errors
	ErrItemNotFound    = errors.New("ErrItemNotFound")    // The key found no records in Databse
	ErrMaxLimitReached = errors.New("ErrMaxLimitReached") // The user has reached the maximum allowed limit of items
	ErrItemExists      = errors.New("ErrItemExists")      // The key to insert already exists in Database

	// Infrastructure errors
	ErrDatabaseFailure          = errors.New("ErrDatabaseFailure")          // Critical Database failure
	ErrInvalidJSON              = errors.New("ErrInvalidJSON")              //Invalid JSON
	ErrInvalidInputResourceType = errors.New("ErrInvalidInputResourceType") // Critical Database failure
	ErrReadyzFailure            = errors.New("ErrReadyzFailure")            // Critical Database failure
	ErrMissingRequiredField     = errors.New("MissingRequiredField")
	ErrHealthzFailure           = errors.New("ErrHealthzFailure")    // Critical Database failure
	ErrInvalidPageNumber        = errors.New("ErrInvalidPageNumber") // Critical Database failure
	ErrInvalidPageSize          = errors.New("ErrInvalidPageSize")   // Critical Database failure
)
