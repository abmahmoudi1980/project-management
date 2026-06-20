<script>
  import { createEventDispatcher, tick } from "svelte";
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

  const FOCUSABLE_SELECTOR = [
    'a[href]',
    'button:not([disabled])',
    'input:not([disabled]):not([type="hidden"])',
    'select:not([disabled])',
    'textarea:not([disabled])',
    '[tabindex]:not([tabindex="-1"])',
  ].join(',');

  let modalRef = $state();
  let previousActiveElement = null;
  const titleId = `modal-title-${Math.random().toString(36).slice(2, 9)}`;

  function getFocusable() {
    if (!modalRef) return [];
    return Array.from(modalRef.querySelectorAll(FOCUSABLE_SELECTOR)).filter(
      (el) => !el.hasAttribute('inert') && el.offsetParent !== null
    );
  }

  function handleBackdropClick(event) {
    if (event.target === event.currentTarget) {
      close();
    }
  }

  function handleKeydown(event) {
    if (event.key === "Escape") {
      event.preventDefault();
      close();
      return;
    }
    if (event.key !== "Tab" || !modalRef) return;

    const focusable = getFocusable();
    if (focusable.length === 0) {
      event.preventDefault();
      modalRef.focus();
      return;
    }
    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    const active = document.activeElement;

    if (event.shiftKey && (active === first || !modalRef.contains(active))) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && (active === last || !modalRef.contains(active))) {
      event.preventDefault();
      first.focus();
    }
  }

  function close() {
    dispatch("close");
  }

  $effect(() => {
    if (show) {
      previousActiveElement = document.activeElement;
      tick().then(() => {
        const focusable = getFocusable();
        if (focusable.length > 0) {
          focusable[0].focus();
        } else if (modalRef) {
          modalRef.focus();
        }
      });
    } else if (previousActiveElement && typeof previousActiveElement.focus === 'function') {
      previousActiveElement.focus();
      previousActiveElement = null;
    }
  });
</script>

<svelte:window onkeydown={show ? handleKeydown : undefined} />

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center sm:justify-center bg-black bg-opacity-50"
    onclick={handleBackdropClick}
    onkeydown={handleKeydown}
    role="presentation"
    transition:fade={{ duration: 200 }}
  >
    <div
      bind:this={modalRef}
      class="bg-white sm:rounded-lg shadow-xl w-full sm:w-full {maxWidthClasses[
        maxWidth
      ]} {fullScreen
        ? 'h-full max-h-full flex flex-col sm:max-h-[90vh]'
        : 'max-h-[90vh] sm:max-h-[90vh] my-auto flex flex-col'} sm:m-0 overflow-hidden outline-none"
      transition:scale={{ duration: 200, start: 0.95 }}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? titleId : undefined}
      tabindex="-1"
    >
      <!-- Modal Header -->
      <div
        class="flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4 border-b border-slate-200 flex-shrink-0"
      >
        <h2 id={titleId} class="text-lg sm:text-xl font-semibold text-slate-900">{title}</h2>
        <button
          onclick={close}
          class="p-2 sm:p-0 min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0 text-slate-400 hover:text-slate-600 hover:bg-slate-100 sm:hover:bg-transparent rounded-lg sm:rounded-none transition-colors"
          aria-label="بستن"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            aria-hidden="true"
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
