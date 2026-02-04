import {formatQueryParams} from "./params";

export function formatQuery(port, params) {
  let query = port
  if (params.length > 0) {
    query += '?' + formatQueryParams(params)
  }
  return query
}

export function hasParams(query) {
  return query.search(/[?]/) > -1
}

