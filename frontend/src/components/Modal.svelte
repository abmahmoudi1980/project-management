<script>
  import { createEventDispatcher } from "svelte";
  import { fade, scale } from "svelte/transition";

  let { show = false, title = "", maxWidth = "2xl" } = $props();

  const dispatch = createEventDispatcher();

  const maxWidthClasses = {
    sm: "max-w-sm",
    md: "max-w-md",
    lg: "max-w-lg",
    xl: "max-w-xl",
    "2xl": "max-w-2xl",
    "3xl": "max-w-3xl",
    "4xl": "max-w-4xl",
    "5xl": "max-w-5xl",
  };

  function handleBackdropClick(event) {
    if (event.target === event.currentTarget) {
      close();
    }
  }

  function handleKeydown(event) {
    if (event.key === "Escape") {
      close();
    }
  }

  function close() {
    dispatch("close");
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50"
    on:click={handleBackdropClick}
    transition:fade={{ duration: 200 }}
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white rounded-lg shadow-xl w-full {maxWidthClasses[
        maxWidth
      ]} max-h-[90vh] flex flex-col"
      transition:scale={{ duration: 200, start: 0.95 }}
    >
      <!-- Modal Header -->
      <div
        class="flex items-center justify-between px-6 py-4 border-b border-slate-200"
      >
        <h2 class="text-xl font-semibold text-slate-900">{title}</h2>
        <button
          on:click={close}
          class="text-slate-400 hover:text-slate-600 transition-colors"
          aria-label="Close modal"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <!-- Modal Body (scrollable) -->
      <div class="flex-1 overflow-y-auto px-6 py-6">
        <slot />
      </div>
    </div>
  </div>
{/if}
