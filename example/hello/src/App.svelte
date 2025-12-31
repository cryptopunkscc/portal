<script>
    import { onMount } from 'svelte';
    import { opUser, opSwarmStatus } from './user.js';

    let loading = true;
    let hasUser = false;
    let userAlias = '';
    let nodeAlias = '';
    let error = null;
    let swarmMembers = [];

    onMount(async () => {
        await refresh();
    });

    async function refresh() {
        loading = true;
        error = null;
        try {
            const [userResult, swarmResult] = await Promise.all([
                opUser(),
                opSwarmStatus()
            ]);

            hasUser = userResult.hasUser;
            userAlias = userResult.userAlias;
            nodeAlias = userResult.nodeAlias;
            swarmMembers = swarmResult.members;
        } catch (e) {
            error = 'Failed to connect to Astral daemon';
        } finally {
            loading = false;
        }
    }

    function formatExpiry(dateStr) {
        if (!dateStr) return 'Unknown';
        const date = new Date(dateStr);
        return date.toLocaleDateString();
    }
</script>

<main>
    <div class="card">
        {#if loading}
            <div class="loading">
                <div class="spinner"></div>
                <p>Connecting to Astral...</p>
            </div>
        {:else if error}
            <div class="error">
                <h2>⚠️ Connection Error</h2>
                <p>{error}</p>
                <button on:click={refresh}>Retry</button>
            </div>
        {:else if hasUser}
            <div class="user-info">
                <div class="avatar">👤</div>
                <h1>Welcome back!</h1>
                <div class="details">
                    <div class="detail-item">
                        <span class="label">User</span>
                        <span class="value">{userAlias}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">Node</span>
                        <span class="value">{nodeAlias}</span>
                    </div>
                </div>

                {#if swarmMembers.length > 0}
                    <div class="swarm-section">
                        <h2>🐝 Swarm Members ({swarmMembers.length})</h2>
                        <div class="swarm-list">
                            {#each swarmMembers as member}
                                <div class="member-card">
                                    <div class="member-header">
                                        <span class="member-alias">{member.alias}</span>
                                        <span class="member-status" class:linked={member.linked}>
                                            {member.linked ? '🔗 Linked' : '⭕ Unlinked'}
                                        </span>
                                    </div>
                                    <div class="member-details">
                                        <span class="member-identity" title={member.identity}>
                                            {member.identity}
                                        </span>
                                        <span class="member-expiry">
                                            Expires: {formatExpiry(member.expiresAt)}
                                        </span>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}

                <button on:click={refresh}>Refresh</button>
            </div>
        {:else}
            <div class="no-user">
                <div class="icon">🔒</div>
                <h1>Please log in</h1>
                <p>No user session found</p>
                <button on:click={refresh}>Retry</button>
            </div>
        {/if}
    </div>
</main>

<style>
    main {
        width: 100%;
    }

    .card {
        background: var(--card-bg);
        border-radius: 16px;
        padding: 2.5rem;
        box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
        text-align: center;
    }

    .loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }

    .spinner {
        width: 40px;
        height: 40px;
        border: 3px solid var(--text-muted);
        border-top-color: var(--primary-color);
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }

    .loading p {
        color: var(--text-muted);
    }

    .error {
        color: var(--error-color);
    }

    .error h2 {
        margin-bottom: 0.5rem;
    }

    .error p {
        color: var(--text-muted);
        margin-bottom: 1.5rem;
    }

    .user-info, .no-user {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }

    .avatar, .icon {
        font-size: 3rem;
        margin-bottom: 0.5rem;
    }

    h1 {
        font-size: 1.75rem;
        font-weight: 600;
        margin-bottom: 0.5rem;
    }

    h2 {
        font-size: 1.25rem;
        font-weight: 600;
        margin-bottom: 1rem;
        color: var(--text-color);
    }

    .details {
        width: 100%;
        margin: 1rem 0;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .detail-item {
        display: flex;
        justify-content: space-between;
        padding: 0.75rem 1rem;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 8px;
    }

    .label {
        color: var(--text-muted);
        font-size: 0.875rem;
    }

    .value {
        font-weight: 500;
        color: var(--success-color);
    }

    .no-user p {
        color: var(--text-muted);
        margin-bottom: 1rem;
    }

    .swarm-section {
        width: 100%;
        margin-top: 1.5rem;
        padding-top: 1.5rem;
        border-top: 1px solid rgba(255, 255, 255, 0.1);
    }

    .swarm-list {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .member-card {
        background: rgba(255, 255, 255, 0.05);
        border-radius: 8px;
        padding: 1rem;
        text-align: left;
    }

    .member-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;
    }

    .member-alias {
        font-weight: 600;
        color: var(--text-color);
    }

    .member-status {
        font-size: 0.75rem;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        background: rgba(239, 68, 68, 0.2);
        color: var(--error-color);
    }

    .member-status.linked {
        background: rgba(16, 185, 129, 0.2);
        color: var(--success-color);
    }

    .member-details {
        display: flex;
        justify-content: space-between;
        font-size: 0.75rem;
        color: var(--text-muted);
    }

    .member-identity {
        font-family: monospace;
    }

    button {
        background: var(--primary-color);
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        transition: background 0.2s ease;
        margin-top: 1rem;
    }

    button:hover {
        background: var(--primary-hover);
    }

    button:active {
        transform: scale(0.98);
    }
</style>
