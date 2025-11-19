import { apiClient } from "./http";
import type {
  Rating,
  RatingsSortBy,
  RatingsSortOrder
} from "../types/ratings";



export function getRatings(params: {
  sortBy?: RatingsSortBy
  sortOrder?: RatingsSortOrder
  offset?: number
  limit?: number
  tickerLike?: string
  companyLike?: string
}) {
  return apiClient<Rating[]>(`/v1/stock_ratings`, {
    method: 'GET',
    body: JSON.stringify(params),
  });
}


// export function createUser(payload: CreateUserInput) {
//   return http<User>(`/users`, {
//     method: 'POST',
//     body: JSON.stringify(payload),
//   });
// }
