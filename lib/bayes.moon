module "bayes", package.seeall

export ^

-- https://github.com/zach-binary/moon-grahams

Config = {
  GoodTokenWeight: 2
  MinTokenCount: 0
  MinCountForInclusion: 5
  MinScore: 0.001
  MaxScore: 0.9999999
  LikelySpamScore: 0.99
  CertainSpamScore: 0.9999
  CertainSpamCount: 10
  InterestingWordCount: 15
}

class Corpus
  -- Pattern to select words that don't begin with a number
  @TokenPattern = '([a-zA-Z]%w+)%W*'

  new: (line) =>
    @Tokens = {}
    @NumTokens = 0

    if line ~= nil
      @ProcessTextLine line

  ProcessTextLine: (line) =>
    for match in string.gmatch line, @@TokenPattern
      if #match > 2
        @AddToken match\lower!


  AddToken: (rawPhrase) =>
    if (@Tokens[rawPhrase])
      @Tokens[rawPhrase] = @Tokens[rawPhrase] + 1
    else
      @Tokens[rawPhrase] = 1
      @NumTokens = @NumTokens + 1

class Filter
  new: (good, bad) =>
    @Good = good
    @Bad = bad

    @CalculateProbabilities!

  CalculateProbabilities: () =>
    @Probabilities = {}

    for token, score in pairs @Good.Tokens
      @CalculateTokenProbability token

    remainingTokens = {k,v for k, v in pairs @Bad.Tokens when not @Probabilities[k]}

    for token, score in pairs remainingTokens
      @CalculateTokenProbability token


  CalculateTokenProbability: (token) =>
    g = if @Good.Tokens[token] then @Good.Tokens[token] * Config.GoodTokenWeight else 0
    b = if @Bad.Tokens[token] then @Bad.Tokens[token] else 0

    if (g + b > Config.MinCountForInclusion)

      goodFactor = math.min 1, g / @Good.NumTokens
      badFactor = math.min 1, b / @Bad.NumTokens

      prob = math.max Config.MinScore, math.min Config.MaxScore, badFactor / (goodFactor + badFactor)

      if g == 0
        prob = if b > Config.CertainSpamCount then Config.CertainSpamScore else Config.LikelySpamScore

      @Probabilities[token] = prob

  -- Returns probability of spam and the list of interesting words
  Test: (message) =>
    probs = {}
    index = 0

    message = message\lower!

    for token in string.gmatch message, Corpus.TokenPattern
      if @Probabilities[token]

        prob = @Probabilities[token]

        -- here we're storing the 'interestingness' of the word as a key
        key = string.format '%.5f', tostring(0.5 - math.abs (0.5 - prob))
        key ..=  token
        key ..= tostring(index + 1)
        index += 1
        probs[key] = prob

    mult = 1 -- abc..n
    comb = 1 -- (1 - a)(1 - b)..(1 - n)
    index = 0

    -- sort the words of a message by how interesting they are, not probability
    probsSorted = {}
    for Interest, Probability in pairs probs
      table.insert probsSorted, {:Interest, :Probability}

    table.sort probsSorted, (a, b) -> return a.Interest < b.Interest

    words = {}
    for i, prob in ipairs probsSorted

      Probability = prob.Probability

      mult *= Probability
      comb *= (1 - Probability)

      Word = string.match(prob.Interest, Corpus.TokenPattern)

      table.insert words, {:Word, :Probability}

      index += 1

      if index > Config.InterestingWordCount
        break

    return mult / (mult + comb), words

class FileFilter extends Filter
  new: (badfile, goodfile) =>
    with io.open badfile, "r"
      text = \read "*a"
      @Bad = Corpus text

    with io.open goodfile, "r"
      text = \read "*a"
      @Good = Corpus text

    @CalculateProbabilities!
