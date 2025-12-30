import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createProjectStore() {
  const { subscribe, set, update } = writable([]);

  return {
    subscribe,
    load: async () => {
      const projects = await api.projects.getAll();
      set(projects);
    },
    create: async (projectData) => {
      const project = await api.projects.create(projectData);
      update(currentProjects => [project, ...currentProjects]);
      return project;
    },
    update: async (id, projectData) => {
      const project = await api.projects.update(id, projectData);
      update(currentProjects => currentProjects.map(p => p.id === id ? project : p));
      return project;
    },
    delete: async (id) => {
      await api.projects.delete(id);
      update(currentProjects => currentProjects.filter(p => p.id !== id));
    }
  };
}

export const projects = createProjectStore();
