import { writable, derived } from 'svelte/store';
import { api } from '../lib/api.js';

function createTaskStore() {
  const { subscribe, set, update } = writable([]);

  let currentPage = 1;
  let pageSize = 10;
  let total = 0;
  let hasMore = false;
  let currentProjectId = null;
  let loadingMore = false;

  const load = async (projectId, reset = false) => {
    try {
      if (reset) {
        currentPage = 1;
        currentProjectId = projectId;
      }

      const response = await api.tasks.getByProject(projectId, currentPage, pageSize);
      set(Array.isArray(response.tasks) ? response.tasks : []);

      total = response.total || 0;
      hasMore = response.has_more || false;
    } catch (error) {
      console.error('Failed to load tasks:', error);
      set([]);
      hasMore = false;
    }
  };

  const loadMore = async () => {
    if (loadingMore || !hasMore || !currentProjectId) return;

    loadingMore = true;
    try {
      currentPage += 1;
      const response = await api.tasks.getByProject(currentProjectId, currentPage, pageSize);
      const newTasks = Array.isArray(response.tasks) ? response.tasks : [];

      update(currentTasks => [...currentTasks, ...newTasks]);

      hasMore = response.has_more || false;
    } catch (error) {
      console.error('Failed to load more tasks:', error);
      currentPage -= 1;
    } finally {
      loadingMore = false;
    }
  };

  const reset = () => {
    currentPage = 1;
    total = 0;
    hasMore = false;
    currentProjectId = null;
    loadingMore = false;
    set([]);
  };

  return {
    subscribe,
    currentPage: derived(subscribe, () => currentPage),
    pageSize: derived(subscribe, () => pageSize),
    total: derived(subscribe, () => total),
    hasMore: derived(subscribe, () => hasMore),
    loadingMore: derived(subscribe, () => loadingMore),
    load,
    loadMore,
    reset,
    create: async (projectId, taskData) => {
      const task = await api.tasks.create(projectId, taskData);
      await load(projectId, true);
      return task;
    },
    update: async (id, taskData) => {
      const task = await api.tasks.update(id, taskData);
      update(currentTasks => currentTasks.map(t => t.id === id ? task : t));
      return task;
    },
    toggleComplete: async (id) => {
      const task = await api.tasks.toggleComplete(id);
      update(currentTasks => currentTasks.map(t => t.id === id ? task : t));
      return task;
    },
    delete: async (id) => {
      await api.tasks.delete(id);
      update(currentTasks => currentTasks.filter(t => t.id !== id));
      total -= 1;
    },
    getById: async (id) => {
      const task = await api.tasks.get(id);
      update(currentTasks => {
        const index = currentTasks.findIndex(t => t.id === id);
        if (index >= 0) {
          currentTasks[index] = task;
        }
        return currentTasks;
      });
      return task;
    }
  };
}

export const tasks = createTaskStore();
