const API_BASE = '/api';

async function apiCall(endpoint, options = {}) {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'An error occurred' }));
    throw new Error(error.error || 'An error occurred');
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
}

export const api = {
  projects: {
    getAll: () => apiCall('/projects'),
    create: (data) => apiCall('/projects', { method: 'POST', body: JSON.stringify(data) }),
    get: (id) => apiCall(`/projects/${id}`),
    update: (id, data) => apiCall(`/projects/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id) => apiCall(`/projects/${id}`, { method: 'DELETE' }),
  },
  tasks: {
    getByProject: (projectId) => apiCall(`/projects/${projectId}/tasks`),
    create: (projectId, data) => apiCall(`/projects/${projectId}/tasks`, { method: 'POST', body: JSON.stringify(data) }),
    get: (id) => apiCall(`/tasks/${id}`),
    update: (id, data) => apiCall(`/tasks/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    toggleComplete: (id) => apiCall(`/tasks/${id}/complete`, { method: 'PATCH' }),
    delete: (id) => apiCall(`/tasks/${id}`, { method: 'DELETE' }),
  },
  timeLogs: {
    getByTask: (taskId) => apiCall(`/tasks/${taskId}/timelogs`),
    create: (taskId, data) => apiCall(`/tasks/${taskId}/timelogs`, { method: 'POST', body: JSON.stringify(data) }),
    get: (id) => apiCall(`/timelogs/${id}`),
    delete: (id) => apiCall(`/timelogs/${id}`, { method: 'DELETE' }),
  },
};
