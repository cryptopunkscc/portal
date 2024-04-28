<script>
  import Apps from "./Apps.svelte";
  import {onDestroy} from "svelte";
  import {apphost} from "portal/portal.js";
  onDestroy(apphost.interrupt)
</script>

<main>
    <Apps/>
</main>

<style>
    main {
        background: #1b2636;
    }
</style>
