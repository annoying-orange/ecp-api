// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/annoying-orange/ecp-api/ent/account"
	"github.com/annoying-orange/ecp-api/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	accountFields := schema.Account{}.Fields()
	_ = accountFields
	// accountDescAddress is the schema descriptor for address field.
	accountDescAddress := accountFields[2].Descriptor()
	// account.DefaultAddress holds the default value on creation for the address field.
	account.DefaultAddress = accountDescAddress.Default.(string)
	// accountDescID is the schema descriptor for id field.
	accountDescID := accountFields[0].Descriptor()
	// account.IDValidator is a validator for the "id" field. It is called by the builders before save.
	account.IDValidator = accountDescID.Validators[0].(func(int) error)
}
