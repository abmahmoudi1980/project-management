import { api } from '../api.js';

/**
 * Get all members of a project
 * @param {string} projectId
 * @returns {Promise<Array>}
 */
export async function getProjectMembers(projectId) {
  return api.call(`/projects/${projectId}/members`);
}

/**
 * Add a member to a project
 * @param {string} projectId
 * @param {Object} data - { user_id, project_role_id }
 * @returns {Promise<Object>}
 */
export async function addProjectMember(projectId, data) {
  return api.call(`/projects/${projectId}/members`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

/**
 * Get eligible users who can be added to a project
 * @param {string} projectId
 * @returns {Promise<Array>}
 */
export async function getEligibleProjectUsers(projectId) {
  return api.call(`/projects/${projectId}/members/eligible-users`);
}

/**
 * Get all active project roles
 * @returns {Promise<Array>}
 */
export async function getProjectRoles() {
  return api.call('/project-roles');
}
