export type RatingsSortBy =
  'target_from' |
  'target_to' |
  'target_delta' |
  'score' |
  'ticker' |
  'company' |
  'brokerage' |
  'action' |
  'rating_from' |
  'rating_to' | undefined
export type RatingsSortOrder = 'asc' | 'desc'
export type RatingSearchBy = 'company' | 'ticker'
export type RatingColumn = 'email' | string


export type Rating = {
  ticker: string
  company: string
  target_from: string
  target_to: string
  action: 'up' | 'down' | 'reiterated'
  rating_from: 'buy' | 'hold' | 'sell' | 'pending'
  rating_to: 'buy' | 'hold' | 'sell' | 'pending'
  at: string
  target_delta: string
  score: number
}

