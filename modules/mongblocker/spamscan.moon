require "lib/bayes"

export ^

Scores = {}
Filter = bayes.FileFilter "etc/spamscan/spam.txt", "etc/spamscan/ham.txt"

class Client
  new: (nick, score = 0, warnings = 0) =>
    @nick = nick
    @score = score
    @warnings = warnings

  SetScore: (delta) =>
    @score += delta

    if @score < -10
      @score = -10
    elseif @score > 10
      @score = 10

Command "SCORES", true, (source, destination, msg) ->
  parc = #msg

  if parc ~= 1
    return "Need a channel to look up scores in"

  channel = msg[1]\upper!

  if Scores[channel] == nil
    return "Scores for #{channel} not found."

  print channel
  print Scores[channel]
  print source

  client.Notice source, "Scores for #{channel}"
  client.Notice source, strings.format("%-20s | %-4s | %s", "Nickname", "Warn", "Score")

  print "Got here"

  for uid, scoree in pairs Scores[channel]
    client.Notice source,
      strings.format("%-20s | %v    | %v", scoree.nick, scoree.warnings, scoree.score)

  "End of list"

Hook "MONGBLOCKER-CHANMSG", (source, destination, msg) ->
  dest = destination.Target!
  src = source.Uid

  if strings.first(msg[1]) == tetra.bot.Config.General.Prefix
    return

  if Scores[dest] == nil
    Scores[dest] = {}

  if Scores[dest][src] == nil
    Scores[dest][src] = Client source.Nick

  score = Filter\Test strings.join(msg, " ")
  clientscore = Scores[dest][src]

  if score > .8
    clientscore\SetScore score * 1.2
    --client.Privmsg destination, "I think that is spam (#{clientscore.score}))"
  else
    if score < 0.01
      clientscore\SetScore -1*(score*10)
    elseif score < 0.3
      clientscore\SetScore -1*(score*3)
    else
      clientscore\SetScore -1*(score*1.5)
    --client.Privmsg destination, "That isn't spam (#{clientscore.score})"

  if clientscore.score > 8 and clientscore.warnings >= 0
    switch clientscore.warnings
      when 2
        client.Privmsg destination, "!kick #{source.Nick}"
        clientscore.score = 3
      when 1
        client.Privmsg destination, "#{source.Nick}: please do not spam. If you continue I will kick you."
        clientscore.warnings = 2
      when 0
        client.Privmsg destination, "#{source.Nick}: please say more constructive things."
        clientscore.warnings = 1
