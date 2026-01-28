<script lang="ts">
    import { notifications } from '../../stores/notifications';
    import { fly } from 'svelte/transition';
</script>

<div class="toast-container">
    {#each $notifications as toast (toast.id)}
        <div 
            class="toast {toast.type}" 
            in:fly="{{ y: 20, duration: 300 }}" 
            out:fly="{{ x: 20, duration: 300 }}"
        >
            <div class="toast-content">
                {toast.message}
            </div>
            <button class="toast-close" on:click={() => notifications.update(n => n.filter(t => t.id !== toast.id))}>×</button>
        </div>
    {/each}
</div>

<style>
    .toast-container {
        position: fixed;
        bottom: 24px;
        right: 24px;
        display: flex;
        flex-direction: column;
        gap: 12px;
        z-index: 9999;
        pointer-events: none; /* Allows clicking through the container area */
    }

    .toast {
        pointer-events: auto;
        min-width: 300px;
        max-width: 450px;
        padding: 16px;
        border-radius: 8px;
        background: var(--bg-panel);
        border: 1px solid var(--border-subtle);
        box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        color: var(--text-primary);
        font-size: 14px;
        backdrop-filter: blur(10px);
    }

    .toast.success { border-left: 4px solid var(--status-success); }
    .toast.error { border-left: 4px solid var(--status-error); }
    .toast.warning { border-left: 4px solid var(--status-warning); }
    .toast.info { border-left: 4px solid var(--accent-blue); }

    .toast-close {
        background: none;
        border: none;
        color: var(--text-secondary);
        font-size: 18px;
        cursor: pointer;
        padding: 0 0 0 12px;
        line-height: 1;
    }
    .toast-close:hover { color: var(--text-primary); }
</style>
