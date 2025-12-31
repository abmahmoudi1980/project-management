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
      if (reset) {
        update(state => ({ ...state, currentPage: 1, currentProjectId: projectId, loadingMore: true }));
      } else {
        update(state => ({ ...state, loadingMore: true }));
      }

      const currentPage = reset ? 1 : initialState.currentPage;
      const response = await api.tasks.getByProject(projectId, currentPage, initialState.pageSize);
      console.log('Initial load response:', { projectId, currentPage, pageSize: initialState.pageSize, response });

      update(state => ({
        ...state,
        tasks: Array.isArray(response.tasks) ? response.tasks : [],
        total: response.total || 0,
        hasMore: response.has_more || false,
        loadingMore: false
      }));
      console.log('After initial load:', { total: initialState.total, hasMore: initialState.hasMore });
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
    update(state => {
      console.log('loadMore called from state:', { loadingMore: state.loadingMore, hasMore: state.hasMore, currentProjectId: state.currentProjectId, currentPage: state.currentPage });
      if (state.loadingMore || !state.hasMore || !state.currentProjectId) return state;
      return { ...state, loadingMore: true };
    });

    let currentState;
    subscribe(state => currentState = state)();

    if (!currentState || currentState.loadingMore || !currentState.hasMore || !currentState.currentProjectId) {
      return;
    }

    try {
      const nextPage = currentState.currentPage + 1;
      console.log('Fetching page:', nextPage);
      const response = await api.tasks.getByProject(currentState.currentProjectId, nextPage, currentState.pageSize);
      const newTasks = Array.isArray(response.tasks) ? response.tasks : [];
      console.log('Response:', response);

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
