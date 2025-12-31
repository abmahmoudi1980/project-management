import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createProjectStore() {
  const { subscribe, set, update } = writable([]);

  return {
    subscribe,
    load: async () => {
      try {
        console.debug('[projectStore] loading projects...');
        const projects = await api.projects.getAll();
        console.debug('[projectStore] projects loaded:', projects);
        set(projects || []);
      } catch (err) {
        console.error('[projectStore] failed to load projects:', err);
        set([]);
        throw err;
      }
    },
    create: async (projectData) => {
      const project = await api.projects.create(projectData);
      update(currentProjects => [project, ...(currentProjects || [])]);
      return project;
    },
    update: async (id, projectData) => {
      const project = await api.projects.update(id, projectData);
      update(currentProjects => (currentProjects || []).map(p => p.id === id ? project : p));
      return project;
    },
    delete: async (id) => {
      await api.projects.delete(id);
      update(currentProjects => (currentProjects || []).filter(p => p.id !== id));
    }
  };
}

export const projects = createProjectStore();
