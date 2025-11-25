import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  RatingColumn,
  RatingSearchBy,
  RatingsSortBy,
  RatingsSortOrder
} from '../types/ratings'


export const useRatingsStore = defineStore('ratings', () => {

  const searchString = ref('')
  const searchBy = ref<RatingSearchBy>('company')
  const sortBy = ref<RatingsSortBy>(undefined)
  const sortOrder = ref<RatingsSortOrder>('desc')
  const visibleColumns = ref<Partial<Record<RatingColumn, boolean>>>({'ticker': true})

  const getApiParams = computed(() => {
    const params: Record<string, string> = {}

    // Sorting
    if (sortBy.value !== undefined) {
      params.sort_by = String(sortBy.value)
    }
    if (sortOrder.value !== undefined) {
      params.sort_order = String(sortOrder.value)
    }

    // Searching
    if (searchBy.value) {
      params[searchBy.value] = searchString.value
    }

    // Pagination
    params.limit = String(10)

    return new URLSearchParams(params).toString()
  })

  return {
    searchString,
    searchBy,
    sortBy,
    sortOrder,
    visibleColumns,
    getApiParams
  }
})
