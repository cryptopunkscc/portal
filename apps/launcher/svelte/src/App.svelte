<script>
  import {AppRepository} from "./app.js";
  import {log} from "../../lib/portal.js";

  const repo = new AppRepository()
  export let app = {
    title: '',
    package: '',
    name: '',
    description: '',
    icon: '',
    url: '',
    installed: false,
  };
  const launch = () => repo.launch(app.package).catch(log)
  const install = () => repo.install(app.package).catch(log)
  const uninstall = () => repo.uninstall(app.package).catch(log)

</script>

<div class="app-item" on:click={launch}>
    <div class="app-icon">
        <img src="{app.icon}" alt="{app.name}"/>
    </div>

    <div class="app-info">
        <div class="app-name">
            <h2>{app.title}</h2>
            <span class="installed-status">{app.installed ? '(Zainstalowano)' : ''}</span>
        </div>
        <p class="app-description">{app.description}</p>
    </div>

<!--    <div class="app-actions">-->
<!--        <button on:click={() => {}}>-->
<!--            {app.installed ? 'Uruchom' : 'Zainstaluj'}-->
<!--        </button>-->
<!--    </div>-->
</div>

<style>
    .app-item {
        display: flex;
        align-items: center;
        padding: 18px 36px;
        /*background-color: #333;*/
        border-radius: 3px;
        border-bottom: 1px solid #222;
        cursor: pointer;
        transition: background-color 0.1s ease-in-out;
    }

    .app-item:hover {
        background-color: #ffffff09;
    }

    .app-icon {
        margin-right: 36px;
    }

    .app-icon img {
        width: 40px;
        height: 40px;
        border-radius: 50%;
    }

    .app-info {
        flex: 1;
    }

    .app-name {
        display: flex;
        align-items: center;
    }

    .app-name h2 {
        margin: 0;
        font-size: 18px;
        font-weight: bold;
        color: #fff;
    }

    .app-name .installed-status {
        margin-left: 12px;
        font-size: 12px;
        color: #999;
    }

    .app-description {
        margin-top: 8px;
        font-size: 16px;
        color: #aaa;
        display: flex;
        /*text-align: center; !* Wyśrodkowanie opisu *!*/
    }

    .app-actions {
        display: flex;
        justify-content: flex-end;
    }

    .app-actions button {
        padding: 8px 15px;
        border: none;
        border-radius: 3px;
        background-color: #007bff;
        color: #fff;
        font-weight: bold;
        cursor: pointer;
        transition: background-color 0.2s ease-in-out;
    }

    .app-actions button:hover {
        background-color: #0062cc;
    }

    /* Styl przycisku w stanie "Zainstalowano" */
    .app-actions button.installed {
        background-color: #4caf50;
    }

    .app-actions button.installed:hover {
        background-color: #388e3c;
    }
</style>
