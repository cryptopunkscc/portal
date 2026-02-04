<script>
  import {rpc} from 'portal'
  import {onMount} from 'svelte'
  import hasUser from "./hasUser.js";

  const client = rpc.target("portald").bind({"user": ["create"]})

  let user = ''
  let error = false
  let nodeAlias = undefined

  const validateAlias = () => !(error = user.length < 2)

  const onAliasInput = (event) => {
    user = event.target.value
    if (error) validateAlias()
  }

  const createUser = async () => {
    if (!validateAlias()) return
    await client.create(user)
    await hasUser.refresh()
  }

  onMount(async () => {
    const profile = await rpc.target("localnode").call(".profile").request()
    nodeAlias = profile.alias
  })
</script>

<main>
    <div class="onboarding-container">
        <div class="input-group">
            <input
                    type="text"
                    class="alias-input"
                    class:error={error}
                    bind:value={user}
                    on:input={onAliasInput}
                    placeholder="user alias"
            />
            {#if error}
                <span class="error-message">User alias requires at least 2 characters.</span>
            {/if}
        </div>

        <button class="primary-button" on:click={createUser}>
            Create new user
        </button>

        <div class="divider-container">
            <div class="divider"></div>
            <span class="divider-text">or</span>
            <div class="divider"></div>
        </div>

        <p class="claim-info">
            Claim <strong class="node-alias">{nodeAlias}</strong> from another device then...
        </p>

        <button class="primary-button" on:click={hasUser.refresh}>
            Confirm claimed
        </button>
    </div>
</main>

<style>
    main {
        margin: -16px -8px;
        margin-outside: 200px;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        width: 100vw;
    }

    .onboarding-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 32px;
        max-width: 400px;
        width: 100%;
    }

    .input-group {
        width: 100%;
        margin-bottom: 8px;
    }

    .alias-input {
        width: 100%;
        padding: 12px 16px;
        font-size: 16px;
        border: 1px solid #ccc;
        border-radius: 4px;
        outline: none;
        transition: border-color 0.2s;
        box-sizing: border-box;
    }

    .alias-input:focus {
        border-color: #1976d2;
    }

    .alias-input.error {
        border-color: #d32f2f;
    }

    .error-message {
        display: block;
        color: #d32f2f;
        font-size: 12px;
        margin-top: 4px;
        margin-left: 4px;
    }

    .primary-button {
        width: 100%;
        padding: 12px 24px;
        font-size: 16px;
        font-weight: 500;
        color: white;
        background-color: #1976d2;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
        margin-bottom: 8px;
    }

    .primary-button:hover {
        background-color: #1565c0;
    }

    .primary-button:active {
        background-color: #0d47a1;
    }

    .divider-container {
        display: flex;
        align-items: center;
        width: 100%;
        margin: 16px 0;
    }

    .divider {
        flex: 1;
        height: 1px;
        background-color: #ccc;
    }

    .divider-text {
        padding: 0 24px;
        color: #666;
        font-size: 14px;
    }

    .claim-info {
        width: 100%;
        text-align: center;
        font-size: 14px;
        color: #555;
        margin: 8px 0 12px 0;
        line-height: 1.5;
    }

    .node-alias {
        font-weight: 700;
        color: #1976d2;
    }
</style>
