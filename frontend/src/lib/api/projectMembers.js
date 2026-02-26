import { api } from '../api.js';

/**
 * Get all members of a project
 * @param {string} projectId
 * @returns {Promise<Array>}
 */
export async function getProjectMembers(projectId) {
  return api.projectMembers.getAll(projectId);
}

/**
 * Add a member to a project
 * @param {string} projectId
 * @param {Object} data - { user_id, project_role_id }
 * @returns {Promise<Object>}
 */
export async function addProjectMember(projectId, data) {
  return api.projectMembers.add(projectId, data);
}

/**
 * Get eligible users who can be added to a project
 * @param {string} projectId
 * @returns {Promise<Array>}
 */
export async function getEligibleProjectUsers(projectId) {
  return api.projectMembers.getEligibleUsers(projectId);
}

/**
 * Get all active project roles
 * @returns {Promise<Array>}
 */
export async function getProjectRoles() {
  return api.projectRoles.getAll();
}

/**
 * Remove a member from a project
 * @param {string} projectId
 * @param {string} userId
 * @returns {Promise<Object>}
 */
export async function removeProjectMember(projectId, userId) {
  return api.projectMembers.remove(projectId, userId);
}
