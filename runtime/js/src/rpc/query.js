
export function splitQuery(query) {
  const index = query.search(/[?.{\[]/)
  if (index === -1) {
    return [query]
  }
  const left = query.slice(0, index)
  let right = query.slice(index, query.length)
  if (/^[.?]/.test(right)) {
    right = right.slice(1)
  }
  return [left, right]
}


export function hasParams(query) {
  return query.search(/[?{\[]/) > -1
}

