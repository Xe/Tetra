package atheme

// Fault constants from atheme/doc/XMLRPC
const (
	FAULT_NEEDMOREPARAMS = iota // Not enough parameters
	FAULT_BADPARAMS             // Parameters invalid somehow
	FAULT_NOSUCH_SOURCE         // Source account does not exist
	FAULT_NOSUCH_TARGET         // Target does not exist
	FAULT_AUTHFAIL              // Bad password or authcookie
	FAULT_NOPRIVS               // Permission denied (not auth)
	FAULT_NOSUCH_KEY            // Requested element on target does not exist
	FAULT_ALREADYEXISTS         // Something conflicting already exists
	FAULT_TOOMANY               // Too many of something
	FAULT_EMAILFAIL             // Sending email failed
	FAULT_NOTVERIFIED           // Account not verified
	FAULT_NOCHANGE              // Object is already in requested state
	FAULT_ALREADY_AUTHED        // Already logged in
	FAULT_UNIMPLEMENTED         // Function not implemented
)
