export async function getNextMeeting() {
  const response = await fetch('/api/meetings/next', {
    headers: {
      'Accept': 'application/json'
    }
  });
  
  if (response.status === 204) {
    return null;
  }
  
  if (!response.ok) {
    throw new Error('Failed to fetch next meeting');
  }
  
  return await response.json();
}

export async function createMeeting(data) {
  const response = await fetch('/api/meetings', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    body: JSON.stringify(data)
  });
  
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || 'Failed to create meeting');
  }
  
  return await response.json();
}

export async function listMeetings() {
  const response = await fetch('/api/meetings', {
    headers: {
      'Accept': 'application/json'
    }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch meetings');
  }
  
  return await response.json();
}
