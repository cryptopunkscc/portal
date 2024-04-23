export function isBottomReached(offset) {
  offset = offset || 1
  const scrollY = window.scrollY
  const scrollH = document.documentElement.scrollHeight
  const windowH = window.innerHeight
  return scrollH - scrollY - windowH < offset
}

export function onScrollBottomReached(callback) {
  let was = isBottomReached()
  let is = was
  if (is) {
    callback()
  }
  const listener = () => {
    is = isBottomReached()
    if (!was && is) {
      callback()
    }
    was = is
  }
  document.addEventListener("scroll", listener)
  return () => document.removeEventListener("scroll", listener)
}
