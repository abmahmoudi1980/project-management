export async function getDashboardData() {
  const response = await fetch('/api/dashboard', {
    headers: {
      'Accept': 'application/json'
    }
  });
  
  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Unauthorized');
    }
    throw new Error('Failed to fetch dashboard data');
  }
  
  return await response.json();
}
