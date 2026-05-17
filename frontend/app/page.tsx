"use client";

import { useState } from "react";
import { motion } from "framer-motion";

type Issue = {
  id: number;
  error: string;
  causes: string[];
  fixes: string[];
  debug_steps: string[];
  tags: string[];
  created_at: string;
};

export default function Home() {
  const [error, setError] = useState("");
  const [causes, setCauses] = useState("");
  const [fixes, setFixes] = useState("");
  const [debugSteps, setDebugSteps] = useState("");
  const [tags, setTags] = useState("");

  const [results, setResults] = useState<Issue[]>([]);
  const [suggestions, setSuggestions] = useState<Issue[]>([]);

  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [searched, setSearched] = useState(false);

  const [similarIssue, setSimilarIssue] = useState<Issue | null>(null);

  const API = process.env.NEXT_PUBLIC_API_URL;

  // 🔍 SEARCH
  const handleSearch = async () => {
    setLoading(true);
    setMessage("");
    setSearched(true);
    setResults([]);

    try {
      const res = await fetch(`${API}/search?error=${error}`);
      const data = await res.json();

      setResults(data.results || []);
    } catch {
      setMessage("Something went wrong ❌");
    }

    setLoading(false);
  };

  // 💾 SAVE
  const handleSave = async () => {
    setLoading(true);
    setMessage("");
    setSimilarIssue(null);

    try {
      const res = await fetch(`${API}/issue`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer devops-secret-key",
        },
        body: JSON.stringify({
          error,

          causes: causes
            .split("\n")
            .map((c) => c.trim())
            .filter(Boolean),

          fixes: fixes
            .split("\n")
            .map((f) => f.trim())
            .filter(Boolean),

          debug_steps: debugSteps
            .split("\n")
            .map((s) => s.trim())
            .filter(Boolean),

          tags: tags
            .split(",")
            .map((tag) => tag.trim())
            .filter(Boolean),
        }),
      });

      const data = await res.json();

      if (data.similar) {
        setSimilarIssue(data.similar);
        setMessage(data.message);
      } else {
        setMessage(data.message || "Saved successfully ✅");

        setError("");
        setCauses("");
        setFixes("");
        setDebugSteps("");
        setTags("");
      }
    } catch {
      setMessage("Failed to save ❌");
    }

    setLoading(false);
  };

  // ⚡ LIVE SUGGESTIONS
  const handleInputChange = async (value: string) => {
    setError(value);

    if (value.length > 2) {
      try {
        const res = await fetch(`${API}/search?error=${value}`);
        const data = await res.json();

        setSuggestions(data);
      } catch {
        setSuggestions([]);
      }
    } else {
      setSuggestions([]);
    }
  };

  const handleDelete = async (id: number) => {

    try {

      const res = await fetch(
        `${API}/delete?id=${id}`,
        {
          method: "DELETE",
          headers: {
            "Authorization": "Bearer devops-secret-key",
          },
        }
      );

      // 🔥 read raw response first
      const text = await res.text();

      console.log("DELETE RESPONSE:", text);

      // safely parse JSON
      let data;

      try {
        data = JSON.parse(text);
      } catch {
        data = {
          message: "Issue deleted successfully ✅",
        };
      }

      setMessage(data.message);

      // 🔥 remove deleted issue instantly
      setResults((prev) =>
        prev.filter((item) => item.id !== id)
      );

    } catch (err) {

      console.error(err);

      setMessage("Failed to delete ❌");

    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="min-h-screen bg-gradient-to-br from-purple-900 via-purple-800 to-pink-700 text-white flex flex-col items-center p-6"
    >
      {/* TITLE */}
      <motion.h1
        initial={{ y: -20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="text-4xl font-bold mb-6 text-center"
      >
        DevOps Memory Assistant 🚀
      </motion.h1>

      {/* FORM */}
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        className="bg-white/10 backdrop-blur-lg p-6 rounded-2xl shadow-lg w-full max-w-md space-y-4"
      >
        {/* ERROR INPUT */}
        <div>
          <input
            value={error}
            onChange={(e) => handleInputChange(e.target.value)}
            placeholder="Error (e.g. CrashLoopBackOff)"
            className="w-full p-3 rounded-lg bg-white/90 text-black focus:ring-2 focus:ring-pink-400"
          />

          {/* LIVE SUGGESTIONS */}
          {suggestions.length > 0 && (
            <div className="bg-white text-black rounded-lg mt-2 shadow-lg max-h-40 overflow-y-auto">
              {suggestions.map((item, index) => (
                <div
                  key={index}
                  className="p-2 hover:bg-gray-200 cursor-pointer"
                  onClick={() => {
                    setError(item.error);
                    setSuggestions([]);
                  }}
                >
                  {item.error}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* CAUSE */}
        <textarea
          value={causes}
          onChange={(e) => setCauses(e.target.value)}
          placeholder="Causes (one per line)"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        {/* FIX */}
        <textarea
          value={fixes}
          onChange={(e) => setFixes(e.target.value)}
          placeholder="Fixes (one per line)"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        {/* STEPS */}
        <textarea
          value={debugSteps}
          onChange={(e) => setDebugSteps(e.target.value)}
          placeholder="Debug steps (one per line)"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        {/* TAGS */}
        <input
          value={tags}
          onChange={(e) => setTags(e.target.value)}
          placeholder="Tags (kubernetes,docker,network)"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        {/* BUTTONS */}
        <div className="flex gap-4">
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={handleSearch}
            disabled={loading}
            className="flex-1 bg-pink-500 py-2 rounded-lg hover:bg-pink-600 transition"
          >
            {loading ? "Searching..." : "Search"}
          </motion.button>

          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={handleSave}
            disabled={loading}
            className="flex-1 bg-purple-500 py-2 rounded-lg hover:bg-purple-600 transition"
          >
            {loading ? "Saving..." : "Save"}
          </motion.button>
        </div>

        {/* MESSAGE */}
        {message && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center text-sm bg-black/30 p-2 rounded"
          >
            {message}
          </motion.div>
        )}

        {/* SIMILAR ISSUE */}
        {similarIssue && (
          <div className="mt-4 bg-yellow-200 text-black p-4 rounded-lg">
            <p className="font-bold text-lg">
              ⚠️ Similar issue already exists
            </p>

            <p className="mt-3">
              <strong>🚨 Error:</strong>
              <br />
              {similarIssue.error}
            </p>

            <div className="mt-3">
              <strong>📌 Causes:</strong>

              <ul className="list-disc ml-5 mt-1">
                {similarIssue.causes?.map((cause, idx) => (
                  <li key={idx}>{cause}</li>
                ))}
              </ul>
            </div>

            <div className="mt-3">
              <strong>✅ Fixes:</strong>

              <ul className="list-disc ml-5 mt-1">
                {similarIssue.fixes?.map((fix, idx) => (
                  <li key={idx}>{fix}</li>
                ))}
              </ul>
            </div>

            <div className="mt-3">
              <strong>🛠 Debug Steps:</strong>

              <ul className="list-disc ml-5 mt-1">
                {similarIssue.debug_steps?.map((step, idx) => (
                  <li key={idx}>{step}</li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </motion.div>

      {/* RESULTS */}
      <div className="mt-8 w-full max-w-md space-y-4">
        {searched && results.length === 0 && !loading && (
          <p className="text-center text-gray-300">
            No results found 👀
          </p>
        )}

        {results.map((item, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
            className="bg-white/10 backdrop-blur-lg p-4 rounded-xl border border-white/20 hover:scale-[1.03] transition"
          >
            <p className="text-pink-300 font-bold text-lg">
              🚨 {item.error}
            </p>

            <div className="mt-3">

              <div className="mt-3">
                <p className="text-yellow-300 font-semibold">
                  📌 Causes
                </p>

                <ul className="list-disc ml-5 text-sm text-gray-200 mt-1">
                  {item.causes?.map((cause, idx) => (
                    <li key={idx}>{cause}</li>
                  ))}
                </ul>
              </div>
            </div>

            <div className="mt-3">
              <p className="text-green-300 font-semibold">
                ✅ Fixes
              </p>

              <ul className="list-disc ml-5 text-sm text-gray-200 mt-1">
                {item.fixes?.map((fix, idx) => (
                  <li key={idx}>{fix}</li>
                ))}
              </ul>
            </div>

            {item.debug_steps &&
              item.debug_steps.length > 0 && (
                <div className="mt-3">
                  <p className="text-blue-300 font-semibold">
                    🛠 Debug Steps
                  </p>

                  <ul className="list-disc ml-5 text-sm text-gray-200 mt-1">
                    {item.debug_steps.map((step, idx) => (
                      <li key={idx}>{step}</li>
                    ))}
                  </ul>
                </div>
              )}

            {item.tags && item.tags.length > 0 && (
              <div className="mt-3 flex flex-wrap gap-2">
                {item.tags.map((tag, idx) => (
                  <span
                    key={idx}
                    className="bg-pink-500/30 text-pink-200 px-2 py-1 rounded-full text-xs"
                  >
                    #{tag}
                  </span>
                ))}
              </div>
            )}

            <>
              {item.created_at && (
                <p className="text-xs text-gray-400 mt-4">
                  🕒 {new Date(item.created_at).toLocaleString()}
                </p>
              )}

              <button
                onClick={() => handleDelete(item.id)}
                className="mt-4 bg-red-500 hover:bg-red-600 px-3 py-1 rounded text-sm"
              >
                🗑 Delete
              </button>
            </>
          </motion.div>

        ))}
      </div>
    </motion.div>
  );
}