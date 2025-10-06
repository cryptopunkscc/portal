<script>
  import {rpc} from 'portal'

  const client = rpc.target("portald").bind({user: ["claim"]})

  let nodeId = ''
  let error = false
  let isLoading = false
  let showDialog = false
  let dialogMessage = ''
  let dialogTitle = ''
  let isSuccess = false

  const validateNodeId = () => !(error = nodeId.length < 2)

  const onNodeIdInput = (event) => {
    nodeId = event.target.value
    if (error) validateNodeId()
  }

  const claimUser = async () => {
    if (!validateNodeId()) return
    
    isLoading = true
    try {
      await client.claim(nodeId)
      dialogTitle = 'Success'
      dialogMessage = `Successfully claimed node: ${nodeId}`
      isSuccess = true
      showDialog = true
      nodeId = ''
      error = false
    } catch (err) {
      dialogTitle = 'Error'
      dialogMessage = err.message || 'Failed to claim node. Please try again.'
      isSuccess = false
      showDialog = true
    } finally {
      isLoading = false
    }
  }

  const closeDialog = () => {
    showDialog = false
  }
</script>

<main>
    <div class="claim-container">
        <h1 class="title">Claim new node</h1>
        
        <div class="input-group">
            <input
                    type="text"
                    class="node-input"
                    class:error={error}
                    bind:value={nodeId}
                    on:input={onNodeIdInput}
                    placeholder="node alias or identity"
                    disabled={isLoading}
            />
            {#if error}
                <span class="error-message">Node alias/id requires at least 2 characters.</span>
            {/if}
        </div>

        <button 
            class="primary-button" 
            on:click={claimUser}
            disabled={isLoading}
        >
            {isLoading ? 'Claiming...' : 'Claim'}
        </button>
    </div>

    {#if showDialog}
        <div class="dialog-overlay" on:click={closeDialog}>
            <div class="dialog" on:click|stopPropagation>
                <div class="dialog-header" class:success={isSuccess} class:error={!isSuccess}>
                    <h2>{dialogTitle}</h2>
                </div>
                <div class="dialog-content">
                    <p>{dialogMessage}</p>
                </div>
                <div class="dialog-actions">
                    <button class="dialog-button" on:click={closeDialog}>
                        Close
                    </button>
                </div>
            </div>
        </div>
    {/if}
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

    .claim-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 32px;
        max-width: 400px;
        width: 100%;
    }

    .title {
        font-size: 24px;
        font-weight: 500;
        color: #ccc;
        margin-bottom: 24px;
    }

    .input-group {
        width: 100%;
        margin-bottom: 8px;
    }

    .node-input {
        width: 100%;
        padding: 12px 16px;
        font-size: 16px;
        border: 1px solid #ccc;
        border-radius: 4px;
        outline: none;
        transition: border-color 0.2s;
        box-sizing: border-box;
    }

    .node-input:focus {
        border-color: #1976d2;
    }

    .node-input.error {
        border-color: #d32f2f;
    }

    .node-input:disabled {
        background-color: #f5f5f5;
        cursor: not-allowed;
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

    .primary-button:hover:not(:disabled) {
        background-color: #1565c0;
    }

    .primary-button:active:not(:disabled) {
        background-color: #0d47a1;
    }

    .primary-button:disabled {
        background-color: #90caf9;
        cursor: not-allowed;
    }

    .dialog-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }

    .dialog {
        background-color: white;
        border-radius: 8px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
        min-width: 320px;
        max-width: 400px;
        overflow: hidden;
    }

    .dialog-header {
        padding: 16px 24px;
        border-bottom: 1px solid #e0e0e0;
    }

    .dialog-header.success {
        background-color: #e8f5e9;
        color: #2e7d32;
    }

    .dialog-header.error {
        background-color: #ffebee;
        color: #c62828;
    }

    .dialog-header h2 {
        margin: 0;
        font-size: 20px;
        font-weight: 500;
    }

    .dialog-content {
        padding: 24px;
    }

    .dialog-content p {
        margin: 0;
        font-size: 14px;
        line-height: 1.5;
        color: #555;
    }

    .dialog-actions {
        padding: 16px 24px;
        display: flex;
        justify-content: flex-end;
        border-top: 1px solid #e0e0e0;
    }

    .dialog-button {
        padding: 8px 24px;
        font-size: 14px;
        font-weight: 500;
        color: #1976d2;
        background-color: transparent;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .dialog-button:hover {
        background-color: rgba(25, 118, 210, 0.08);
    }

    .dialog-button:active {
        background-color: rgba(25, 118, 210, 0.16);
    }
</style>
