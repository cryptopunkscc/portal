/**
 * Formats an array of query parameters into a single query string.
 *
 * @param {Array} params - An array of query parameter objects and primitive type to format.
 * @return {string} A formatted query string by joining all formatted parameters with '&'.
 * @throws {TypeError} If the input is not an array.
 */
export function formatQueryParams(params) {
  if (!Array.isArray(params)) throw new TypeError('Expected an array of parameters.');
  return params.map(formatQueryParam).join('&')
}

function formatQueryParam(param) {
  if (param === null) return `_=null`
  if (param === undefined) return `_=undefined`
  if (!param) return `_=${encodeURIComponent(param)}`
  if (Array.isArray(param)) throw new TypeError('Expected a non-array.');
  if (typeof param === 'object') return Object.entries(param).map(e =>
    e.map(encodeURIComponent).join('=')
  ).join('&')

  return `_=${encodeURIComponent(param)}`
}

/**
 * Parses a query string into an object where keys map to their corresponding values.
 *
 * @param {string} query - The query string to be parsed. It should be in the format of key=value pairs separated by '&'.
 * @return {Object} - An object representing the parsed query parameters. Keys are strings, and values are strings or arrays of strings for repeated keys.
 */
export function parseQueryParams(query) {
  if (typeof query !== 'string') throw new TypeError('Expected a string.');
  let acc = {};
  query.split('&').map(parseQueryParam).forEach(([key, value]) => {
    if (key in acc) {
      acc[key] = Array.isArray(acc[key]) ? acc[key].concat(value) : [acc[key], value];
    } else if (key === '_') {
      acc[key] = [value];
    } else {
      acc[key] = value;
    }
  })
  return acc
}

function parseQueryParam(param) {
  if (typeof param !== 'string') throw new TypeError('Expected a string.');
  let [key, value] = param.split('=');
  key = decodeURIComponent(key)
  value = decodeURIComponent(value)
  value = parseToPrimitive(value)
  return [key, value];
}

function parseToPrimitive(value) {
  if (value === null || value === undefined) return value;
  if (value === "") return value;

  const num = Number(value);
  if (!isNaN(num)) return num;

  const lower = value.toLowerCase();
  if (lower === "true") return true;
  if (lower === "false") return false;
  if (lower === "null") return null;
  if (lower === "undefined") return undefined;

  return value;

}