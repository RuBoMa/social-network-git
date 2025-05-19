export default function CreateGroup({ onClose }) {
  return (
    <div className="absolute inset-20 bg-white bg-opacity-0 flex items-center justify-center z-50">
      <div className="bg-grey p-6 rounded-lg shadow-lg w-[400px] max-w-full">
        <h2 className="text-xl font-bold mb-4">Create Group</h2>
        <form>
          <input
            type="text"
            className="block w-full mb-3 border p-2 rounded"
            placeholder="Group Name"
          />
          <textarea
            className="block w-full mb-3 border p-2 rounded"
            placeholder="Group Description"
          />
          <div className="flex justify-end space-x-2">
            <button
              type="button"
              className="px-4 py-2 border rounded"
              onClick={onClose}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded"
            >
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
