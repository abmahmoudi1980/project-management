import { writable } from 'svelte/store';

function createToastStore() {
  const { subscribe, update } = writable([]);

  let nextId = 1;

  function push(message, type = 'info', duration = 4000) {
    const id = nextId++;
    const toast = { id, message, type, duration };

    update((toasts) => [...toasts, toast]);

    if (duration > 0) {
      setTimeout(() => dismiss(id), duration);
    }

    return id;
  }

  function dismiss(id) {
    update((toasts) => toasts.filter((t) => t.id !== id));
  }

  return {
    subscribe,
    success: (message, duration) => push(message, 'success', duration),
    error: (message, duration) => push(message, 'error', duration),
    info: (message, duration) => push(message, 'info', duration),
    warning: (message, duration) => push(message, 'warning', duration),
    dismiss,
  };
}

export const toasts = createToastStore();
