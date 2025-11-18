import { defineStore } from 'pinia'
import { ref } from 'vue'
import { RatingSearchBy, RatingsSortBy, RatingsSortOrder } from '../types/ratings'






export const useCounterStore = defineStore('counter', () => {

  const searchString = ref('')
  const searchBy = ref<RatingSearchBy>('company')
  const sortBy = ref<RatingsSortBy>(undefined)
  const sortOrder = ref<RatingsSortOrder>('desc')

  return {
    searchString,
    searchBy,
    sortBy,
    sortOrder
  }
})
