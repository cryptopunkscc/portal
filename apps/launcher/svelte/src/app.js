import portal from "./portal.js";

export class AppRepository {
  constructor() {
    this.launch = portal.launch
    this.install = portal.install
    this.uninstall = portal.uninstall
  }
}