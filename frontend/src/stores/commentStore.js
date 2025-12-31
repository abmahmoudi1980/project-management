import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createCommentStore() {
  const { subscribe, set, update } = writable([]);

  return {
    subscribe,
    load: async (taskId) => {
      try {
        const response = await api.comments.getByTask(taskId);
        set(response.data?.comments || []);
      } catch (error) {
        console.error('Failed to load comments:', error);
        set([]);
      }
    },
    create: async (taskId, commentData) => {
      try {
        const response = await api.comments.create(taskId, commentData);
        update(currentComments => [...currentComments, response.data?.comment]);
        return response.data?.comment;
      } catch (error) {
        console.error('Failed to create comment:', error);
        throw error;
      }
    },
    update: async (id, commentData) => {
      try {
        const response = await api.comments.update(id, commentData);
        update(currentComments => currentComments.map(c => c.id === id ? response.data?.comment : c));
        return response.data?.comment;
      } catch (error) {
        console.error('Failed to update comment:', error);
        throw error;
      }
    },
    delete: async (id) => {
      try {
        await api.comments.delete(id);
        update(currentComments => currentComments.filter(c => c.id !== id));
      } catch (error) {
        console.error('Failed to delete comment:', error);
        throw error;
      }
    }
  };
}

export const comments = createCommentStore();
