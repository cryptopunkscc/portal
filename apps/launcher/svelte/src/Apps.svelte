<script>
  import App from "./App.svelte";
  import {AppsRepository} from "./apps.js";
  import {onScrollBottomReached} from "./utils.js";
  import {onDestroy} from "svelte";

  const apps = new AppsRepository()
  onScrollBottomReached(() => apps.loadMore())
  onDestroy(() => apps.cancel())
</script>

<div class="all apps">
    {#each $apps as app, index}
        <App app={app}/>
    {/each}
</div>

<style>
    .all.apps {
        margin-outside: 200px;

        -webkit-user-select: none; /* Safari */
        -ms-user-select: none; /* IE 10 and IE 11 */
        user-select: none; /* Standard syntax */
    }
</style>
