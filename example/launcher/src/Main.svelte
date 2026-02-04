<script>
  import AppList from './AppList.svelte';
  import Onboarding from './Onboarding.svelte';
  import hasUser from "./hasUser.js";
</script>

<main>
    {#if $hasUser === null}
        <div class="loading-container">
            <p>Loading...</p>
        </div>
    {:else if $hasUser}
        <AppList/>
    {:else}
        <Onboarding on:complete={hasUser.refresh}/>
    {/if}
</main>

<style>
    main {
        width: 100%;
        height: 100vh;
    }

    .loading-container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        width: 100%;
    }

    .loading-container p {
        font-size: 18px;
        color: #666;
    }
</style>
