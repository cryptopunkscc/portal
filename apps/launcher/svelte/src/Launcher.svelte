<script>
  import Apps from "./Apps.svelte";
  import {onDestroy} from "svelte";
  import {apphost} from "portal";
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
