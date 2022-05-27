package constants

const (
	LINKSUFFIX_REGEX     = "^[a-zA-Z0-9-]{1,110}$"
	LINKSUFFIX_REGEX_ERR = "%s must have between 1 and 110 characters, only containing letters, numbers, and dashes."

	DEFAULT_PAYMENT_METHOD_NUMBER_REGEX     = "^[0-9]{15,16}$"
	DEFAULT_PAYMENT_METHOD_NUMBER_REGEX_ERR = "Number must have either 15 or 16 digits."

	DEFAULT_PAYMENT_METHOD_EXP_MONTH_REGEX     = "^[0-9]{1,2}$"
	DEFAULT_PAYMENT_METHOD_EXP_MONTH_REGEX_ERR = "Please enter a valid expiration month. For example, 01 or 12"

	DEFAULT_PAYMENT_METHOD_EXP_YEAR_REGEX     = "^[0-9]{4}$"
	DEFAULT_PAYMENT_METHOD_EXP_YEAR_REGEX_ERR = "Expiration year must have 4 digits."

	DEFAULT_PAYMENT_METHOD_CVC_REGEX     = "^[0-9]+$"
	DEFAULT_PAYMENT_METHOD_CVC_REGEX_ERR = "CVC must have only digits."
)
