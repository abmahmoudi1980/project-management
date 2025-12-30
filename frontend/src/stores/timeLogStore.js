import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createTimeLogStore() {
  const { subscribe, set } = writable([]);

  return {
    subscribe,
    load: async (taskId) => {
      const timeLogs = await api.timeLogs.getByTask(taskId);
      set(timeLogs);
    },
    create: async (taskId, timeLogData) => {
      const timeLog = await api.timeLogs.create(taskId, timeLogData);
      set(currentLogs => [timeLog, ...currentLogs]);
      return timeLog;
    },
    delete: async (id) => {
      await api.timeLogs.delete(id);
      set(currentLogs => currentLogs.filter(l => l.id !== id));
    }
  };
}

export const timeLogs = createTimeLogStore();
