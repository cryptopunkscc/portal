<script>
  import {log, rpc} from 'portal'
  import {onDestroy, onMount} from "svelte";
  import AppItem from "./AppItem.svelte";
  import readableSet from "./readableSet.js";

  const apps = readableSet((app) => app.package)
  const client = rpc.target("portald").bind(
    "open",
    {
      "app": [
        "install",
        "uninstall",
        "observe",
      ]
    }
  )
  onMount(async () => client.observe.filter((app) => {
    log(JSON.stringify(app))
    app.onClick = () => {
      if (app.installed) client.open(app.package)
      else client.install(app.package)
    }
    apps.set(app)
  }).collect({scope: "gui|srv"}))
  onDestroy(rpc.interrupt)
</script>

<main>
    <div class="all-apps">
        {#each $apps as app, index}
            <AppItem app={app}/>
        {/each}
    </div>
</main>


<style>
    .all-apps {
        margin-outside: 200px;
        /*margin: -8px;*/
        -webkit-user-select: none; /* Safari */
        -ms-user-select: none; /* IE 10 and IE 11 */
        user-select: none; /* Standard syntax */
    }
</style>
