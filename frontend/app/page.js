export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-800 via-purple-900 to-pink-700 text-white flex flex-col items-center justify-center p-6">

      <h1 className="text-4xl font-bold mb-6">
        DevOps Memory Assistant 🚀
      </h1>

      <input
        placeholder="Enter error (e.g. CrashLoopBackOff)"
        className="w-full max-w-md p-3 rounded-lg text-black mb-4"
      />

      <div className="flex gap-4">
        <button className="bg-pink-500 px-5 py-2 rounded-lg hover:bg-pink-600">
          Search
        </button>

        <button className="bg-purple-500 px-5 py-2 rounded-lg hover:bg-purple-600">
          Save
        </button>
      </div>

    </div>
  );
}