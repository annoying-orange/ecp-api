// Code generated by entc, DO NOT EDIT.

package account

const (
	// Label holds the string label denoting the account type in the database.
	Label = "account"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldReferrers holds the string denoting the referrers field in the database.
	FieldReferrers = "referrers"
	// Table holds the table name of the account in the database.
	Table = "accounts"
)

// Columns holds all SQL columns for account fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldAddress,
	FieldCode,
	FieldReferrers,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultAddress holds the default value on creation for the "address" field.
	DefaultAddress string
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(int) error
)
