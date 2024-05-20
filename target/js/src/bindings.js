export const bindings = {}

export function inject(platform, adapter) {
  if (platform !== undefined) {
    Object.assign(bindings, {
      platform: platform,
      ...adapter()
    })
  }
}

export default bindings