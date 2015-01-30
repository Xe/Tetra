package atheme

// Fault constants from atheme/doc/XMLRPC
const (
	FaultNEEDMOREPARAMS = iota // Not enough parameters
	FaultBADPARAMS             // Parameters invalid somehow
	FaultNOSUCHSOURCE          // Source account does not exist
	FaultNOSUCHTARGET          // Target does not exist
	FaultAUTHFAIL              // Bad password or authcookie
	FaultNOPRIVS               // Permission denied (not auth)
	FaultNOSUCHKEY             // Requested element on target does not exist
	FaultALREADYEXISTS         // Something conflicting already exists
	FaultTOOMANY               // Too many of something
	FaultEMAILFAIL             // Sending email failed
	FaultNOTVERIFIED           // Account not verified
	FaultNOCHANGE              // Object is already in requested state
	FaultALREADYAUTHED         // Already logged in
	FaultUNIMPLEMENTED         // Function not implemented
)
