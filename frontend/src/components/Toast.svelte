<script>
  import { fly } from "svelte/transition";
  import { cubicOut } from "svelte/easing";
  import { toasts } from "../lib/toastStore.js";

  const typeStyles = {
    success: "bg-success-50 border-success-500 text-success-700",
    error: "bg-danger-50 border-danger-500 text-danger-700",
    warning: "bg-warning-50 border-warning-500 text-warning-700",
    info: "bg-info-50 border-info-500 text-info-700",
  };

  const typeIcons = {
    success: "M5 13l4 4L19 7",
    error: "M6 18L18 6M6 6l12 12",
    warning: "M12 9v2m0 4h.01M5 19h14a2 2 0 001.84-2.75L13.74 4a2 2 0 00-3.48 0L3.16 16.25A2 2 0 005 19z",
    info: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
  };

  const typeAriaRole = {
    success: "status",
    info: "status",
    warning: "alert",
    error: "alert",
  };
</script>

<div
  class="fixed bottom-4 left-1/2 -translate-x-1/2 sm:left-auto sm:translate-x-0 sm:end-4 z-[100] flex flex-col gap-2 w-full max-w-sm px-4 sm:px-0 pointer-events-none"
  aria-live="polite"
  aria-atomic="false"
>
  {#each $toasts as toast (toast.id)}
    <div
      role={typeAriaRole[toast.type]}
      class="pointer-events-auto flex items-start gap-3 p-4 rounded-lg border-s-4 shadow-lg bg-white {typeStyles[toast.type]}"
      transition:fly={{ y: 20, duration: 250, easing: cubicOut }}
    >
      <svg
        class="w-5 h-5 flex-shrink-0 mt-0.5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d={typeIcons[toast.type]}
        />
      </svg>
      <p class="flex-1 text-sm font-medium text-slate-800">{toast.message}</p>
      <button
        type="button"
        onclick={() => toasts.dismiss(toast.id)}
        class="flex-shrink-0 text-slate-400 hover:text-slate-600 transition-colors p-1 -m-1"
        aria-label="بستن اعلان"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {/each}
</div>
