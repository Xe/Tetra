// Package modes contains a bunch of constants and lookup tables that are
// a pain to use without this file.
package modes

// Converted from a python file

const (
	CHFL_PEON   = 0x0000 // No channel status
	CHFL_VOICE  = 0x0001 // Voiced
	CHFL_HALFOP = 0x0002 // Channel halfop
	CHFL_CHANOP = 0x0004 // Channel operator
	CHFL_ADMIN  = 0x0008 // Channel admin
	CHFL_OWNER  = 0x0010 // Channel owner

	//Channel properties
	PROP_NONE       = 0x00000000 // No properties
	PROP_MUTE       = 0x00000001 // Old +m, mute
	PROP_PRIVATE    = 0x00000002 // Old +p, private channel
	PROP_INVITE     = 0x00000004 // Old +i, invite only
	PROP_TOPICREST  = 0x00000008 // Old +t, only ops can set topic
	PROP_INTERNAL   = 0x00000010 // Old +n, only users in channel can send to it
	PROP_SECRET     = 0x00000020 // Old +s, only users in channel know it exists
	PROP_NOCTCP     = 0x00000040 // Old +C, no CTCP messages
	PROP_NOACTION   = 0x00000080 // Old +D, no CTCP ACTION messages
	PROP_NOKICKS    = 0x00000100 // Old +E, operators cannot kick
	PROP_NOCAPS     = 0x00000200 // OLD +G, ALL CAPITAL LETTER MESSAGES ARE BLOCKED
	PROP_NOREJOIN   = 0x00000400 // Old +J, no immediate rejoin after KICK
	PROP_LARGELIST  = 0x00000800 // Old +L, larger channel lists
	PROP_NOOPERKICK = 0x00001000 // Old +M, staff cannot be kicked
	PROP_OPERONLY   = 0x00002000 // Old +O, only opers may join
	PROP_PERMANENT  = 0x00004000 // Old +P, channel persists without users
	PROP_DISFORWARD = 0x00008000 // Old +Q, channel may not be forwarded to
	PROP_NONOTICE   = 0x00010000 // Old +T, channel may not be NOTICE'd to
	PROP_NOCOLOR    = 0x00020000 // Old +c, channel color codes are stripped
	PROP_NONICKS    = 0x00040000 // Old +d, nick changes are forbidden when in channel
	PROP_FREEINVITE = 0x00080000 // Old +g, invite is freely usable
	PROP_HIDEBANS   = 0x00100000 // Old +u, ban list is hidden without the proper STATUS
	PROP_OPMOD      = 0x00200000 // Old +z, channel messages blocked by something are sent to ops
	PROP_FREEFWD    = 0x00400000 // Old +F, free forwarding
	PROP_NOREPEAT   = 0x00800000 // Old +K, no repeating messages

	//User properties
	UPROP_NONE       = 0x000000
	UPROP_INVISIBLE  = 0x000001 // Old +i, invisible client
	UPROP_CALLERID   = 0x000002 // Old +g, "caller id"
	UPROP_IRCOP      = 0x000004 // Old +o, user is an IRC operator
	UPROP_CLOAKED    = 0x000008 // Old +x, user has a cloaked IP address
	UPROP_ADMIN      = 0x000010 // Old +a, user is an IRC administrator
	UPROP_OVERRIDE   = 0x000020 // Old +p, implicit chanop access
	UPROP_NOCTCP     = 0x000040 // Old +C, prevents receiving CTCP messages other than ACTION (/me)
	UPROP_DEAF       = 0x000080 // Old +D, ignoes all channel messages
	UPROP_DISFORWARD = 0x000100 // Old +Q, prevents channel forwarding
	UPROP_REGPM      = 0x000200 // Old +R, requires people to be registered with services to pm
	UPROP_SOFTCALL   = 0x000400 // Old +G, Soft caller ID, caller id exempting common channels
	UPROP_NOINVITE   = 0x000800 // Old +V, prevents user from getting invites
	UPROP_NOSTALK    = 0x001000 // Old +I, doesn't show channel list in whois
	UPROP_SSLCLIENT  = 0x002000 // Old +Z, client is connected over SSL

	//Channel lists
	LIST_BAN    = 0x0001
	LIST_QUIET  = 0x0002
	LIST_EXCEPT = 0x0004
	LIST_INVEX  = 0x0008
)

// This is a handy lookup table from channel mode letters to bitmasks.
var CHANMODES = []map[string]int{
	map[string]int{
		"q": LIST_QUIET,
		"b": LIST_BAN,
		"e": LIST_EXCEPT,
		"I": LIST_INVEX,
	},
	map[string]int{
		"C": PROP_NOCTCP,
		"D": PROP_NOACTION,
		"E": PROP_NOKICKS,
		"G": PROP_NOCAPS,
		"J": PROP_NOREJOIN,
		"F": PROP_FREEFWD,
		"K": PROP_NOREPEAT,
		"L": PROP_LARGELIST,
		"M": PROP_NOOPERKICK,
		"O": PROP_OPERONLY,
		"P": PROP_PERMANENT,
		"Q": PROP_DISFORWARD,
		"T": PROP_NONOTICE,
		"c": PROP_NOCOLOR,
		"d": PROP_NONICKS,
		"g": PROP_FREEINVITE,
		"i": PROP_INVITE,
		"m": PROP_MUTE,
		"n": PROP_INTERNAL,
		"p": PROP_PRIVATE,
		"s": PROP_SECRET,
		"t": PROP_TOPICREST,
		"u": PROP_HIDEBANS,
		"z": PROP_OPMOD,
	},
	map[string]int{
		"h": CHFL_HALFOP,
		"o": CHFL_CHANOP,
		"v": CHFL_VOICE,
		"a": CHFL_ADMIN,
		"y": CHFL_OWNER,
	},
}

// This is a handy lookup table from user mode flags to bitmasks.
var UMODES = map[string]int{
	"i": UPROP_INVISIBLE,
	"g": UPROP_CALLERID,
	"o": UPROP_IRCOP,
	"a": UPROP_ADMIN,
	"p": UPROP_OVERRIDE,
	"C": UPROP_NOCTCP,
	"D": UPROP_DEAF,
	"Q": UPROP_DISFORWARD,
	"R": UPROP_REGPM,
	"G": UPROP_SOFTCALL,
	"V": UPROP_NOINVITE,
	"I": UPROP_NOSTALK,
	"Z": UPROP_SSLCLIENT,
}

// This is a lookup table for channel prefixes to bitmask flags.
var PREFIXES = map[string]int{
	"+": CHFL_VOICE,
	"%": CHFL_HALFOP,
	"@": CHFL_CHANOP,
	"!": CHFL_ADMIN,
	"~": CHFL_OWNER,
}
