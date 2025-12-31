<script>
  import { createEventDispatcher } from "svelte";
  import { fade, scale } from "svelte/transition";

  let { show = false, title = "", maxWidth = "2xl", children, fullScreen = true } = $props();
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

<svelte:window onkeydown={handleKeydown} />

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center sm:p-4 bg-black bg-opacity-50"
    onclick={handleBackdropClick}
    onkeydown={handleKeydown}
    tabindex="-1"
    transition:fade={{ duration: 200 }}
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white sm:rounded-lg shadow-xl w-full sm:w-full {maxWidthClasses[
        maxWidth
      ]} {fullScreen
        ? 'h-full sm:max-h-[90vh] sm:flex sm:flex-col'
        : 'max-h-[90vh] sm:max-h-[90vh] my-auto'} sm:m-0"
      transition:scale={{ duration: 200, start: 0.95 }}
    >
      <!-- Modal Header -->
      <div
        class="flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4 border-b border-slate-200 flex-shrink-0"
      >
        <h2 class="text-lg sm:text-xl font-semibold text-slate-900">{title}</h2>
        <button
          onclick={close}
          class="p-2 sm:p-0 min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0 text-slate-400 hover:text-slate-600 hover:bg-slate-100 sm:hover:bg-transparent rounded-lg sm:rounded-none transition-colors"
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
      <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4 sm:py-6">
        {@render children()}
      </div>
    </div>
  </div>
{/if}
