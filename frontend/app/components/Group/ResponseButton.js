export default function ResponseButton({ groupId, requestId, status, onResponse }) {
  const isAccept = status === 'accepted';

  const handleClick = async () => {
    await fetch(`http://localhost:8080/api/request`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group: { group_id: Number(groupId) },
        status: status,
        request_id: Number(requestId),
      }),
    });
    onResponse();
  };

  return (
    <button
      onClick={handleClick}
      className={`${isAccept ? 'bg-green-700 hover:bg-green-600' : 'bg-red-700 hover:bg-red-600'} text-white py-1 px-3 mx-2 mt-5 rounded`}
    >
      {isAccept ? 'Accept' : 'Decline'}
    </button>
  );
}

