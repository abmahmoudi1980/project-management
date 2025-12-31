import { writable, derived } from 'svelte/store';
import { api } from '../lib/api.js';

function createTaskStore() {
  const initialState = {
    tasks: [],
    currentPage: 1,
    pageSize: 10,
    total: 0,
    hasMore: false,
    currentProjectId: null,
    loadingMore: false
  };

  const { subscribe, set, update } = writable(initialState);

  const load = async (projectId, reset = false) => {
    try {
      // Always update currentProjectId when loading, and reset page to 1
      update(state => ({ ...state, currentPage: 1, currentProjectId: projectId, loadingMore: true }));

      const response = await api.tasks.getByProject(projectId, 1, initialState.pageSize);
      console.log('Initial load response:', { projectId, pageSize: initialState.pageSize, taskCount: response.tasks?.length, total: response.total, hasMore: response.has_more });

      update(state => ({
        ...state,
        tasks: Array.isArray(response.tasks) ? response.tasks : [],
        total: response.total || 0,
        hasMore: response.has_more || false,
        loadingMore: false
      }));
    } catch (error) {
      console.error('Failed to load tasks:', error);
      update(state => ({
        ...state,
        tasks: [],
        hasMore: false,
        loadingMore: false
      }));
    }
  };

  const loadMore = async () => {
    // First check if we can load more
    let canLoadMore = false;
    let currentState;
    
    subscribe(state => {
      currentState = state;
      if (!state.loadingMore && state.hasMore && state.currentProjectId) {
        canLoadMore = true;
      }
    })();

    if (!canLoadMore) {
      console.log('Cannot load more:', { 
        loadingMore: currentState?.loadingMore, 
        hasMore: currentState?.hasMore, 
        currentProjectId: currentState?.currentProjectId 
      });
      return;
    }

    // Set loading state
    update(state => ({ ...state, loadingMore: true }));

    try {
      // Re-subscribe to get the latest state after update
      let latestState;
      subscribe(state => latestState = state)();
      
      const nextPage = latestState.currentPage + 1;
      console.log('Fetching page:', nextPage, 'for project:', latestState.currentProjectId);
      const response = await api.tasks.getByProject(latestState.currentProjectId, nextPage, latestState.pageSize);
      const newTasks = Array.isArray(response.tasks) ? response.tasks : [];
      console.log('LoadMore response:', { nextPage, tasksReceived: newTasks.length, hasMore: response.has_more });

      update(state => ({
        ...state,
        tasks: [...state.tasks, ...newTasks],
        currentPage: nextPage,
        hasMore: response.has_more || false,
        loadingMore: false
      }));
    } catch (error) {
      console.error('Failed to load more tasks:', error);
      update(state => ({ ...state, loadingMore: false }));
    }
  };

  const reset = () => {
    set(initialState);
  };

  return {
    subscribe,
    currentPage: derived(subscribe, s => s.currentPage),
    pageSize: derived(subscribe, s => s.pageSize),
    total: derived(subscribe, s => s.total),
    hasMore: derived(subscribe, s => s.hasMore),
    loadingMore: derived(subscribe, s => s.loadingMore),
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
