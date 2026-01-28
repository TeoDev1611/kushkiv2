<script>
  export let htmlContent = "";

  let iframeElement;

  // Reactive statement to inject styles whenever content updates
  $: if (iframeElement && htmlContent) {
    // Wait for iframe to load new content (simple approach)
    setTimeout(() => {
        try {
            const doc = iframeElement.contentDocument || iframeElement.contentWindow.document;
            if (doc) {
                const style = doc.createElement('style');
                style.textContent = `
                    body { 
                        background-color: transparent !important; 
                        margin: 0; 
                        padding: 0;
                        overflow: hidden;
                        font-family: inherit;
                    }
                    /* Force ECharts container to fit and be clean */
                    .go-echarts { 
                        width: 100% !important; 
                        height: 100% !important; 
                        border: none !important;
                    }
                `;
                doc.head.appendChild(style);
            }
        } catch (e) {
            console.warn("Could not inject chart styles", e);
        }
    }, 100);
  }
</script>

<div class="chart-wrapper">
  {#if htmlContent}
    <iframe
      bind:this={iframeElement}
      srcdoc={htmlContent}
      title="Estadísticas"
      class="chart-frame"
    ></iframe>
  {:else}
    <div class="loading-chart">
        <div class="spinner"></div>
        <span>Cargando datos...</span>
    </div>
  {/if}
</div>

<style>
    .chart-wrapper {
        width: 100%;
        height: 400px;
        border-radius: var(--radius-md);
        overflow: hidden;
        background: var(--bg-surface); /* Matches card background */
        border: 1px solid var(--border-subtle);
        position: relative;
        transition: border-color 0.3s;
    }
    
    .chart-wrapper:hover {
        border-color: var(--border-medium);
    }

    .chart-frame {
        width: 100%;
        height: 100%;
        border: none;
        background: transparent;
        display: block;
    }

    .loading-chart {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: var(--text-secondary);
        gap: 12px;
    }

    .spinner {
        width: 24px;
        height: 24px;
        border: 2px solid var(--border-subtle);
        border-top-color: var(--accent-mint);
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin { to { transform: rotate(360deg); } }
</style>