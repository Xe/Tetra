use "strings"

export ^ -- Go style exporting

RPL_LOAD2HI = "Oops! You're using that command too often! Wait a while to spread the love some more!"
RPL_UNKNOWN = "Oops! I don't know who that is!"
RPL_SUCCESS = "Your message has been sent!"

[[
[X] service_bind_command(loveserv, &ls_admirer);
[ ] service_bind_command(loveserv, &ls_rose);
[ ] service_bind_command(loveserv, &ls_chocolate);
[ ] service_bind_command(loveserv, &ls_candy);
[X] service_bind_command(loveserv, &ls_hug);
[ ] service_bind_command(loveserv, &ls_kiss);
[X] service_bind_command(loveserv, &ls_lovenote);
[X] service_bind_command(loveserv, &ls_apology);
[ ] service_bind_command(loveserv, &ls_thankyou);
[ ] service_bind_command(loveserv, &ls_chocobo);
]]

-- Rate limiting
-- Rate limiting will be up to 5 valentine messages per IP address per hour

export rates = {}

--- CheckRates takes in a *tetra.Client and returns true or false if the
--  user has "permission" to do the valentine's day message.
--
--  Returns true if the user is allowed and false if they are not allowed.
CheckRates = (user) ->
  if not rates[user.Uid] -- If the user is not enrolled into LoveServ
    rates[user.Uid] = {} -- Initialise their entry.
    return true

  if #rates[user.Uid] > 4
    return false

  true

--- AddDing adds a "ding" to a user, marking a sucessful use of one of the
--  LoveServ commands.
AddDing = (user) ->
  if not rates[user.Uid]
    rates[user.Uid] = {}

  table.insert rates[user.Uid], os.time!

Hook "CRON-HEARTBEAT", ->
  -- Every 5 minutes, scan over everyone's rates and remove old dings.
  now = os.time!

  for uid, userdings in pairs rates
    for i, ding in pairs userdings
      if now - ding > 30 -- 3600 seconds in an hour, TODO fix
        table.remove userdings, i
        print strings.format "Removed ding at %d for %s", ding, uid

BaseMessage = (source, destination, message, anonymous=false) ->
  if CheckRates source
    AddDing source
    -- TODO: send message
    if not tetra.Clients.ByNick[destination\upper!]
      return RPL_UNKNOWN

    target = tetra.Clients.ByNick[destination\upper!]

    if anonymous
      client.Notice target, message
    else
      client.Notice target, source.Nick .. message

    return RPL_SUCCESS
  RPL_LOAD2HI

Command "HUG", (source, destination, args) ->
  BaseMessage source, args[1], " sent you a darling hug! Adorable!"

Command "ADMIRER", (source, destination, args) ->
  BaseMessage source, args[1], "You have a secret admirer!", true

Command "LOVENOTE", (source, destination, args) ->
  BaseMessage source, args[1], " sent you a love note! Awwwwww!"

Command "APOLOGY", (source, destination, args) ->
  BaseMessage source, args[1], " sent you an apology! Forgiveness is key!"

Command "FORGIVE", (source, destination, args) ->
  BaseMessage source, args[1], " forgave you! Be sure to thank them!"
