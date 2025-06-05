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
      className={`${isAccept ? 'bg-sky-600/60 hover:bg-sky-800-60' : 'bg-red-600/80 hover:bg-red-800/80'} text-white py-1 px-3 rounded font-bold`}
    >
      {isAccept ? 'Accept' : 'Decline'}
    </button>
  );
}

