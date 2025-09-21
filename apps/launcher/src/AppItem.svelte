<script>
  export let app = {
    title: '',
    description: '',
    icon: '',
    installed: false,
    onClick: () => undefined
  }
</script>

<div class="app-item" on:click={app.onClick}>
    <div class="app-info">
        <p class="app-name">{app.title}</p>
        <p class="app-description">{app.description}</p>
        <p class="app-status">{app.installed ? '' : 'not installed'}</p>
    </div>
</div>

<style>
    .app-item {
        display: flex;
        align-items: center;
        background-color: #333;
        padding: 18px 36px;
        border-radius: 3px;
        border-bottom: 1px solid #222;
        cursor: pointer;
        transition: background-color 0.1s ease-in-out;
    }

    .app-item:hover {
        background-color: #ffffff09;
    }

    .app-info {
        flex: 1;
    }

    .app-name {
        display: flex;
        align-items: center;
        margin: 0;
        font-size: 18px;
        font-weight: bold;
        color: #fff;
    }

    .app-status {
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
</style>

