const API_BASE = 'http://localhost:3000/api';

async function apiCall(endpoint, options = {}) {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include', // Always send cookies with requests
    ...options,
  });

  if (!response.ok) {
    // Handle 401 Unauthorized - user needs to login
    if (response.status === 401) {
      // Redirect to login by reloading the page
      // The authStore.checkAuth() will handle showing login screen
      window.location.reload();
      return;
    }
    
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
    getByProject: (projectId, page = 1, limit = 10) => {
      console.log(`API call: GET /projects/${projectId}/tasks?page=${page}&limit=${limit}`);
      return apiCall(`/projects/${projectId}/tasks?page=${page}&limit=${limit}`);
    },
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
  comments: {
    getByTask: (taskId) => apiCall(`/tasks/${taskId}/comments`),
    create: (taskId, data) => apiCall(`/tasks/${taskId}/comments`, { method: 'POST', body: JSON.stringify(data) }),
    update: (id, data) => apiCall(`/comments/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id) => apiCall(`/comments/${id}`, { method: 'DELETE' }),
  },
  users: {
    getAll: (page = 1, limit = 20) => apiCall(`/users?page=${page}&limit=${limit}`),
    getById: (id) => apiCall(`/users/${id}`),
    updateRole: (id, role) => apiCall(`/users/${id}/role`, { method: 'PUT', body: JSON.stringify({ role }) }),
    updateActivation: (id, isActive) => apiCall(`/users/${id}/activate`, { method: 'PUT', body: JSON.stringify({ is_active: isActive }) }),
  },
  dashboard: {
    get: () => apiCall('/dashboard'),
  },
  meetings: {
    getNext: () => apiCall('/meetings/next'),
    create: (data) => apiCall('/meetings', { method: 'POST', body: JSON.stringify(data) }),
    getAll: () => apiCall('/meetings'),
    getById: (id) => apiCall(`/meetings/${id}`),
  },
};
