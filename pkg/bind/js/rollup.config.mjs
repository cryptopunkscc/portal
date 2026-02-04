import generatePackageJson from "rollup-plugin-generate-package-json";

export default {
  input: "all.js",
  output: {
    file: "embed/portal/index.js",
    format: "es",
  },
  plugins: [generatePackageJson(
    {
      baseContents: (pkg) => {
        pkg["scripts"] = undefined
        pkg["devDependencies"] = {}
        pkg["module"] = "index.js"
        return pkg
      }
    }
  )]
}
